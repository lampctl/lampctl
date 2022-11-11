package db

// Register provides database storage for register info.
type Register struct {
	ID       int64  `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	DataPin  int    `gorm:"not null"`
	LatchPin int    `gorm:"not null"`
	ClockPin int    `gorm:"not null"`
	Width    int64  `gorm:"not null"`
}
