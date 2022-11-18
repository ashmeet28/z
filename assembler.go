package main

import "fmt"

func ASMGetLines(data []uint8) [][]uint8 {
	var Lines [][]uint8
	var LineLength int
	var LastIndex int
	LastIndex = 0
	for i := range data {
		if data[i] == 0x0a {
			if LineLength != 0 {
				Lines = append(Lines, data[LastIndex:(LastIndex+LineLength)])
			}
			LastIndex = i + 1
			LineLength = 0
		} else {
			LineLength++
		}
	}
	return Lines
}
func main() {
	d := []uint8{0x0a, 0x0a, 2, 5, 7, 9, 0x0a, 0x0a, 12, 55, 0x0a, 0x0a}
	fmt.Println(ASMGetLines(d))
}
