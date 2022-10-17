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
	}

	fmt.Println(opcode, inst, imm, rs1, rd)
	return c
}

func main() {
	var c ZVMContext

	c = ZVMReset(c)

	c = ZVMTick(c)
}
