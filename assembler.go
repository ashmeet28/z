package main

import (
	"fmt"
	"os"
	"strings"
)

type ZCContext struct {
	pc     uint32
	labels map[string]uint32
	lines  []string
}

func ZCReset(c ZCContext) ZCContext {
	c.pc = 0x8000000
	c.labels = make(map[string]uint32)
	c.lines = make([]string, 0, 0)
	return c
}

func ZCGetLines(c ZCContext, s string) ZCContext {
	for _, v := range strings.Split(s, "\x0a") {
		if v != "" && v[0:1] != "#" {
			c.lines = append(c.lines, v)
		}
	}
	return c
}

func ZCGetLabels(c ZCContext) ZCContext {
	var words []string
	var pc uint32

	pc = c.pc
	for _, v := range c.lines {
		words = strings.Split(v, " ")
		if words[0] == "label" {
			c.labels[words[1]] = pc
		} else {
			pc = pc + 4
		}
	}
	return c
}

func main() {
	var c ZCContext
	c = ZCReset(c)
	d, _ := os.ReadFile("/home/xenon/a")
	c = ZCGetLines(c, string(d))
	c = ZCGetLabels(c)
	fmt.Println(c)
}
