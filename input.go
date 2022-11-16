package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mathewbyrne/chip-8/chip8"
)

type Input struct{}

var keyMap = map[chip8.Key]ebiten.Key{
	0x0: ebiten.Key0,
	0x1: ebiten.Key1,
	0x2: ebiten.KeyArrowUp,
	0x3: ebiten.Key3,
	0x4: ebiten.KeyArrowLeft,
	0x5: ebiten.Key4,
	0x6: ebiten.KeyArrowRight,
	0x7: ebiten.Key7,
	0x8: ebiten.KeyArrowDown,
	0x9: ebiten.Key9,
	0xA: ebiten.KeyA,
	0xB: ebiten.KeyB,
	0xC: ebiten.KeyC,
	0xD: ebiten.KeyD,
	0xE: ebiten.KeyE,
	0xF: ebiten.KeyF,
}

func (i *Input) State(k chip8.Key) bool {
	return ebiten.IsKeyPressed(keyMap[k])
}

func (i *Input) Wait() chip8.Key {
	panic("not implemented")
}
