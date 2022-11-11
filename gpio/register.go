package gpio

import (
	register_db "github.com/lampctl/lampctl/gpio/db"
	"github.com/stianeikeland/go-rpio/v4"
)

// Register represents an individual shift register chain.
type Register struct {
	*register_db.Register
	channels []bool
}

// NewRegister creates and initializes a new Register instance.
func NewRegister(register *register_db.Register) *Register {
	r := &Register{
		Register: register,
		channels: make([]bool, register.Width),
	}
	rpio.Pin(r.DataPin).Output()
	rpio.Pin(r.LatchPin).Output()
	rpio.Pin(r.ClockPin).Output()
	return r
}

// Cycle applies state changes to the channels in the register.
func (r *Register) Cycle() {

	// Put latch down to start sending data
	rpio.Pin(r.ClockPin).Low()
	rpio.Pin(r.LatchPin).Low()
	rpio.Pin(r.ClockPin).High()

	// Write the channel values, one at a time
	for _, s := range r.channels {
		rpio.Pin(r.ClockPin).Low()
		if s {
			rpio.Pin(r.DataPin).High()
		} else {
			rpio.Pin(r.DataPin).Low()
		}
		rpio.Pin(r.ClockPin).High()
	}

	// Put latch up to store data in register
	rpio.Pin(r.ClockPin).Low()
	rpio.Pin(r.LatchPin).High()
	rpio.Pin(r.ClockPin).High()
}
