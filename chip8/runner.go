package chip8

import (
	"fmt"
	"sync"
	"time"
)

type Runner struct {
	mu     sync.Mutex
	paused bool
}

func (r *Runner) Pause() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.paused = !r.paused
}

func (r *Runner) Run(c *Chip8) {
	cycle := time.NewTicker(time.Second / 500)
	defer cycle.Stop()
	timers := time.NewTicker(time.Second / 60)
	defer timers.Stop()

	for {
		if !r.paused {
			select {
			case <-timers.C:
				c.Tick()
			case <-cycle.C:
				op := c.Cycle()
				fmt.Printf("%04X [%04X] %s\n", c.pc, uint16(op), op)
				fmt.Printf("%s\n", c)
			}
		}
	}
}
