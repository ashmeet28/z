package main

import "fmt"

type ZVMContext struct {
	pc uint32
	r  []uint32
	m  []uint8
	s  uint32
}

func ZVMReset(c ZVMContext) ZVMContext {
	c.m = make([]uint8, 4294967296, 4294967296)
	c.r = make([]uint32, 32, 32)
	c.pc = 0x08000000
	c.s = 2
	return c
}

func ZVMHandleECall(c ZVMContext) ZVMContext {
	c.s = 5
	return c
}

func ZVMTick(c ZVMContext) ZVMContext {
	var opcode uint32
	var inst uint32
	var imm uint32
	var rs1 uint32
	var rs2 uint32
	var rd uint32
	var funct3 uint32

	var pc uint32
	var r []uint32
	var m []uint8
	var s uint32

	pc = c.pc
	r = c.r
	m = c.m
	s = c.s

	s = 4
	r[0] = 0

	inst = uint32(m[pc]) | (uint32(m[pc+1]) << 8) | (uint32(m[pc+2]) << 16) | (uint32(m[pc+3]) << 24)

	opcode = inst & 0x7f

	if opcode == 0b0110111 {
		// LUI
		rd = (inst >> 7) & 0x1f
		imm = inst & 0xfffff000

		r[rd] = imm

		pc = pc + 4

	} else if opcode == 0b0010111 {
		// AUIPC
		rd = (inst >> 7) & 0x1f
		imm = inst & 0xfffff000

		r[rd] = imm + pc

		pc = pc + 4

	} else if opcode == 0b1101111 {
		// JAL
		imm = ((inst & 0x80000000) >> 11) | (inst & 0xff000) | ((inst >> 9) & 0x800) | ((inst >> 20) & 0x7fe)
		if (imm & 0x100000) == 0x100000 {
			imm = imm | 0xffe00000
		}
		rd = (inst >> 7) & 0x1f

		r[rd] = pc + 4
		pc = pc + imm

	} else if opcode == 0b1100111 {

		rd = (inst >> 7) & 0x1f
		rs1 = (inst >> 15) & 0x1f
		imm = inst >> 20
		if (imm & 0x800) == 0x800 {
			imm = imm | 0xfffff000
		}
		funct3 = (inst >> 12) & 0x7

		if funct3 == 0b000 {
			// JALR
			r[rd] = pc + 4
			pc = (imm + r[rs1]) & 0xfffffffe

		} else {
			s = 5
		}

	} else if opcode == 0b1100011 {

		imm = ((inst & 0x80000000) >> 19) | ((inst & 0x80) << 4) | ((inst >> 20) & 0x7e0) | ((inst >> 7) & 0x1e)
		if (imm & 0x1000) == 0x1000 {
			imm = imm | 0xffffe000
		}
		rs1 = (inst >> 15) & 0x1f
		rs2 = (inst >> 20) & 0x1f
		funct3 = (inst >> 12) & 0x7

		if funct3 == 0b000 {
			// BEQ
			if r[rs1] == r[rs2] {
				pc = imm + pc
			} else {
				pc = pc + 4
			}

		} else if funct3 == 0b001 {
			// BNE
			if r[rs1] != r[rs2] {
				pc = imm + pc
			} else {
				pc = pc + 4
			}

		} else if funct3 == 0b100 {
			// BLT
			if ((r[rs1] >> 31) == 0x0) && ((r[rs2] >> 31) == 0x0) {

				if r[rs1] < r[rs2] {
					pc = imm + pc
				} else {
					pc = pc + 4
				}

			} else if ((r[rs1] >> 31) == 0x1) && ((r[rs2] >> 31) == 0x1) {

				if r[rs1] > r[rs2] {
					pc = imm + pc
				} else {
					pc = pc + 4
				}

			} else if (r[rs1] >> 31) == 0x1 {
				pc = imm + pc
			} else {
				pc = pc + 4
			}

		} else if funct3 == 0b101 {
			// BGE
			if ((r[rs1] >> 31) == 0x0) && ((r[rs2] >> 31) == 0x0) {

				if r[rs1] >= r[rs2] {
					pc = imm + pc
				} else {
					pc = pc + 4
				}

			} else if ((r[rs1] >> 31) == 0x1) && ((r[rs2] >> 31) == 0x1) {

				if r[rs1] <= r[rs2] {
					pc = imm + pc
				} else {
					pc = pc + 4
				}

			} else if (r[rs2] >> 31) == 0x1 {
				pc = imm + pc
			} else {
				pc = pc + 4
			}

		} else if funct3 == 0b110 {
			// BLTU
			if r[rs1] < r[rs2] {
				pc = imm + pc
			} else {
				pc = pc + 4
			}

		} else if funct3 == 0b111 {
			// BGEU
			if r[rs1] >= r[rs2] {
				pc = imm + pc
			} else {
				pc = pc + 4
			}

		} else {
			s = 5
		}

	} else if opcode == 0b0000011 {

		rd = (inst >> 7) & 0x1f
		rs1 = (inst >> 15) & 0x1f
		funct3 = (inst >> 12) & 0x7
		imm = inst >> 20
		if (imm & 0x800) == 0x800 {
			imm = imm | 0xfffff000
		}

		if funct3 == 0b000 {
			// LB
			r[rd] = uint32(m[r[rs1]+imm])
			if (r[rd] & 0x80) == 0x80 {
				r[rd] = r[rd] | 0xffffff80
			}
			pc = pc + 4

		} else if funct3 == 0b001 {
			// LH
			r[rd] = uint32(m[r[rs1]+imm]) | (uint32(m[r[rs1]+imm+1]) << 8)
			if (r[rd] & 0x8000) == 0x8000 {
				r[rd] = r[rd] | 0xffff0000
			}
			pc = pc + 4

		} else if funct3 == 0b010 {
			// LW
			r[rd] = uint32(m[r[rs1]+imm]) | (uint32(m[r[rs1]+imm+1]) << 8) | (uint32(m[r[rs1]+imm+2]) << 16) | (uint32(m[r[rs1]+imm+3]) << 24)
			if (r[rd] & 0x8000) == 0x8000 {
				r[rd] = r[rd] | 0xffff0000
			}
			pc = pc + 4

		} else if funct3 == 0b100 {
			// LBU
			r[rd] = uint32(m[r[rs1]+imm])
			pc = pc + 4

		} else if funct3 == 0b101 {
			// LHU
			r[rd] = uint32(m[r[rs1]+imm]) | (uint32(m[r[rs1]+imm+1]) << 8)
			pc = pc + 4

		} else {
			s = 5
		}

	} else if opcode == 0b0100011 {
		rs1 = (inst >> 15) & 0x1f
		rs2 = (inst >> 20) & 0x1f
		funct3 = (inst >> 12) & 0x7
		imm = ((inst >> 7) & 0x1f) | ((inst >> 20) & 0xfe0)

		if (imm & 0x800) == 0x800 {
			imm = imm | 0xfffff000
		}

		if funct3 == 0b000 {
			// SB
			m[r[rs1]+imm] = uint8(r[rs2] & 0xff)
            pc = pc + 4

		} else if funct3 == 0b001 {
			// SH
			m[r[rs1]+imm] = uint8(r[rs2] & 0xff)
			m[r[rs1]+imm+1] = uint8((r[rs2] >> 8) & 0xff)
            pc = pc + 4

		} else if funct3 == 0b010 {
			// SW
			m[r[rs1]+imm] = uint8(r[rs2] & 0xff)
			m[r[rs1]+imm+1] = uint8((r[rs2] >> 8) & 0xff)
			m[r[rs1]+imm+2] = uint8((r[rs2] >> 16) & 0xff)
			m[r[rs1]+imm+3] = uint8((r[rs2] >> 24) & 0xff)
            pc = pc + 4

		}

	} else if opcode == 0b1110011 {
		// ECALL
		c = ZVMHandleECall(c)

	} else {
		s = 5
	}

	if s == 4 {
		s = 2
	}

	c.pc = pc
	c.r = r
	c.m = m
	c.s = s

	return c
}

func ZVMRun(c ZVMContext) ZVMContext {
	for c.s == 2 {
		c = ZVMTick(c)
	}
	return c
}

func main() {
	var c ZVMContext
	c = ZVMReset(c)

	data := []uint32{0b1_00001_0110111, 0b1_00010_0010111}
	var a uint32
	a = c.pc

	for _, v := range data {
		c.m[a+0] = uint8((v >> 0) & 0xff)
		c.m[a+1] = uint8((v >> 8) & 0xff)
		c.m[a+2] = uint8((v >> 16) & 0xff)
		c.m[a+3] = uint8((v >> 24) & 0xff)
		a += 4
	}

	c = ZVMRun(c)

	fmt.Println(c.r)
	fmt.Println(c.pc)
	fmt.Println(c.s)
}
