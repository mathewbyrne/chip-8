package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mathewbyrne/chip-8/chip8"
)

type Game struct {
	c *chip8.Chip8
}

func (g *Game) Update() error {
	return g.c.Tick()
}

func (g *Game) Draw(screen *ebiten.Image) {
	var buff [4 * 64 * 32]byte
	var pixel = [4]byte{0xFF, 0xFF, 0xFF, 0xFF}
	fb := g.c.FrameBuffer()
	for i := range fb {
		for j := 0; j < 8; j++ {
			if fb[i]>>j&0x1 == 0x1 {
				copy(buff[32*i+4*j:], pixel[:])
			}
		}
	}

	screen.WritePixels(buff[:])
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 64, 32
}

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

	k := &Input{}

	c, err := chip8.NewChip8(f, k)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}

	ebiten.SetWindowSize(640, 320)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{c}); err != nil {
		log.Fatal(err)
	}
}
