//go:build !unix

package gpio

import (
	register_db "github.com/lampctl/lampctl/gpio/db"
)

func initRPIO() error {
	return nil
}

func closeRPIO() {}

type Register struct {
	*register_db.Register
	channels []bool
}

func NewRegister(register *register_db.Register) *Register {
	return &Register{
		Register: register,
		channels: make([]bool, register.Width),
	}
}

func (r *Register) Cycle() {}
