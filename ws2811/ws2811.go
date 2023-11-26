package ws2811

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lampctl/lampctl/db"
	"github.com/lampctl/lampctl/registry"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	KeyNumLeds = "ws2811.numLeds"

	ProviderID = "ws2811"
	GroupID    = "ws2811"
)

var errNoLEDs = errors.New("LED count is set to 0")

// Ws2811 implements the Provider interface for ws2811.
type Ws2811 struct {
	mutex   sync.RWMutex
	db      *db.Conn
	ws      *ws2811.WS2811
	numLEDs int
}

func (w *Ws2811) init() error {

	// (Re)create the string if there are more than 0 LEDs
	if w.numLEDs > 0 {

		// Override the default options to max the brightness & set count
		o := ws2811.DefaultOptions
		o.Channels[0].Brightness = 255
		o.Channels[0].LedCount = w.numLEDs

		// Create the string
		ws, err := ws2811.MakeWS2811(&o)
		if err != nil {
			return err
		}

		// Initialize the string
		if err := ws.Init(); err != nil {
			return err
		}

		// Assign the string *only* on success
		w.ws = ws
	}

	return nil
}

func (w *Ws2811) free() {
	if w.ws != nil {
		w.ws.Fini()
		w.ws = nil
	}
}

func New(cfg *Config) (*Ws2811, error) {
	v, err := cfg.DB.GetIntSetting(KeyNumLeds, 0)
	if err != nil {
		return nil, err
	}
	w := &Ws2811{
		db:      cfg.DB,
		numLEDs: v,
	}
	if err := w.init(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Ws2811) ID() string {
	return ProviderID
}

func (w *Ws2811) Name() string {
	return "ws2811"
}

func (w *Ws2811) Init(api *gin.RouterGroup) error {
	api.POST("/ws2811", w.api_ws2811_leds_POST)
	return nil
}

func (w *Ws2811) Close() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.free()
}

func (w *Ws2811) Groups() []*registry.Group {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	if w.ws != nil {
		return []*registry.Group{
			{
				ID:   GroupID,
				Name: "LEDs on GPIO 18",
			},
		}
	} else {
		return []*registry.Group{}
	}
}

func (w *Ws2811) Lamps() []*registry.Lamp {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	lamps := []*registry.Lamp{}
	for i := 0; i < w.numLEDs; i++ {
		lamps = append(lamps, &registry.Lamp{
			ID:      fmt.Sprint(i),
			Name:    fmt.Sprintf("LED %d", i),
			GroupID: GroupID,
			State:   false,
		})
	}
	return lamps
}

func (w *Ws2811) Apply(changes []*registry.Change) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.ws == nil {
		return errNoLEDs
	}
	for _, c := range changes {
		if c.GroupID != GroupID {
			return fmt.Errorf("invalid group ID %s", c.GroupID)
		}
		i, err := strconv.Atoi(c.LampID)
		if err != nil {
			return err
		}
		if i < 0 || i >= w.numLEDs {
			return fmt.Errorf("invalid lamp ID %d", i)
		}
		var color uint32
		if c.State {
			color = 0xffffff
		}
		w.ws.Leds(0)[i] = color
	}
	return w.ws.Render()
}

func (w *Ws2811) ApplyToAll(change *registry.Change) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.ws == nil {
		return errNoLEDs
	}
	var color uint32
	if change.State {
		color = 0xffffff
	}
	for i := 0; i < len(w.ws.Leds(0)); i++ {
		w.ws.Leds(0)[i] = color
	}
	return w.ws.Render()
}
