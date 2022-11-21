package chip8

import (
	"fmt"
	"strings"
	"sync"
)

const (
	FB_WIDTH  uint8 = 64
	FB_HEIGHT uint8 = 32
	FB_LEN          = (uint16(FB_WIDTH) * uint16(FB_HEIGHT)) >> 3
)

// 1 bit per pixel
type FrameBuffer struct {
	b     [FB_LEN]byte
	mut   sync.Mutex
	dirty bool
}

// draw 8 bit wide sprite at x, y position. out of bounds drawing wraps.
// Drawing is done via XOR.  Returns true if any existing bits are set to 0.
func (f *FrameBuffer) draw(sprite []byte, x, y uint8) bool {
	f.mut.Lock()
	defer f.mut.Unlock()

	f.dirty = true

	// wrap these values if they themselves are outside the bounds
	x = x % FB_WIDTH
	y = y % FB_HEIGHT

	var collision bool
	li := (x >> 3) + (8 * y)
	ri := li + 1
	offset := x % 8

	for i := range sprite {

		ls := sprite[i] >> offset
		if !collision {
			collision = (f.b[li] & ls) > 0
		}
		f.b[li] = f.b[li] ^ ls

		if offset > 0 && ri%8 != 0 {
			rs := sprite[i] << (8 - offset)
			if !collision {
				collision = (f.b[ri] & rs) > 0
			}
			f.b[ri] = f.b[ri] ^ rs
		}

		// do not wrap around to the top if the sprite is drawn at the bottom of the bounds. Wrapping only occurs if the sprite
		// is explicitly drawn past the bounds.
		if uint16(li)+8 > 255 {
			break
		}
		li += 8
		ri += 8
	}

	return collision
}

func (f *FrameBuffer) clear() {
	f.mut.Lock()
	defer f.mut.Unlock()
	f.dirty = true
	for i := range f.b {
		f.b[i] = 0
	}
}

func (f *FrameBuffer) Data() []byte {
	f.mut.Lock()
	defer f.mut.Unlock()
	return f.b[:]
}

func (f *FrameBuffer) String() string {
	f.mut.Lock()
	defer f.mut.Unlock()
	var sb strings.Builder
	for y := 0; y < 32; y++ {
		sb.Write([]byte(fmt.Sprintf("%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b\n", f.b[y*8], f.b[y*8+1], f.b[y*8+2], f.b[y*8+3], f.b[y*8+4], f.b[y*8+5], f.b[y*8+6], f.b[y*8+7])))
	}
	return sb.String()
}

func (f *FrameBuffer) Dirty() bool {
	f.mut.Lock()
	defer f.mut.Unlock()
	d := f.dirty
	f.dirty = false
	return d
}
