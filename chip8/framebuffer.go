package chip8

import "fmt"

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
	fmt.Printf("draw %d %d\n", x, y)
	var collision bool
	li := (x >> 3) + (8 * y)
	ri := li + 1
	if (ri % 8) == 0 {
		ri -= 8
	}
	offset := x % 8

	for i := range sprite {

		fmt.Printf("\tdraw %08b to (%d, %d) offset: %d\n", sprite[i], li, ri, offset)

		ls := sprite[i] >> offset
		collision = (f[li] & ls) > 0
		f[li] = f[li] ^ ls
		fmt.Printf("\tl sprite %08b\n", ls)

		if offset > 0 {
			rs := sprite[i] << (8 - offset)
			collision = (f[ri] & rs) > 0
			f[ri] = f[ri] ^ rs
			fmt.Printf("\tr sprite %08b\n", rs)
		}

		// uint8 will have them wrap automatically back to the top
		li += 8
		ri += 8
	}

	fmt.Println("< ====")
	for y := 0; y < 32; y++ {
		fmt.Printf("\t%08b|%08b|%08b|%08b|%08b|%08b|%08b|%08b\n", f[y*8], f[y*8+1], f[y*8+2], f[y*8+3], f[y*8+4], f[y*8+5], f[y*8+6], f[y*8+7])
	}
	fmt.Println("< ====")

	fmt.Printf("collision: %t\n", collision)
	return collision
}
