# CHIP-8 Interpretter in Go

A simple [CHIP-8](https://en.wikipedia.org/wiki/CHIP-8) interpretter written in Go.

```
go run . [path to CHIP-8 ROM]
```

This is a very simple implementation that uses [`ebiten`](https://github.com/hajimehoshi/ebiten) for rendering and input.  Use` `chip8.Runner` to run a `chip8.Chip8` on a separate thread and communicate.

Passes all tests in https://github.com/Timendus/chip8-test-suite for the CHIP-8 target.
