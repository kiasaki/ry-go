package main

import (
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/kiasaki/ry/frontends"
)

type Editor struct {
	Running bool

	Buffers  *Buffers
	Windows  *Windows
	Frontend frontends.Frontend

	EvChan chan (frontends.Event)

	Height int
	Width  int
}

func NewEditor(f frontends.Frontend) *Editor {
	return &Editor{
		Running:  false,
		Buffers:  &Buffers{},
		Windows:  &Windows{},
		Frontend: f,
	}
}

func (e *Editor) Init() error {
	err := e.Frontend.Init()
	if err != nil {
		return err
	}
	e.Width, e.Height = e.Frontend.Size()
	e.EvChan = make(chan frontends.Event, 100)
	e.Running = true
	return nil
}

func (e *Editor) Close() error {
	if e.Running {
		e.Running = false
		return e.Frontend.Close()
	}
	return nil
}

func (e *Editor) handleEvent(event frontends.Event) {
	if event.Type() == frontends.EventResize {
		e.Height = event.Height()
		e.Width = event.Width()
	}
}

func (e *Editor) Update() {
	select {
	case ev := <-e.EvChan:
		e.handleEvent(ev)
	case <-time.After(50 * time.Millisecond):
		break
	}
	// updates
}

func (e *Editor) Draw() {
	text := []byte(strconv.Itoa(e.Height))
	textlen := utf8.RuneCount(text)
	drawNFirstRunes(e.Frontend, 2, 2, textlen, frontends.ColorWhite, frontends.ColorDefault, text)

	e.Frontend.Flush()
}
