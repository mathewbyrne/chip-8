package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mathewbyrne/chip-8/chip8"
)

type Game struct {
	r chip8.Runner
	f *chip8.FrameBuffer
}

func (g *Game) Update() error {
	var m uint16
	for k, e := range keyMap {
		if ebiten.IsKeyPressed(e) {
			m |= 1 << k
		}
	}
	g.r.SetKeyMap(m)

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.r.Pause()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.r.Step()
	}
	return nil
}

var pixel = [4]byte{0xFF, 0xFF, 0xFF, 0xFF}

func (g *Game) Draw(screen *ebiten.Image) {
	var buff [4 * 64 * 32]byte
	if g.f.Dirty() {
		for i, b := range g.f.Data() {
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
		log.Fatal("usage: chip8 [romfile]")
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("could not open rom file for reading: %v", err)
	}

	c, err := chip8.New(f)
	if err != nil {
		log.Fatalf("%v", err)
	}

	r := chip8.Run(c)
	defer r.Close()

	ebiten.SetWindowSize(640, 320)
	ebiten.SetWindowTitle(os.Args[1])
	ebiten.SetScreenClearedEveryFrame(false)

	if err := ebiten.RunGame(&Game{r, c.FrameBuffer()}); err != nil {
		log.Fatal(err)
	}
}
