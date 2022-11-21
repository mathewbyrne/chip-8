package chip8

import (
	"fmt"
	"time"
)

type Runner interface {
	Pause()
	Step()
	Close()
	SetKeyMap(uint16)
}

type runner chan uint32

const (
	CMD_START uint32 = iota
	CMD_PAUSE
	CMD_STEP
	CMD_CLOSE
	CMD_KEYMAP
)

func (r runner) Close() {
	r <- CMD_CLOSE
}

func (r runner) Pause() {
	r <- CMD_PAUSE
}

func (r runner) Step() {
	r <- CMD_STEP
}

func (r runner) SetKeyMap(keymap uint16) {
	r <- uint32(keymap)<<16 | CMD_KEYMAP
}

// Run runs a CHIP-8 chip as a go routine and returns a struct that can interface to the running simulation.
func Run(c *Chip8) runner {
	r := make(runner)

	paused := false
	cycle := time.NewTicker(time.Second / 1000)
	timers := time.NewTicker(time.Second / 60)

	doCycle := func() {
		op, pc := c.Cycle()
		fmt.Printf("%04X [%04X] %s\n", pc, uint16(op), op)
		fmt.Printf("%s\n", c)
	}

	go func() {
		for {
			select {
			case <-timers.C:
				c.Tick()
			case <-cycle.C:
				if !paused {
					doCycle()
				}
			case cmd := <-r:
				switch cmd & 0x00FF {
				case CMD_PAUSE:
					paused = !paused
				case CMD_STEP:
					if paused {
						doCycle()
					}
				case CMD_CLOSE:
					cycle.Stop()
					timers.Stop()
					return
				case CMD_KEYMAP:
					c.KeyMap.Store(cmd >> 16)
				}
			}
		}

	}()

	return r
}
