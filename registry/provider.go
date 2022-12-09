package registry

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidGroup = errors.New("invalid group specified")
	ErrInvalidLamp  = errors.New("invalid lamp specified")
)

// Group provides a logical grouping for lamps in a provider.
type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Lamp provides information about a specific lamp that can be controlled.
type Lamp struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	GroupID string `json:"group_id"`
	State   bool   `json:"state"`
}

// Change represents a request to change the state of a lamp.
type Change struct {
	GroupID    string  `json:"group_id"`
	LampID     string  `json:"lamp_id"`
	State      bool    `json:"state"`
	Duration   int64   `json:"duration"`
	Brightness float64 `json:"brightness"`
}

// Provider represents a group of lamps. The interface provides methods for
// initializing, enumerating, and controlling them.
type Provider interface {

	// ID returns a machine-friendly name for the provider.
	ID() string

	// Name returns a human-friendly name for the provider.
	Name() string

	// Init initializes the provider, allowing it to register any API routes.
	Init(api *gin.RouterGroup) error

	// Close frees all resources associated with the provider.
	Close()

	// Groups returns a list of groups in the provider.
	Groups() []*Group

	// Lamps returns a list of all lamps managed by the provider.
	Lamps() []*Lamp

	// Apply applies a list of state changes to the lamps in the provider.
	Apply(changes []*Change) error
}
