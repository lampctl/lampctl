package manager

// Lamp provides information about a specific lamp that can be controlled.
type Lamp interface {

	// ID returns a machine-friendly name for the lamp.
	ID() string

	// Name returns a human-friendly name for the lamp.
	Name() string

	// State returns the current state of the lamp.
	State() bool

	// SetState changes the current state of the lamp. Note that the provider's
	// Apply() method may be required to complete the change.
	SetState(value bool)
}
