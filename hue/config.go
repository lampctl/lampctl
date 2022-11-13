package hue

import (
	"github.com/lampctl/lampctl/db"
)

// Config provides the configuration for the Hue provider.
type Config struct {
	DB *db.Conn
}
