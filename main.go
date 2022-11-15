package main

import (
	"fmt"
	"os"

	"github.com/mathewbyrne/chip-8/chip8"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: chip8 [romfile]")
		return
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open rom file for reading: %v", err)
		return
	}

	c, err := chip8.NewChip8(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}

	c.Tick()
	c.Tick()
	c.Tick()
	c.Tick()
}
