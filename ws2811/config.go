package ws2811

import (
	"github.com/lampctl/lampctl/db"
)

// Config provides the configuration for working with ws2811 channels.
type Config struct {
	DB *db.Conn
}
