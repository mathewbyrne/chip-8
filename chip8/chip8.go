package chip8

import (
	"fmt"
	"io"
	"math/rand"
	"sync/atomic"
	"time"

	"encoding/binary"
)

type Key uint8

type Input interface {
	State(k Key) bool
	Wait() Key
}

type Chip8 struct {
	// registers
	r [16]uint8
	i uint16

	// stack
	stack [32]uint16
	sp    uint8
	pc    uint16

	// buffers
	m  [4096]byte
	fb FrameBuffer

	// timers
	dt uint8
	st uint8

	rnd rand.Rand

	KeyMap          atomic.Uint32
	lastCycleKeyMap uint32
}

func New(rom io.Reader) (*Chip8, error) {
	c := &Chip8{}

	// font data occupies the first 80 bytes
	copy(c.m[:], fontData[:])

	_, err := rom.Read(c.m[0x0200:])
	if err != nil {
		return nil, fmt.Errorf("error loading rom: %v", err)
	}

	c.pc = 0x0200
	c.rnd = *rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	return c, nil
}

func (c *Chip8) carry(val bool) {
	if val {
		c.r[0xF] = 1
	} else {
		c.r[0xF] = 0
	}
}

func (c *Chip8) Tick() {
	if c.dt > 0 {
		c.dt -= 1
	}
	if c.st > 0 {
		c.st -= 1
	}
}

func (c *Chip8) Cycle() (opcode, uint16) {
	pc := c.pc
	op := opcode(binary.BigEndian.Uint16(c.m[c.pc : c.pc+2]))

	var pcInc uint16 = OP_WIDTH

	if op.equal(OP_CLS) {
		c.fb.clear()
	} else if op.equal(OP_RET) {
		c.pc = c.stack[c.sp]
		c.sp--
		pcInc = 0
	} else if op.equal(OP_SYS) {
		// noop
	} else if op.equal(OP_JP_ADDR) {
		c.pc = op.addr()
		pcInc = 0
	} else if op.equal(OP_CALL_ADDR) {
		c.sp++
		c.stack[c.sp] = c.pc + OP_WIDTH
		c.pc = op.addr()
		pcInc = 0
	} else if op.equal(OP_SE_VX_BYTE) {
		if c.r[op.vx()] == op.byte() {
			pcInc += OP_WIDTH
		}
	} else if op.equal(OP_SNE_VX_BYTE) {
		if c.r[op.vx()] != op.byte() {
			pcInc += OP_WIDTH
		}
	} else if op.equal(OP_SE_VX_VY) {
		if c.r[op.vx()] == c.r[op.vy()] {
			pcInc += OP_WIDTH
		}
	} else if op.equal(OP_LD_VX_BYTE) {
		c.r[op.vx()] = op.byte()
	} else if op.equal(OP_ADD_VX_BYTE) {
		c.r[op.vx()] += op.byte()
	} else if op.equal(OP_LD_VX_VY) {
		c.opLdVxVy(op.vx(), op.vy())
	} else if op.equal(OP_OR_VX_VY) {
		c.opOrVxVy(op.vx(), op.vy())
	} else if op.equal(OP_AND_VX_VY) {
		c.opAndVxVy(op.vx(), op.vy())
	} else if op.equal(OP_XOR_VX_VY) {
		c.opXorVxVy(op.vx(), op.vy())
	} else if op.equal(OP_ADD_VX_VY) {
		c.opAddVxVy(op.vx(), op.vy())
	} else if op.equal(OP_SUB_VX_VY) {
		c.opSubVxVy(op.vx(), op.vy())
	} else if op.equal(OP_SHR_VX_VY) {
		c.opShrVxVy(op.vx(), op.vy())
	} else if op.equal(OP_SUBN_VX_VY) {
		c.opSubnVxVy(op.vx(), op.vy())
	} else if op.equal(OP_SHL_VX_VY) {
		c.opShlVxVy(op.vx(), op.vy())
	} else if op.equal(OP_SNE_VX_VY) {
		if c.r[op.vx()] != c.r[op.vy()] {
			pcInc += OP_WIDTH
		}
	} else if op.equal(OP_LD_I_ADDR) {
		c.i = op.addr()
	} else if op.equal(OP_JP_V0_ADDR) {
		c.pc = op.addr() + uint16(c.r[0])
		pcInc = 0
	} else if op.equal(OP_RND_VX_BYTE) {
		c.r[op.vx()] = uint8(c.rnd.Uint32()) & op.byte()
	} else if op.equal(OP_DRW_VX_VY_NIBBLE) {
		carry := c.fb.draw(c.m[c.i:c.i+uint16(op.nibble())], c.r[op.vx()], c.r[op.vy()])
		c.carry(carry)
	} else if op.equal(OP_SKP_VX) {
		pcInc = c.opSkpVx(op.vx())
	} else if op.equal(OP_SKNP_VX) {
		pcInc = c.opSknpVx(op.vx())
	} else if op.equal(OP_LD_VX_DT) {
		c.r[op.vx()] = c.dt
	} else if op.equal(OP_LD_VX_K) {
		pcInc = c.opLdVxK(op.vx())
	} else if op.equal(OP_LD_DT_VX) {
		c.dt = c.r[op.vx()]
	} else if op.equal(OP_LD_ST_VX) {
		c.st = c.r[op.vx()]
	} else if op.equal(OP_ADD_I_VX) {
		c.i += uint16(c.r[op.vx()])
	} else if op.equal(OP_LD_F_VX) {
		c.opLdFVx(op.vx())
	} else if op.equal(OP_LD_B_VX) {
		c.opLdBVx(op.vx())
	} else if op.equal(OP_LD_VX_I) {
		c.opLdVxI(op.vx())
	} else if op.equal(OP_LD_I_VX) {
		c.opLdIVx(op.vx())
	} else {
		panic(fmt.Errorf("unrecognised opcode %x", op))
	}

	c.pc += pcInc
	c.lastCycleKeyMap = c.KeyMap.Load()

	return op, pc
}

