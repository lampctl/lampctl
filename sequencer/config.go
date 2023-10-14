package sequencer

import (
	"github.com/lampctl/lampctl/registry"
)

// Config provides the configuration for the sequencer.
type Config struct {
	Registry *registry.Registry
}
