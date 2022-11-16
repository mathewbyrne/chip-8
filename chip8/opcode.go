package chip8

import "fmt"

const (
	OP_SYS           = 0x0000_F000
	OP_CLS           = 0x00E0_FFFF
	OP_RET           = 0x00EE_FFFF
	OP_JP_ADDR       = 0x1000_F000
	OP_CALL_ADDR     = 0x2000_F000
	OP_SE_VX_VAL     = 0x3000_F000
	OP_LD_VX_VAL     = 0x6000_F000
	OP_ADD_VX_VAL    = 0x7000_F000
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

func (o opcode) nibble() uint8 {
	return uint8(o & 0x000F)
}

func (op opcode) String() string {
	if op.equal(OP_CLS) {
		return "CLS"
	} else if op.equal(OP_RET) {
		return "RET"
	} else if op.equal(OP_SYS) {
		return "SYS"
	} else if op.equal(OP_JP_ADDR) {
		return fmt.Sprintf("JP %x", op.addr())
	} else if op.equal(OP_CALL_ADDR) {
		return fmt.Sprintf("CALL %x", op.addr())
	} else if op.equal(OP_SE_VX_VAL) {
		return fmt.Sprintf("SE %x %x", op.r1(), op.val())
	} else if op.equal(OP_LD_VX_VAL) {
		return fmt.Sprintf("LD %x %x", op.r1(), op.val())
	} else if op.equal(OP_ADD_VX_VAL) {
		return fmt.Sprintf("ADD %x %x", op.r1(), op.val())
	} else if op.equal(OP_LD_I_ADDR) {
		return fmt.Sprintf("LD I %x", op.addr())
	} else if op.equal(OP_DRW_VX_VY_SPR) {
		return fmt.Sprintf("DRW %x %x %x\n", op.r1(), op.r2(), op.nibble())
	} else if op.equal(OP_SKP_VX) {
		return fmt.Sprintf("SKP %x", op.r1())
	} else if op.equal(OP_SKNP_VX) {
		return fmt.Sprintf("SKNP %x", op.r1())
	} else if op.equal(OP_LD_VX_DT) {
		return fmt.Sprintf("LD %x DT", op.r1())
	} else if op.equal(OP_LD_DT_VX) {
		return fmt.Sprintf("LD DT %x", op.r1())
	} else {
		panic(fmt.Errorf("unrecognised opcode %x", uint(op)))
	}
}
