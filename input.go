package main

import "github.com/mathewbyrne/chip-8/chip8"

type Input struct{}

func (i *Input) State(k chip8.Key) bool {
	return false
}

func (i *Input) Wait() chip8.Key {
	return 0
}
