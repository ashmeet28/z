package main

import (
	"fmt"
	"os"
	"strings"
)

var InstOpcodes = map[string]uint32{
	"lui":   0b0110111,
	"auipc": 0b0010111,
	"jal":   0b1101111,
	"jalr":  0b1100111,
	"beq":   0b000_00000_1100011,
	"bne":   0b001_00000_1100011,
	"blt":   0b100_00000_1100011,
	"bge":   0b101_00000_1100011,
	"bltu":  0b110_00000_1100011,
	"bgeu":  0b111_00000_1100011,
	"lb":    0b000_00000_0000011,
	"lh":    0b001_00000_0000011,
	"lw":    0b010_00000_0000011,
	"lbu":   0b100_00000_0000011,
	"lhu":   0b101_00000_0000011,
	"sb":    0b000_00000_0100011,
	"sh":    0b001_00000_0100011,
	"sw":    0b010_00000_0100011,
	"addi":  0b000_00000_0010011,
	"slti":  0b010_00000_0010011,
	"sltiu": 0b011_00000_0010011,
	"xori":  0b100_00000_0010011,
	"ori":   0b110_00000_0010011,
	"andi":  0b111_00000_0010011,
	"slli":  0b0_000000000000000_001_00000_0010011,
	"srli":  0b0_000000000000000_101_00000_0010011,
	"srai":  0b1_000000000000000_101_00000_0010011,
	"add":   0b0_000000000000000_000_00000_0110011,
	"sub":   0b1_000000000000000_000_00000_0110011,
	"sll":   0b0_000000000000000_001_00000_0110011,
	"slt":   0b0_000000000000000_010_00000_0110011,
	"sltu":  0b0_000000000000000_011_00000_0110011,
	"xor":   0b0_000000000000000_100_00000_0110011,
	"srl":   0b0_000000000000000_101_00000_0110011,
	"sra":   0b1_000000000000000_101_00000_0110011,
	"or":    0b0_000000000000000_110_00000_0110011,
	"and":   0b0_000000000000000_111_00000_0110011,
	"ecall": 0b1110011,
}

func ZCCompile(s string) []uint8 {
	var pc uint32
	var labels map[string]uint32
	var instructions []string
	var words []string
	// var inst uint32

	pc = 0x8000000
	labels = make(map[string]uint32)
	instructions = make([]string, 0, 0)

	for _, v := range strings.Split(s, "\x0a") {
		if v != "" && v[0:1] != "#" {
			words = strings.Split(v, " ")
			if words[0] == "label" {
				labels[words[1]] = pc
			} else {
				instructions = append(instructions, v)
				pc = pc + 4
			}
		}
	}

	pc = 0x8000000
	return nil
}

func main() {
	d, _ := os.ReadFile("/home/xenon/a")

	ZCCompile(string(d))
	fmt.Println(InstOpcodes)
}
