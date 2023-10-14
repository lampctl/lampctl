package server

import (
	"github.com/lampctl/lampctl/registry"
	"github.com/lampctl/lampctl/sequencer"
)

// Config provides the configuration for the web server.
type Config struct {
	Addr      string
	Debug     bool
	Registry  *registry.Registry
	Sequencer *sequencer.Sequencer
}
