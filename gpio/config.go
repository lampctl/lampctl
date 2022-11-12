package gpio

import (
	"github.com/lampctl/lampctl/db"
)

// Config provides the configuration for working with registers.
type Config struct {
	DB *db.Conn
}
