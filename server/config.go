package server

import (
	"github.com/lampctl/lampctl/registry"
)

// Config provides the configuration for the web server.
type Config struct {
	Addr     string
	Debug    bool
	Registry *registry.Registry
}
