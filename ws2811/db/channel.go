package db

// Channel represents a ws2811 channel.
type Channel struct {
	ID       int64 `gorm:"primaryKey" json:"id"`
	GpioPin  int   `gorm:"not null" json:"gpio_pin"`
	LedCount int   `gorm:"not null" json:"led_count"`
}
