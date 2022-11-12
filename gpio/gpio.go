package gpio

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/lampctl/lampctl/db"
	register_db "github.com/lampctl/lampctl/gpio/db"
	"github.com/lampctl/lampctl/registry"
)

const ProviderID = "gpio"

var (
	errInvalidRegister = errors.New("invalid register specified")
	errInvalidChannel  = errors.New("invalid channel specified")
)

// GPIO implements the Provider interface for shift registers connected to
// GPIO pins on a Raspberry Pi.
type GPIO struct {
	mutex     sync.RWMutex
	db        *db.Conn
	registers map[int64]*Register
}

func (g *GPIO) findRegister(id string) (*Register, error) {
	v, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	r, ok := g.registers[v]
	if !ok {
		return nil, errInvalidRegister
	}
	return r, nil
}

// New creates a new GPIO instance.
func New(cfg *Config) (*GPIO, error) {
	g := &GPIO{
		db:        cfg.DB,
		registers: make(map[int64]*Register),
	}
	if err := initRPIO(); err != nil {
		return nil, err
	}
	if err := g.db.AutoMigrate(&register_db.Register{}); err != nil {
		return nil, err
	}
	registers := []*register_db.Register{}
	if err := g.db.Find(&registers).Error; err != nil {
		return nil, err
	}
	for _, r := range registers {
		g.registers[r.ID] = NewRegister(r)
	}
	return g, nil
}

func (g *GPIO) ID() string {
	return ProviderID
}

func (g *GPIO) Name() string {
	return "GPIO Shift Register"
}

func (g *GPIO) Close() {
	closeRPIO()
}

func (g *GPIO) Groups() []*registry.Group {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	groups := []*registry.Group{}
	for _, r := range g.registers {
		groups = append(groups, &registry.Group{
			ID:   fmt.Sprint(r.ID),
			Name: r.Name,
		})
	}
	return groups
}

func (g *GPIO) Lamps() []*registry.Lamp {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	lamps := []*registry.Lamp{}
	for _, r := range g.registers {
		for i, s := range r.channels {
			lamps = append(lamps, &registry.Lamp{
				ID:      fmt.Sprint(i),
				Name:    fmt.Sprintf("Channel %d", i+1),
				GroupID: fmt.Sprint(r.Register.ID),
				State:   s,
			})
		}
	}
	return lamps
}

func (g *GPIO) Apply(changes []*registry.Change) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	dirtyRegisters := make(map[*Register]interface{})
	for _, c := range changes {
		r, err := g.findRegister(c.GroupID)
		if err != nil {
			return err
		}
		v, err := strconv.ParseInt(c.LampID, 10, 64)
		if err != nil {
			return err
		}
		if v < 0 || v >= r.Width {
			return errInvalidChannel
		}
		r.channels[v] = c.State
		dirtyRegisters[r] = nil
	}
	for r := range dirtyRegisters {
		r.Cycle()
	}
	return nil
}
