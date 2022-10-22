package main

import "fmt"

type ZVMContext struct {
	r  []uint32
	m  []uint8
	pc uint32
}

func ZVMReset(c ZVMContext) ZVMContext {
	c.m = make([]uint8, 4294967296, 4294967296)
	c.r = make([]uint32, 32, 32)
	c.pc = 0x08000000
	return c
}

func ZVMTick(c ZVMContext) ZVMContext {
	var opcode uint32
	var inst uint32
	var imm uint32
	var rs1 uint32
	var rd uint32
	var funct3 uint32

	c.r[0] = 0

	inst = uint32(c.m[c.pc]) | uint32(c.m[c.pc+1])<<8 | uint32(c.m[c.pc+2])<<16 | uint32(c.m[c.pc+3])<<24
	opcode = inst & 0x7f

	if opcode == 0b0110111 {
		// LUI
		rd = (inst >> 7) & 0x1f
		c.r[rd] = inst & 0xfffff000
	} else if opcode == 0b0010111 {
		// AUIPC
		rd = (inst >> 7) & 0x1f
		c.r[rd] = (inst & 0xfffff000) + c.pc
	} else if opcode == 0b1101111 {
		// JAL
		imm = ((inst & 0x80000000) >> 11) | (inst & 0xff000) | ((inst >> 9) & 0x800) | ((inst >> 20) & 0x7fe)

		if (imm & 0x100000) == 0x100000 {
			imm = imm | 0xffe00000
		}

		rd = (inst >> 7) & 0x1f
		c.r[rd] = c.pc + 4
		c.pc = c.pc + imm
	} else if opcode == 0b1100111 {
		// JALR
		rd = (inst >> 7) & 0x1f
		rs1 = (inst >> 15) & 0x1f
		imm = inst >> 20

		if (imm & 0x800) == 0x800 {
			imm = imm | 0xfffff000
		}

		c.r[rd] = c.pc + 4
		c.pc = (imm + c.r[rs1]) & 0xfffffffe

	} else if opcode == 0b1100011 {
		imm = ((inst & 0x80000000) >> 19) | ((inst & 0x80) << 4) | ((inst >> 20) & 0x7e0) | ((inst >> 7) & 0x1e)

		funct3 = (inst >> 12) & 0x7
		if funct3 == 0b000 {

		}

	}

	fmt.Println(opcode, inst, imm, rs1, rd)
	return c
}

func main() {
	var c ZVMContext

	c = ZVMReset(c)

	c = ZVMTick(c)
}