func (c *Chip8) opLdVxVy(vx, vy uint8) {
	c.r[vx] = c.r[vy]
}

func (c *Chip8) opOrVxVy(vx, vy uint8) {
	defer c.carry(false)
	c.r[vx] |= c.r[vy]
}

func (c *Chip8) opAndVxVy(vx, vy uint8) {
	defer c.carry(false)
	c.r[vx] &= c.r[vy]
}

func (c *Chip8) opXorVxVy(vx, vy uint8) {
	defer c.carry(false)
	c.r[vx] ^= c.r[vy]
}

func (c *Chip8) opAddVxVy(vx, vy uint8) {
	defer c.carry(uint16(c.r[vx])+uint16(c.r[vy]) > 0xFF)
	c.r[vx] += c.r[vy]
}

func (c *Chip8) opSubVxVy(vx, vy uint8) {
	defer c.carry(c.r[vx] > c.r[vy])
	c.r[vx] -= c.r[vy]
}

func (c *Chip8) opShrVxVy(vx, vy uint8) {
	defer c.carry(c.r[vy]&0x01 == 0x01)
	c.r[vx] = c.r[vy] >> 1
}

func (c *Chip8) opSubnVxVy(vx, vy uint8) {
	defer c.carry(c.r[vy] > c.r[vx])
	c.r[vx] = c.r[vy] - c.r[vx]
}

func (c *Chip8) opShlVxVy(vx, vy uint8) {
	defer c.carry(c.r[vy]&0x80 == 0x80)
	c.r[vx] = c.r[vy] << 1
}

func (c *Chip8) opSkpVx(vx uint8) uint16 {
	keymap := c.KeyMap.Load()
	key := uint32(1) << c.r[vx]
	if keymap&key == key {
		return OP_WIDTH * 2
	}
	return OP_WIDTH
}

func (c *Chip8) opSknpVx(vx uint8) uint16 {
	keymap := c.KeyMap.Load()
	key := uint32(1) << c.r[vx]
	if keymap&key != key {
		return OP_WIDTH * 2
	}
	return OP_WIDTH
}

func (c *Chip8) opLdVxI(vx uint8) {
	copy(c.m[c.i:], c.r[0:vx+1])
	c.i += uint16(vx) + 1
}

func (c *Chip8) opLdIVx(vx uint8) {
	copy(c.r[0:vx+1], c.m[c.i:])
	c.i += uint16(vx) + 1
}

func (c *Chip8) opLdVxK(vx uint8) uint16 {
	r := c.lastCycleKeyMap &^ c.KeyMap.Load()
	if r > 0 {
		for i := uint8(0); i < 16; i++ {
			if (1<<i)&r == 1<<i {
				c.r[vx] = i
			}
		}
		return OP_WIDTH
	}
	return 0
}

func (c *Chip8) opLdFVx(vx uint8) {
	c.i = uint16(c.r[vx]) * 5
}

func (c *Chip8) opLdBVx(vx uint8) {
	c.m[c.i+0] = (c.r[vx] / 100) % 10
	c.m[c.i+1] = (c.r[vx] / 10) % 10
	c.m[c.i+2] = c.r[vx] % 10
}

func (c *Chip8) String() string {
	return fmt.Sprintf(`pc:%04x sp: %02d i:%04x dt:%02d st:%02d km:%016b
0:%02x 1:%02x 2:%02x 3:%02x 4:%02x 5:%02x 6:%02x 7:%02x 8:%02x 9:%02x A:%02x B:%02x C:%02x D:%02x E:%02x F:%02x
`,
		c.pc, c.sp, c.i, c.dt, c.st, c.KeyMap.Load(),
		c.r[0x0], c.r[0x1], c.r[0x2], c.r[0x3],
		c.r[0x4], c.r[0x5], c.r[0x6], c.r[0x7],
		c.r[0x8], c.r[0x9], c.r[0xA], c.r[0xB],
		c.r[0xC], c.r[0xD], c.r[0xE], c.r[0xF],
	)
}

func (c *Chip8) FrameBuffer() *FrameBuffer {
	return &c.fb
}
