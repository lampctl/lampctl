package hue

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lampctl/lampctl/db"
	hue_db "github.com/lampctl/lampctl/hue/db"
	"github.com/lampctl/lampctl/registry"
)

const ProviderID = "hue"

// Hue implements the Provider interface for Philips Hue wireless products that
// are connected to a bridge accessible over the network.
type Hue struct {
	mutex   sync.RWMutex
	db      *db.Conn
	bridges map[string]*Bridge
}

// New creates a new Hue instance.
func New(cfg *Config) (*Hue, error) {
	h := &Hue{
		db:      cfg.DB,
		bridges: make(map[string]*Bridge),
	}
	if err := h.db.AutoMigrate(&hue_db.Bridge{}); err != nil {
		return nil, err
	}
	bridges := []*hue_db.Bridge{}
	if err := h.db.Find(&bridges).Error; err != nil {
		return nil, err
	}
	for _, b := range bridges {
		v := NewBridge(b)
		if err := v.Init(); err != nil {
			return nil, err
		}
		h.bridges[b.ID] = v
	}
	return h, nil
}

func (h *Hue) ID() string {
	return ProviderID
}

func (h *Hue) Name() string {
	return "Philips Hue"
}

func (h *Hue) Init(api *gin.RouterGroup) error {
	api.POST("/bridges", h.api_hue_bridges_POST)
	api.DELETE("/bridges/:id", h.api_hue_bridges_id_DELETE)
	return nil
}

func (h *Hue) Close() {}

func (h *Hue) Groups() []*registry.Group {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	groups := []*registry.Group{}
	for _, b := range h.bridges {
		groups = append(groups, &registry.Group{
			ID:   fmt.Sprint(b.ID),
			Name: b.Host,
		})
	}
	return groups
}

func (h *Hue) Lamps() []*registry.Lamp {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	lights := []*registry.Lamp{}
	for _, b := range h.bridges {
		for _, r := range b.resources {
			lights = append(lights, &registry.Lamp{
				ID:      r.Resource.ID,
				Name:    r.Name,
				GroupID: fmt.Sprint(b.ID),
				State:   r.Resource.On.On,
			})
		}
	}
	return lights
}

func (h *Hue) Apply(changes []*registry.Change) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	for _, c := range changes {
		b, ok := h.bridges[c.GroupID]
		if !ok {
			return registry.ErrInvalidGroup
		}
		if err := b.setState(c.LampID, c.State, c.Brightness, c.Color, c.Duration); err != nil {
			return err
		}
	}
	return nil
}

func (h *Hue) ApplyToAll(change *registry.Change) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	for _, b := range h.bridges {
		if err := b.setState(
			b.allResourceID,
			change.State,
			change.Brightness,
			change.Color,
			change.Duration,
		); err != nil {
			return err
		}
	}
	return nil
}
