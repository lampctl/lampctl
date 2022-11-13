package db

// Bridge provides database storage for Hue connection information.
type Bridge struct {
	ID       int64  `gorm:"primaryKey"`
	Host     string `gorm:"not null"`
	Username string `gorm:"not null"`
}
