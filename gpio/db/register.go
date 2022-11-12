package db

// Register provides database storage for register info.
type Register struct {
	ID       int64  `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	DataPin  int    `gorm:"not null" json:"data_pin"`
	LatchPin int    `gorm:"not null" json:"latch_pin"`
	ClockPin int    `gorm:"not null" json:"clock_pin"`
	Width    int64  `gorm:"not null" json:"width"`
}
