package chip8

import (
	"fmt"
	"io"
)

type chip8 struct {
	r [16]byte
	i uint16

	stack [64]byte
	sp    byte
	pc    uint16

	memory [4096]byte
	fb     [256]byte

	timer      byte
	soundTimer byte
}

func NewChip8(rom io.Reader) (*chip8, error) {
	c := new(chip8)
	_, err := rom.Read(c.memory[0x0200:])
	if err != nil {
		return nil, fmt.Errorf("error loading rom: %v", err)
	}

	c.pc = 0x0200

	return c, nil
}

func (c *chip8) Tick() error {

	op := c.memory[c.pc : c.pc+2]
	op1 := (op[0] >> 4) & 0x0F
	op2 := (op[0] >> 0) & 0x0F
	op3 := (op[1] >> 4) & 0x0F
	op4 := (op[1] >> 0) & 0x0F

	fmt.Printf("%x - %x%x%x%x\n", op, op1, op2, op3, op4)

	c.pc += 2

	return nil
}
