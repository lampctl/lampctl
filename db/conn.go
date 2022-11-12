package db

import (
	"path"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Conn maintains a connection to the database.
type Conn struct {
	*gorm.DB
}

// New attempts to connect to or create the database.
func New(cfg *Config) (*Conn, error) {
	d, err := gorm.Open(
		sqlite.Open(path.Join(cfg.Path, "db.sqlite3")),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}
	return &Conn{
		DB: d,
	}, nil
}

// Transaction runs the provided function in a transaction.
func (c *Conn) Transaction(fn func(*Conn) error) error {
	return c.DB.Transaction(func(tx *gorm.DB) error {
		return fn(&Conn{DB: tx})
	})
}

// Close closes the database connection.
func (c *Conn) Close() {
	db, _ := c.DB.DB()
	db.Close()
}
