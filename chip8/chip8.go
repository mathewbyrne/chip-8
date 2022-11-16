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

type Chip8 struct {
	r [16]uint8
	i uint16

	stack [32]uint16
	sp    uint8
	pc    uint16

	m  [4096]byte
	fb FrameBuffer

	t  uint8
	dt uint8

	input Input
}

func NewChip8(rom io.Reader, i Input) (*Chip8, error) {
	c := new(Chip8)
	_, err := rom.Read(c.m[0x0200:])
	if err != nil {
		return nil, fmt.Errorf("error loading rom: %v", err)
	}

	c.pc = 0x0200
	c.input = i

	return c, nil
}

func (c *Chip8) next() {
	c.pc += 2
}

func (c *Chip8) Tick() {
	if c.t > 0 {
		c.t -= 1
	}
	if c.dt > 0 {
		c.dt -= 1
	}
}

func (c *Chip8) Cycle() {

	op := opcode(binary.BigEndian.Uint16(c.m[c.pc : c.pc+2]))
	fmt.Printf("%04X [%04X] %s\n", c.pc, uint16(op), op)

	c.next()

	if op.equal(OP_CLS) {
		for i := range c.fb {
			c.fb[i] = 0
		}
	} else if op.equal(OP_RET) {
		c.pc = c.stack[c.sp]
		c.sp--
	} else if op.equal(OP_SYS) {
		// noop
	} else if op.equal(OP_JP_ADDR) {
		c.pc = op.addr()
	} else if op.equal(OP_CALL_ADDR) {
		c.sp++
		c.stack[c.sp] = c.pc
		c.pc = op.addr()
	} else if op.equal(OP_SE_VX_VAL) {
		if c.r[op.r1()] == op.val() {
			c.next()
		}
	} else if op.equal(OP_SNE_VX_VAL) {
		if c.r[op.r1()] != op.val() {
			c.next()
		}
	} else if op.equal(OP_LD_VX_VAL) {
		c.r[op.r1()] = op.val()
	} else if op.equal(OP_ADD_VX_VAL) {
		c.r[op.r1()] += op.val()
	} else if op.equal(OP_LD_I_ADDR) {
		c.i = op.addr()
	} else if op.equal(OP_DRW_VX_VY_SPR) {

		x := c.r[op.r1()]
		y := c.r[op.r2()]
		n := op.nibble()

		c.r[0xF] = 0
		if c.fb.draw(c.m[c.i:c.i+uint16(n)], x, y) {
			c.r[0xF] = 1
		}

		fmt.Println(&c.fb)
	} else if op.equal(OP_SKP_VX) {
		if c.input.State(Key(c.r[op.r1()])) {
			c.next()
		}
	} else if op.equal(OP_SKNP_VX) {
		if !c.input.State(Key(c.r[op.r1()])) {
			c.next()
		}
	} else if op.equal(OP_LD_VX_DT) {
		c.r[op.r1()] = c.dt
	} else if op.equal(OP_LD_DT_VX) {
		c.dt = c.r[op.r1()]
	} else if op.equal(OP_ADD_I_VX) {
		c.i += uint16(c.r[op.r1()])
	} else {
		panic(fmt.Errorf("unrecognised opcode %x", op))
	}

	fmt.Printf("%s\n", c)
}

func (c *Chip8) String() string {
	return fmt.Sprintf(`pc:%04x sp: %02d i:%04x t: %02d dt: %02d
0:%02x 1:%02x 2:%02x 3:%02x 4:%02x 5:%02x 6:%02x 7:%02x 8:%02x 9:%02x A:%02x B:%02x C:%02x D:%02x E:%02x F:%02x
`,
		c.pc, c.sp, c.i, c.t, c.dt,
		c.r[0x0], c.r[0x1], c.r[0x2], c.r[0x3],
		c.r[0x4], c.r[0x5], c.r[0x6], c.r[0x7],
		c.r[0x8], c.r[0x9], c.r[0xA], c.r[0xB],
		c.r[0xC], c.r[0xD], c.r[0xE], c.r[0xF],
	)
}

func (c *Chip8) FrameBuffer() []byte {
	return c.fb[:]
}
