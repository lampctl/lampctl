package ws2811

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lampctl/lampctl/db"
	"github.com/lampctl/lampctl/registry"
	ws2811_db "github.com/lampctl/lampctl/ws2811/db"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const ProviderID = "ws2811"

// Ws2811 implements the Provider interface for ws2811.
type Ws2811 struct {
	mutex    sync.RWMutex
	db       *db.Conn
	ws       *ws2811.WS2811
	channels []*ws2811_db.Channel
}

func New(cfg *Config) (*Ws2811, error) {
	w := &Ws2811{
		db:       cfg.DB,
		channels: []*ws2811_db.Channel{},
	}
	if err := w.db.AutoMigrate(&ws2811_db.Channel{}); err != nil {
		return nil, err
	}
	if err := w.db.Find(&w.channels).Error; err != nil {
		return nil, err
	}
	channelOptions := []ws2811.ChannelOption{}
	for _, c := range w.channels {
		channelOptions = append(channelOptions, ws2811.ChannelOption{
			GpioPin:    c.GpioPin,
			LedCount:   c.LedCount,
			Brightness: 255,
			StripeType: ws2811.WS2812Strip,
		})
	}
	ws, err := ws2811.MakeWS2811(&ws2811.Option{
		Frequency: ws2811.TargetFreq,
		DmaNum:    ws2811.DefaultDmaNum,
		Channels:  channelOptions,
	})
	if err != nil {
		return nil, err
	}
	if err := ws.Init(); err != nil {
		return nil, err
	}
	w.ws = ws
	return w, nil
}

func (w *Ws2811) ID() string {
	return ProviderID
}

func (w *Ws2811) Name() string {
	return "ws2811"
}

func (w *Ws2811) Init(api *gin.RouterGroup) error {
	api.POST("/ws2811/channels", w.api_ws2811_channels_POST)
	return nil
}

func (w *Ws2811) Close() {
	w.ws.Fini()
}

func (w *Ws2811) Groups() []*registry.Group {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	groups := []*registry.Group{}
	for _, c := range w.channels {
		groups = append(groups, &registry.Group{
			ID:   fmt.Sprint(c.ID),
			Name: fmt.Sprintf("GPIO Pin %d", c.GpioPin),
		})
	}
	return groups
}

func (w *Ws2811) Lamps() []*registry.Lamp {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	lamps := []*registry.Lamp{}
	for i, c := range w.channels {
		for j := 0; j < c.LedCount; j++ {
			lamps = append(lamps, &registry.Lamp{
				ID:      fmt.Sprint(j),
				Name:    fmt.Sprintf("LED %03d", j+1),
				GroupID: fmt.Sprint(i),
				State:   false,
			})
		}
	}
	return lamps
}

func (w *Ws2811) Apply(changes []*registry.Change) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	for _, c := range changes {
		chOff, err := strconv.Atoi(c.GroupID)
		if err != nil {
			return err
		}
		if chOff < 0 || chOff >= len(w.channels) {
			return fmt.Errorf("invalid group ID %d", chOff)
		}
		idx, err := strconv.Atoi(c.LampID)
		if err != nil {
			return err
		}
		if idx < 0 || idx >= w.channels[chOff].LedCount {
			return fmt.Errorf("invalid lamp ID %d", idx)
		}
		var color uint32
		if c.State {
			color = 0xffffff
		} else {
			color = 0x000000
		}
		w.ws.Leds(chOff)[idx] = color
	}
	return nil
}

func (w *Ws2811) ApplyToAll(change *registry.Change) error {
	return nil
}
