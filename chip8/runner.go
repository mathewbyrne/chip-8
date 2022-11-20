package chip8

import (
	"fmt"
	"time"
)

type Runner interface {
	Pause()
	Step()
	Close()
}

type runner chan uint

const (
	CMD_START uint = iota
	CMD_PAUSE
	CMD_STEP
	CMD_CLOSE
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
				switch cmd {
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
				}
			}
		}

	}()

	return r
}
