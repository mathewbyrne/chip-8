package chip8

import (
	"fmt"
	"io"

	"encoding/binary"
)

type Key uint8

type Input interface {
	State(k Key) bool
	Wait() Key
}

type chip8 struct {
	r [16]uint8
	i uint16

	stack [32]uint16
	sp    uint8
	pc    uint16

	m  [4096]byte
	fb [256]byte

	t  uint8
	dt uint8

	input Input
}

func NewChip8(rom io.Reader, i Input) (*chip8, error) {
	c := new(chip8)
	_, err := rom.Read(c.m[0x0200:])
	if err != nil {
		return nil, fmt.Errorf("error loading rom: %v", err)
	}

	c.pc = 0x0200
	c.input = i

	return c, nil
}

const (
	OP_SYS           = 0x0000_F000
	OP_CLS           = 0x00E0_FFFF
	OP_RET           = 0x00EE_FFFF
	OP_JP_ADDR       = 0x1000_F000
	OP_CALL_ADDR     = 0x2000_F000
	OP_SE_VX_VAL     = 0x3000_F000
	OP_LD_VX_VAL     = 0x6000_F000
	OP_LD_I_ADDR     = 0xA000_F000
	OP_DRW_VX_VY_SPR = 0xD000_F000
	OP_SKP_VX        = 0xE09E_F0FF
	OP_SKNP_VX       = 0xE0A1_F0FF
	OP_LD_VX_DT      = 0xF007_F0FF
	OP_LD_DT_VX      = 0xF015_F0FF
)

type opcode uint16

func (o opcode) equal(code uint) bool {
	val := uint16(code >> 16)
	mask := uint16(code)
	return (uint16(o) & mask) == val
}

func (o opcode) addr() uint16 {
	return uint16(o) & 0x0FFF
}

func (o opcode) r1() uint8 {
	return uint8(o & 0x0F00 >> 8)
}

func (o opcode) r2() uint8 {
	return uint8(o & 0x00F0 >> 4)
}

func (o opcode) val() uint8 {
	return uint8(o & 0x00FF)
}

func (c *chip8) next() {
	c.pc += 2
}

func (c *chip8) Tick() error {

	op := opcode(binary.BigEndian.Uint16(c.m[c.pc : c.pc+2]))
	c.next()
	fmt.Printf("<- %04x\n", op)

	if op.equal(OP_CLS) {
		for i := range c.fb {
			c.fb[i] = 0
		}
		fmt.Println("clear screen")
	} else if op.equal(OP_RET) {
		c.pc = c.stack[c.sp]
		c.sp--

		fmt.Println("return")
	} else if op.equal(OP_SYS) {
		fmt.Println("system")
	} else if op.equal(OP_JP_ADDR) {
		c.pc = op.addr()

		fmt.Printf("jump %x\n", op.addr())
	} else if op.equal(OP_CALL_ADDR) {
		c.sp++
		c.stack[c.sp] = c.pc
		c.pc = op.addr()

		fmt.Printf("call %x\n", c.pc)
	} else if op.equal(OP_SE_VX_VAL) {
		if c.r[op.r1()] == op.val() {
			c.next()
		}

		fmt.Printf("skip\n")
	} else if op.equal(OP_LD_VX_VAL) {
		c.r[op.r1()] = op.val()

		fmt.Printf("load %x %x\n", op.r1(), op.val())
	} else if op.equal(OP_LD_I_ADDR) {
		c.i = op.addr()

		fmt.Printf("load I %x\n", op.addr())
	} else if op.equal(OP_DRW_VX_VY_SPR) {
		fmt.Printf("draw\n")
	} else if op.equal(OP_SKP_VX) {
		if c.input.State(Key(c.r[op.r1()])) {
			c.next()
		}

		fmt.Printf("skip if key %x\n", c.r[op.r1()])
	} else if op.equal(OP_SKNP_VX) {
		if !c.input.State(Key(c.r[op.r1()])) {
			c.next()
		}

		fmt.Printf("skip if not key %x\n", c.r[op.r1()])
	} else if op.equal(OP_LD_VX_DT) {
		c.r[op.r1()] = c.dt

		fmt.Printf("load vx to dt\n")
	} else if op.equal(OP_LD_DT_VX) {
		c.dt = c.r[op.r1()]

		fmt.Printf("load dt to vx\n")
	} else {
		panic(fmt.Errorf("unrecognised opcode %x\n", op))
	}

	fmt.Printf("-> PC: %04x I: %04x\n", c.pc, c.i)
	fmt.Printf("\t%02x %02x %02x %02x\n", c.r[0x0], c.r[0x1], c.r[0x2], c.r[0x3])
	fmt.Printf("\t%02x %02x %02x %02x\n", c.r[0x4], c.r[0x5], c.r[0x6], c.r[0x7])
	fmt.Printf("\t%02x %02x %02x %02x\n", c.r[0x8], c.r[0x9], c.r[0xA], c.r[0xB])
	fmt.Printf("\t%02x %02x %02x %02x\n", c.r[0xC], c.r[0xD], c.r[0xE], c.r[0xF])

	return nil
}
