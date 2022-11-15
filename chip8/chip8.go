package chip8

import (
	"fmt"
	"io"

	"encoding/binary"
)

type chip8 struct {
	r [16]byte
	i uint16

	stack [32]uint16
	sp    uint8
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

const (
	OP_SYS  = 0x0000
	OP_CLS  = 0x00E0
	OP_RET  = 0x00EE
	OP_JP   = 0x1000
	OP_CALL = 0x2000
	OP_LD   = 0x6000
)

func (c *chip8) Tick() error {

	op := binary.BigEndian.Uint16(c.memory[c.pc : c.pc+2])
	c.pc += 2
	fmt.Printf("<- %04x\n", op)

	if op == OP_CLS {
		fmt.Println("clear screen")
	} else if op == OP_RET {
		c.pc = c.stack[c.sp]
		c.sp--
		fmt.Println("return")
	} else if (op & 0xF000) == OP_SYS {
		fmt.Println("system")
	} else if (op & 0xF000) == OP_JP {
		fmt.Println("jump")
	} else if (op & 0xF000) == OP_CALL {
		c.sp++
		c.stack[c.sp] = c.pc
		c.pc = op & 0x0FFF

		fmt.Printf("call %x\n", c.pc)
	} else if (op & 0xF000) == OP_LD {
		fmt.Printf("load %x %x\n", op&0x0F00>>8, op&0x00FF)
	} else {
		fmt.Printf("unrecognised opcode %x\n", op)
	}

	fmt.Printf("-> %04x\n", c.pc)

	return nil
}
