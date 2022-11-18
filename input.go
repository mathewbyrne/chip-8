package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mathewbyrne/chip-8/chip8"
)

type Input struct{}

var keyMap = map[chip8.Key]ebiten.Key{
	// Row 1
	0x1: ebiten.Key1,
	0x2: ebiten.Key2,
	0x3: ebiten.Key3,
	0xC: ebiten.Key4,
	// Row 2
	0x4: ebiten.KeyQ,
	0x5: ebiten.KeyW,
	0x6: ebiten.KeyE,
	0xD: ebiten.KeyR,
	// Row 3
	0x7: ebiten.KeyA,
	0x8: ebiten.KeyS,
	0x9: ebiten.KeyD,
	0xE: ebiten.KeyF,
	// Row 4
	0xA: ebiten.KeyZ,
	0x0: ebiten.KeyX,
	0xB: ebiten.KeyC,
	0xF: ebiten.KeyV,
}

func (i *Input) State(k chip8.Key) bool {
	return ebiten.IsKeyPressed(keyMap[k])
}

func (i *Input) Wait() chip8.Key {
	ticker := time.NewTicker(time.Microsecond * 17) // ~60Hz poll
	defer ticker.Stop()
	for {
		<-ticker.C
		for k := range keyMap {
			if ebiten.IsKeyPressed(keyMap[k]) {
				return k
			}
		}
	}
}
