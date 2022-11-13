package db

// Bridge provides database storage for Hue connection information.
type Bridge struct {
	ID       string `gorm:"primaryKey" json:"id"`
	Host     string `gorm:"not null" json:"host"`
	Username string `gorm:"not null" json:"username"`
}
