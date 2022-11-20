package chip8

import (
	"fmt"
	"time"
)

type Runner interface {
	Close()
	Pause()
	Step()
}

type runner chan uint

const (
	CMD_PAUSE uint = iota
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

	go func() {
		for {
			select {
			case <-timers.C:
				c.Tick()
			case <-cycle.C:
				if !paused {
					op := c.Cycle()
					fmt.Printf("%04X [%04X] %s\n", c.pc, uint16(op), op)
					fmt.Printf("%s\n", c)
				}
			case cmd := <-r:
				switch cmd {
				case CMD_PAUSE:
					paused = !paused
				case CMD_STEP:
					if paused {
						op := c.Cycle()
						fmt.Printf("%04X [%04X] %s\n", c.pc, uint16(op), op)
						fmt.Printf("%s\n", c)
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
