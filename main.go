package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mathewbyrne/chip-8/chip8"
)

type Game struct {
	c      *chip8.Chip8
	paused bool
}

func (g *Game) Update() error {
	// if g.paused && inpututil.IsKeyJustPressed(ebiten.KeyS) {
	// 	// no tick :/
	// 	g.c.Cycle()
	// }

	// if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
	// 	g.paused = !g.paused
	// }

	// if g.paused && inpututil.IsKeyJustPressed(ebiten.KeyF) {
	// 	fmt.Printf("%s\n", g.c.FrameBuffer())
	// }

	return nil
}

var pixel = [4]byte{0xFF, 0xFF, 0xFF, 0xFF}

func (g *Game) Draw(screen *ebiten.Image) {
	var buff [4 * 64 * 32]byte
	fb := g.c.FrameBuffer()
	if fb.Dirty() {
		for i, b := range fb.Data() {
			for j := 0; j < 8; j++ {
				if (b>>j)&0x1 == 0x1 {
					copy(buff[32*i+4*(7-j):], pixel[:])

				}
			}
		}
		screen.WritePixels(buff[:])
	}
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
	go runChip8(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}

	ebiten.SetWindowSize(640, 320)
	ebiten.SetWindowTitle(os.Args[1])
	ebiten.SetScreenClearedEveryFrame(false)
	if err := ebiten.RunGame(&Game{c, false}); err != nil {
		log.Fatal(err)
	}
}

func runChip8(c *chip8.Chip8) {
	cycle := time.Tick(time.Second / 500)
	timers := time.Tick(time.Second / 60)
	for {
		select {
		case <-timers:
			c.Tick()
		case <-cycle:
			// default:
			c.Cycle()
		}
	}
}
