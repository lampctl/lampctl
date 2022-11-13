package hue

import (
	"fmt"
	"sort"
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
	bridges map[int64]*Bridge
}

// New creates a new Hue instance.
func New(cfg *Config) (*Hue, error) {
	h := &Hue{
		db:      cfg.DB,
		bridges: make(map[int64]*Bridge),
	}
	if err := h.db.AutoMigrate(&hue_db.Bridge{}); err != nil {
		return nil, err
	}
	bridges := []*hue_db.Bridge{}
	if err := h.db.Find(&bridges).Error; err != nil {
		return nil, err
	}
	for _, b := range bridges {
		v, err := NewBridge(b)
		if err != nil {
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

type byName []*registry.Lamp

func (n byName) Len() int           { return len(n) }
func (n byName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n byName) Less(i, j int) bool { return n[i].Name < n[j].Name }

func (h *Hue) Lamps() []*registry.Lamp {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	lights := []*registry.Lamp{}
	for _, b := range h.bridges {
		for _, l := range b.lights {
			lights = append(lights, &registry.Lamp{
				ID:      l.ID,
				Name:    l.Metadata.Name,
				GroupID: fmt.Sprint(b.ID),
				State:   l.On.On,
			})
		}
	}
	sort.Sort(byName(lights))
	return lights
}

func (h *Hue) Apply(changes []*registry.Change) error {
	return nil
}
