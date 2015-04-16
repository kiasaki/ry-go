package main

import (
	"github.com/kiasaki/ry/frontends"
)

type NormalMode struct {
	Buffer *Buffer
}

func NewNormalMode(b *Buffer) Mode {
	return Mode(NormalMode{b})
}

func (NormalMode) Major() bool {
	return true
}

func (NormalMode) Name() string {
	return "normal"
}

func (m NormalMode) HandleKey(e *Editor, event frontends.Event) {
	e.SetStatusLeft(append(e.StatusLeft, byte(event.Character())), frontends.ColorWhite)
}
