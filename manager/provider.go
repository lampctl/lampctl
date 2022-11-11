package manager

import (
	"github.com/lampctl/lampctl/db"
)

// Provider represents a group of lamps. The interface provides methods for
// initializing, enumerating, and controlling them.
type Provider interface {

	// ID returns a machine-friendly name for the provider.
	ID() string

	// Name returns a human-friendly name for the provider.
	Name() string

	// Init initializes the provider, including database models.
	Init(*db.Conn) error

	// Free closes all resources associated with the provider.
	Free()

	// GetLamps returns a list of all lamps managed by the provider.
	GetLamps() []Lamp

	// Apply causes any state changes to immediately take effect.
	Apply() error
}
