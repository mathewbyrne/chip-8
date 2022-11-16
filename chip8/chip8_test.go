package chip8

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpAddVxVy(t *testing.T) {
	c := Chip8{}
	c.r[0] = 0xFE
	c.r[1] = 1

	c.opAddVxVy(0, 1)
	require.EqualValues(t, 0xFF, c.r[0])
	require.EqualValues(t, 0, c.r[0xF])

	c.opAddVxVy(0, 1)
	require.EqualValues(t, 0, c.r[0])
	require.EqualValues(t, 1, c.r[0xF])
}

func TestOpLdVxI(t *testing.T) {
	c := Chip8{}

	copy(c.r[:], []byte{0x01, 0x02, 0x03, 0x04, 0x05})
	c.i = 0x400

	c.opLdVxI(0)
	require.Equal(t, []byte{0x01, 0x00}, c.m[0x400:0x400+2])
	require.EqualValues(t, 0x401, c.i)

	c.i = 0x400
	c.m[0x405] = 0x0F
	c.opLdVxI(4)
	require.Equal(t, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x0F}, c.m[0x400:0x400+6])
	require.EqualValues(t, 0x405, c.i)

	c.i = 0x400
	c.r[15] = 0x0F
	c.opLdVxI(15)
	require.Equal(t, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F}, c.m[0x400:0x400+16])
	require.EqualValues(t, 0x410, c.i)
}

func TestOpLdIVx(t *testing.T) {
	c := Chip8{}

	copy(c.m[0x0400:], []byte{0x01, 0x02, 0x03, 0x04, 0x05})
	c.i = 0x400

	c.opLdIVx(0)
	require.Equal(t, []byte{0x01, 0x00}, c.r[0:2])
	require.EqualValues(t, 0x401, c.i)

	c.opLdIVx(2)
	require.Equal(t, []byte{0x02, 0x03, 0x04, 0x00}, c.r[0:4])
	require.EqualValues(t, 0x404, c.i)
}

func TestOpLdBVx(t *testing.T) {
	c := Chip8{}
	c.i = 0x200

	c.r[4] = 9
	c.opLdBVx(4)
	require.Equal(t, []byte{0, 0, 9}, c.m[0x200:0x200+3])

	c.r[4] = 123
	c.opLdBVx(4)
	require.Equal(t, []byte{1, 2, 3}, c.m[0x200:0x200+3])

	c.r[4] = 255
	c.opLdBVx(4)
	require.Equal(t, []byte{2, 5, 5}, c.m[0x200:0x200+3])
}
