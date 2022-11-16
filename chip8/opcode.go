package chip8

import "fmt"

const (
	OP_SYS              uint = 0x0000_F000
	OP_CLS              uint = 0x00E0_FFFF
	OP_RET              uint = 0x00EE_FFFF
	OP_JP_ADDR          uint = 0x1000_F000
	OP_CALL_ADDR        uint = 0x2000_F000
	OP_SE_VX_BYTE       uint = 0x3000_F000
	OP_SNE_VX_BYTE      uint = 0x4000_F000
	OP_LD_VX_BYTE       uint = 0x6000_F000
	OP_ADD_VX_BYTE      uint = 0x7000_F000
	OP_LD_VX_VY         uint = 0x8000_F00F
	OP_OR_VX_VY         uint = 0x8001_F00F
	OP_AND_VX_VY        uint = 0x8002_F00F
	OP_XOR_VX_VY        uint = 0x8003_F00F
	OP_ADD_VX_VY        uint = 0x8004_F00F
	OP_SUB_VX_VY        uint = 0x8005_F00F
	OP_SHR_VX           uint = 0x8006_F00F
	OP_SUBN_VX_VY       uint = 0x8007_F00F
	OP_SHL_VX           uint = 0x800E_F00F
	OP_LD_I_ADDR        uint = 0xA000_F000
	OP_RND_VX_BYTE      uint = 0xC000_F000
	OP_DRW_VX_VY_NIBBLE uint = 0xD000_F000
	OP_SKP_VX           uint = 0xE09E_F0FF
	OP_SKNP_VX          uint = 0xE0A1_F0FF
	OP_LD_VX_DT         uint = 0xF007_F0FF
	OP_LD_DT_VX         uint = 0xF015_F0FF
	OP_ADD_I_VX         uint = 0xF01E_F0FF
	OP_LD_F_VX          uint = 0xF029_F0FF
	OP_LD_B_VX          uint = 0xF033_F0FF
	OP_LD_VX_I          uint = 0xF055_F0FF
	OP_LD_I_VX          uint = 0xF065_F0FF
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

func (o opcode) vx() uint8 {
	return uint8(o & 0x0F00 >> 8)
}

func (o opcode) vy() uint8 {
	return uint8(o & 0x00F0 >> 4)
}

func (o opcode) byte() uint8 {
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
		return fmt.Sprintf("JP_ADDR %x", op.addr())
	} else if op.equal(OP_CALL_ADDR) {
		return fmt.Sprintf("CALL_ADDR %x", op.addr())
	} else if op.equal(OP_SE_VX_BYTE) {
		return fmt.Sprintf("SE_VX_BYTE %x %x", op.vx(), op.byte())
	} else if op.equal(OP_SNE_VX_BYTE) {
		return fmt.Sprintf("SNE_VX_BYTE %x %x", op.vx(), op.byte())
	} else if op.equal(OP_LD_VX_BYTE) {
		return fmt.Sprintf("LD_VX_BYTE %x %x", op.vx(), op.byte())
	} else if op.equal(OP_ADD_VX_BYTE) {
		return fmt.Sprintf("ADD_VX_BYTE %x %x", op.vx(), op.byte())
	} else if op.equal(OP_LD_VX_VY) {
		return fmt.Sprintf("LD_VX_VY %x %x", op.vx(), op.vy())
	} else if op.equal(OP_OR_VX_VY) {
		return fmt.Sprintf("OR_VX_VY %x %x", op.vx(), op.vy())
	} else if op.equal(OP_AND_VX_VY) {
		return fmt.Sprintf("AND_VX_VY %x %x", op.vx(), op.vy())
	} else if op.equal(OP_XOR_VX_VY) {
		return fmt.Sprintf("XOR_VX_VY %x %x", op.vx(), op.vy())
	} else if op.equal(OP_ADD_VX_VY) {
		return fmt.Sprintf("ADD_VX_VY %x %x", op.vx(), op.vy())
	} else if op.equal(OP_SUB_VX_VY) {
		return fmt.Sprintf("SUB_VX_VY %x %x", op.vx(), op.vy())
	} else if op.equal(OP_SHR_VX) {
		return fmt.Sprintf("SHR_VX %x %x", op.vx(), op.vy())
	} else if op.equal(OP_SUBN_VX_VY) {
		return fmt.Sprintf("SUBN_VX_VY %x %x", op.vx(), op.vy())
	} else if op.equal(OP_SHL_VX) {
		return fmt.Sprintf("SHL_VX %x %x", op.vx(), op.vy())
	} else if op.equal(OP_LD_I_ADDR) {
		return fmt.Sprintf("LD_I_ADD %x", op.addr())
	} else if op.equal(OP_RND_VX_BYTE) {
		return fmt.Sprintf("RND_VX_BYTE %x %x", op.vx(), op.byte())
	} else if op.equal(OP_DRW_VX_VY_NIBBLE) {
		return fmt.Sprintf("DRW_VX_VY_NIBBLE %x %x %x", op.vx(), op.vy(), op.nibble())
	} else if op.equal(OP_SKP_VX) {
		return fmt.Sprintf("SKP_VX %x", op.vx())
	} else if op.equal(OP_SKNP_VX) {
		return fmt.Sprintf("SKNP_VX %x", op.vx())
	} else if op.equal(OP_LD_VX_DT) {
		return fmt.Sprintf("LD_VX_DT %x", op.vx())
	} else if op.equal(OP_LD_DT_VX) {
		return fmt.Sprintf("LD_DT_VX %x", op.vx())
	} else if op.equal(OP_ADD_I_VX) {
		return fmt.Sprintf("ADD_I_VX %x", op.vx())
	} else if op.equal(OP_LD_F_VX) {
		return fmt.Sprintf("LD_F_VX %x", op.vx())
	} else if op.equal(OP_LD_B_VX) {
		return fmt.Sprintf("LD_B_VX %x", op.vx())
	} else if op.equal(OP_LD_VX_I) {
		return fmt.Sprintf("LD_VX_I %x", op.vx())
	} else if op.equal(OP_LD_I_VX) {
		return fmt.Sprintf("LD_I_VX %x", op.vx())
	} else {
		panic(fmt.Errorf("unrecognised opcode %x", uint(op)))
	}
}
