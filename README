# CHIP-8 Interpretter in Go

A simple [CHIP-8](https://en.wikipedia.org/wiki/CHIP-8) interpretter written in Go.

```
go run . [path to CHIP-8 ROM]
```

This is a very simple implementation that uses [`ebiten`](https://github.com/hajimehoshi/ebiten) for rendering and input.  Emulation runs on a separate thread so that it can appropriately block for input when `OP_LD_VX_K` opcode is encountered.  CHIP-8 cycles run around 500Hz.