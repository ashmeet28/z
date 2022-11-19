package main

import "fmt"

type ZVMContext struct {
	pc uint32
	r  []uint32
	m  []uint8
	s  uint32
}

func ZVMReset(c ZVMContext, data []byte) ZVMContext {
	if len(data) < 0x8000000 {
		c.s = 0x1
		return c
	}

	c.m = make([]uint8, 0x100000000, 0x100000000)
	c.r = make([]uint32, 32, 32)
	c.pc = 0x8000000

	for i := range data {
		c.m[i+0x8000000] = uint8(data[i])
	}

	c.s = 0x2

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
	var funct7 uint32
	var funct12 uint32
	var shamt uint32

	var pc uint32
	var r []uint32
	var m []uint8
	var s uint32

	pc = c.pc
	r = c.r
	m = c.m
	s = c.s

	r[0] = 0x0
	s = 0x4

	inst = uint32(m[pc]) | (uint32(m[pc+1]) << 8) | (uint32(m[pc+2]) << 16) | (uint32(m[pc+3]) << 24)
	opcode = inst & 0x7f

	if opcode == 0b0110111 {
		// LUI
		rd = (inst >> 7) & 0x1f
		imm = inst & 0xfffff000

		r[rd] = imm
		pc = pc + 4

		s = 0x2
	} else if opcode == 0b0010111 {
		// AUIPC
		rd = (inst >> 7) & 0x1f
		imm = inst & 0xfffff000

		r[rd] = imm + pc
		pc = pc + 4

		s = 0x2
	} else if opcode == 0b1101111 {
		// JAL
		rd = (inst >> 7) & 0x1f
		imm = ((inst & 0x80000000) >> 11) | (inst & 0xff000) | ((inst >> 9) & 0x800) | ((inst >> 20) & 0x7fe)
		if (imm >> 20) == 0x1 {
			imm = imm | 0xffe00000
		}

		r[rd] = pc + 4
		pc = pc + imm

		s = 0x2
	} else if opcode == 0b1100111 {
		rd = (inst >> 7) & 0x1f
		funct3 = (inst >> 12) & 0x7
		rs1 = (inst >> 15) & 0x1f
		imm = inst >> 20
		if (imm >> 11) == 0x1 {
			imm = imm | 0xfffff000
		}

		if funct3 == 0b000 {
			// JALR
			r[rd] = pc + 4
			pc = (imm + r[rs1]) & 0xfffffffe

			s = 0x2
		}
	} else if opcode == 0b1100011 {
		imm = ((inst & 0x80000000) >> 19) | ((inst & 0x80) << 4) | ((inst >> 20) & 0x7e0) | ((inst >> 7) & 0x1e)
		if (imm & 0x1000) == 0x1000 {
			imm = imm | 0xffffe000
		}
		funct3 = (inst >> 12) & 0x7
		rs1 = (inst >> 15) & 0x1f
		rs2 = (inst >> 20) & 0x1f

		if funct3 == 0b000 {
			// BEQ
			if r[rs1] == r[rs2] {
				pc = imm + pc
			} else {
				pc = pc + 4
			}

			s = 0x2
		} else if funct3 == 0b001 {
			// BNE
			if r[rs1] != r[rs2] {
				pc = imm + pc
			} else {
				pc = pc + 4
			}

			s = 0x2
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

			s = 0x2
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

			s = 0x2
		} else if funct3 == 0b110 {
			// BLTU
			if r[rs1] < r[rs2] {
				pc = imm + pc
			} else {
				pc = pc + 4
			}

			s = 0x2
		} else if funct3 == 0b111 {
			// BGEU
			if r[rs1] >= r[rs2] {
				pc = imm + pc
			} else {
				pc = pc + 4
			}

			s = 0x2
		}
	} else if opcode == 0b0000011 {
		rd = (inst >> 7) & 0x1f
		funct3 = (inst >> 12) & 0x7
		rs1 = (inst >> 15) & 0x1f
		imm = inst >> 20
		if (imm >> 11) == 0x1 {
			imm = imm | 0xfffff000
		}

		if funct3 == 0b000 {
			// LB
			r[rd] = uint32(m[r[rs1]+imm])
			if (r[rd] & 0x80) == 0x80 {
				r[rd] = r[rd] | 0xffffff80
			}
			pc = pc + 4

			s = 0x2
		} else if funct3 == 0b001 {
			// LH
			r[rd] = uint32(m[r[rs1]+imm]) | (uint32(m[r[rs1]+imm+1]) << 8)
			if (r[rd] & 0x8000) == 0x8000 {
				r[rd] = r[rd] | 0xffff0000
			}
			pc = pc + 4

			s = 0x2
		} else if funct3 == 0b010 {
			// LW
			r[rd] = uint32(m[r[rs1]+imm]) | (uint32(m[r[rs1]+imm+1]) << 8) | (uint32(m[r[rs1]+imm+2]) << 16) | (uint32(m[r[rs1]+imm+3]) << 24)
			if (r[rd] & 0x8000) == 0x8000 {
				r[rd] = r[rd] | 0xffff0000
			}
			pc = pc + 4

			s = 0x2
		} else if funct3 == 0b100 {
			// LBU
			r[rd] = uint32(m[r[rs1]+imm])
			pc = pc + 4

			s = 0x2
		} else if funct3 == 0b101 {
			// LHU
			r[rd] = uint32(m[r[rs1]+imm]) | (uint32(m[r[rs1]+imm+1]) << 8)
			pc = pc + 4

			s = 0x2
		}
	} else if opcode == 0b0100011 {
		imm = ((inst >> 7) & 0x1f) | ((inst >> 20) & 0xfe0)
		if (imm >> 11) == 0x1 {
			imm = imm | 0xfffff000
		}
		funct3 = (inst >> 12) & 0x7
		rs1 = (inst >> 15) & 0x1f
		rs2 = (inst >> 20) & 0x1f

		if funct3 == 0b000 {
			// SB
			m[r[rs1]+imm] = uint8(r[rs2] & 0xff)
			pc = pc + 4

			s = 0x2
		} else if funct3 == 0b001 {
			// SH
			m[r[rs1]+imm] = uint8(r[rs2] & 0xff)
			m[r[rs1]+imm+1] = uint8((r[rs2] >> 8) & 0xff)
			pc = pc + 4

			s = 0x2
		} else if funct3 == 0b010 {
			// SW
			m[r[rs1]+imm] = uint8(r[rs2] & 0xff)
			m[r[rs1]+imm+1] = uint8((r[rs2] >> 8) & 0xff)
			m[r[rs1]+imm+2] = uint8((r[rs2] >> 16) & 0xff)
			m[r[rs1]+imm+3] = uint8((r[rs2] >> 24) & 0xff)
			pc = pc + 4

			s = 0x2
		}
	} else if opcode == 0b0010011 {
		rd = (inst >> 7) & 0x1f
		funct3 = (inst >> 12) & 0x7
		rs1 = (inst >> 15) & 0x1f
		imm = inst >> 20
		if (imm & 0x800) == 0x800 {
			imm = imm | 0xfffff000
		}

		if funct3 == 0b000 {
			// ADDI
			r[rd] = r[rs1] + imm
			pc = pc + 4

			s = 0x2
		} else if funct3 == 0b010 {
			// SLTI
			if ((r[rs1] >> 31) == 0x0) && ((imm >> 31) == 0x0) {
				if r[rs1] < imm {
					r[rd] = 0x1
				} else {
					r[rd] = 0x0
				}
			} else if ((r[rs1] >> 31) == 0x1) && ((r[rs2] >> 31) == 0x1) {
				if r[rs1] > imm {
					r[rd] = 0x1
				} else {
					r[rd] = 0x0
				}
			} else if (r[rs1] >> 31) == 0x1 {
				r[rd] = 0x1
			} else {
				r[rd] = 0x0
			}
			pc = pc + 4

			s = 0x2
		} else if funct3 == 0b011 {
			// SLTIU
			if r[rs1] < imm {
				r[rd] = 0x1
			} else {
				r[rd] = 0x0
			}
			pc = pc + 4

			s = 0x2
		} else if funct3 == 0b100 {
			// XORI
			r[rd] = r[rs1] ^ imm
			pc = pc + 4

			s = 0x2
		} else if funct3 == 0b110 {
			// ORI
			r[rd] = r[rs1] | imm
			pc = pc + 4

			s = 0x2
		} else if funct3 == 0b111 {
			// ANDI
			r[rd] = r[rs1] & imm
			pc = pc + 4

			s = 0x2
		}
	} else if opcode == 0b0010011 {
		rd = (inst >> 7) & 0x1f
		funct3 = (inst >> 12) & 0x7
		rs1 = (inst >> 15) & 0x1f
		shamt = inst & 0x1f
		imm = inst >> 20

		if (funct3 == 0b001) && ((imm & 0xfe0) == 0x0) {
			// SLLI
			r[rd] = r[rs1] << shamt
			pc = pc + 4

			s = 0x2
		} else if funct3 == 0b101 {
			if (imm & 0xfe0) == 0x0 {
				// SRLI
				r[rd] = r[rs1] >> shamt
				pc = pc + 4

				s = 0x2
			} else if (imm & 0xfe0) == 0x400 {
				// SRAI
				r[rd] = r[rs1] >> shamt
				if (r[rs1] & 0x80000000) == 0x80000000 {
					r[rd] = r[rd] | (0xffffffff << (32 - shamt))
				}
				pc = pc + 4

				s = 0x2
			}
		}
	} else if opcode == 0b0110011 {
		rd = (inst >> 7) & 0x1f
		funct3 = (inst >> 12) & 0x7
		rs1 = (inst >> 15) & 0x1f
		rs2 = (inst >> 20) & 0x1f
		funct7 = inst >> 25

		if funct3 == 0b000 {
			if funct7 == 0x0 {
				// ADD
				r[rd] = r[rs1] + r[rs2]
				pc = pc + 4

				s = 0x2
			} else if funct7 == 0x20 {
				// SUB
				r[rd] = r[rs1] - r[rs2]
				pc = pc + 4

				s = 0x2
			}
		} else if (funct3 == 0b001) && (funct7 == 0x0) {
			// SLL
			r[rd] = r[rs1] << (r[rs2] & 0x1f)
			pc = pc + 4

			s = 0x2
		} else if (funct3 == 0b010) && (funct7 == 0x0) {
			// SLT
			if ((r[rs1] >> 31) == 0x0) && ((imm >> 31) == 0x0) {
				if r[rs1] < r[rs2] {
					r[rd] = 0x1
				} else {
					r[rd] = 0x0
				}
			} else if ((r[rs1] >> 31) == 0x1) && ((r[rs2] >> 31) == 0x1) {
				if r[rs1] > r[rs2] {
					r[rd] = 0x1
				} else {
					r[rd] = 0x0
				}
			} else if (r[rs1] >> 31) == 0x1 {
				r[rd] = 0x1
			} else {
				r[rd] = 0x0
			}
			pc = pc + 4

			s = 0x2
		} else if (funct3 == 0b011) && (funct7 == 0x0) {
			// SLTU
			if r[rs1] < imm {
				r[rd] = 0x1
			} else {
				r[rd] = 0x0
			}
			pc = pc + 4

			s = 0x2
		} else if (funct3 == 0b100) && (funct7 == 0x0) {
			// XOR
			r[rd] = r[rs1] ^ r[rs2]
			pc = pc + 4

			s = 0x2
		} else if funct3 == 0b101 {
			if funct7 == 0x0 {
				// SRL
				r[rd] = r[rs1] >> (r[rs2] & 0x1f)
				pc = pc + 4

				s = 0x2
			} else if funct7 == 0x20 {
				// SRA
				r[rd] = r[rs1] >> (r[rs2] & 0x1f)
				if (r[rs1] >> 31) == 0x1 {
					r[rd] = r[rd] | (0xffffffff << (32 - (r[rs2] & 0x1f)))
				}
				pc = pc + 4

				s = 0x2
			}
		} else if (funct3 == 0b110) && (funct7 == 0x0) {
			// OR
			r[rd] = r[rs1] | r[rs2]
			pc = pc + 4

			s = 0x2
		} else if (funct3 == 0b111) && (funct7 == 0x0) {
			// AND
			r[rd] = r[rs1] & r[rs2]
			pc = pc + 4

			s = 0x2
		}
	} else if opcode == 0b1110011 {
		// ECALL
		rd = (inst >> 7) & 0x1f
		funct3 = (inst >> 12) & 0x7
		rs1 = (inst >> 15) & 0x1f
		funct12 = inst >> 20

		if (rd == 0x0) && (funct3 == 0x0) && (rs1 == 0x0) && (funct12 == 0x0) {
			s = 0x5
		}
	}

	if s == 0x4 {
		s = 0x8
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
	data := make([]uint8, 1000000)

	c = ZVMReset(c, data)

	c = ZVMRun(c)

	fmt.Println(c.r)
	fmt.Println(c.pc)
	fmt.Println(c.s)
}
