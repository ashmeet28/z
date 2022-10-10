package main

import "fmt"

type ZVMContext struct {
	r  []uint32
	m  []uint8
	ec bool
}

func ZVMReset(c ZVMContext) ZVMContext {
	c.m = make([]uint8, 4294967296)
	c.r = make([]uint32, 32)
	c.ec = false
	return c
}

func ZVMTick(c ZVMContext) ZVMContext {
	return c
}

func main() {
	var c ZVMContext

	c = ZVMReset(c)

	c = ZVMTick(c)

	fmt.Println(c.m[0xffffffff])
}
