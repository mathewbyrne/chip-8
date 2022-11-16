package chip8

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFrameBufferCornerEdgeCase(t *testing.T) {
	fb := &FrameBuffer{}

	spr := []byte{
		0b11111111,
		0b11000001,
		0b10011001,
		0b10011011,
		0b10000101,
		0b11111111,
	}

	// draw sprite precisely in the corner
	fb.draw(spr, 60, 29)

	fmt.Println(fb)

	require.EqualValues(t, 0b00001001, fb[255])
}
