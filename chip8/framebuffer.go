package chip8

import (
	"fmt"
	"strings"
)

const (
	FB_WIDTH  uint8 = 64
	FB_HEIGHT uint8 = 32
	FB_LEN          = (uint16(FB_WIDTH) * uint16(FB_HEIGHT)) >> 3
)

// 1 bit per pixel
type FrameBuffer [FB_LEN]byte

// draw 8 bit wide sprite at x, y position. out of bounds drawing wraps.
// Drawing is done via XOR.  Returns true if any existing bits are set to 0.
func (f *FrameBuffer) draw(sprite []byte, x, y uint8) bool {
	// wrap these values if they themselves are outside the bounds
	x = x % FB_WIDTH
	y = y % FB_HEIGHT

	var collision bool
	li := (x >> 3) + (8 * y)
	ri := li + 1
	if (ri % 8) == 0 {
		ri -= 8
	}
	offset := x % 8

	for i := range sprite {

		ls := sprite[i] >> offset
		collision = (f[li] & ls) > 0
		f[li] = f[li] ^ ls

		if offset > 0 {
			rs := sprite[i] << (8 - offset)
			collision = (f[ri] & rs) > 0
			f[ri] = f[ri] ^ rs
		}

		// uint8 will have them wrap automatically back to the opposite side in the correct position.  This isn't good programming,
		// just a coincidence (or nice property I guess) of the size of the chip-8 framebuffer.
		li += 8
		ri += 8
	}

	return collision
}

func (f *FrameBuffer) clear() {
	for i := range f {
		f[i] = 0
	}
}

func (f *FrameBuffer) String() string {
	var sb strings.Builder
	for y := 0; y < 32; y++ {
		sb.Write([]byte(fmt.Sprintf("%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b\n", f[y*8], f[y*8+1], f[y*8+2], f[y*8+3], f[y*8+4], f[y*8+5], f[y*8+6], f[y*8+7])))
	}
	return sb.String()
}
