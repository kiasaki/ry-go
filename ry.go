package ry

import (
	"fmt"
	"os"
	"time"

	"github.com/jonvaldes/termo"
)

var editor *Editor

type Editor struct {
	width       int
	height      int
	framebuffer *termo.Framebuffer

	keyChan chan (termo.ScanCode)
	errChan chan (error)

	ticker <-chan (time.Time)
}

func (e *Editor) Start() {
	if err := termo.Init(); err != nil {
		panic(err)
	}
	termo.ShowCursor()

	defer func() {
		termo.Stop()
		if err := recover(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	// Get size
	var err error
	e.width, e.height, err = termo.Size()
	if err != nil {
		panic(err)
	}

	// Create frame buffer
	e.framebuffer = termo.NewFramebuffer(e.width, e.height)

	e.startKeyReadLoop()

	e.MainLoop()
}

func (e *Editor) startKeyReadLoop() {
	// Read keys
	e.keyChan = make(chan termo.ScanCode, 100)
	e.errChan = make(chan error)
	termo.StartKeyReadLoop(e.keyChan, e.errChan)
	e.ticker = time.Tick(100 * time.Millisecond)
}

func (e *Editor) MainLoop() {
	for {
		// Check for terminal resize
		if _w, _h, _ := termo.Size(); e.width != _w || e.height != _h {
			e.width = _w
			e.height = _h
			e.framebuffer = termo.NewFramebuffer(e.width, e.height)
		}

		// Clear framebuffer
		e.framebuffer.Clear()

		e.Render(e.framebuffer)

		// Read keyboard
		select {
		case <-e.ticker:
			// Periodically flush framebuffer to screen
			e.framebuffer.Flush()
		case s := <-e.keyChan:
			e.handleScanCode(s)
		case err := <-e.errChan:
			panic(err)
		}
	}
}

func (e *Editor) Render(f *termo.Framebuffer) {
	statusLine := []rune("--**")
	f.AttribText(0, e.height-1, termo.CellState{
		termo.AttrNone,
		termo.ColorBlack,
		termo.ColorGray,
	}, string(padString(statusLine, '-', e.width)))
	f.AttribRect(0, 0, e.width, e.height-1, termo.CellState{
		termo.AttrNone,
		termo.ColorGray,
		termo.ColorBlack,
	})
}

func (e *Editor) handleScanCode(s termo.ScanCode) {
	if s.IsMouseMoveEvent() || s.IsMouseDownEvent() || s.IsMouseUpEvent() {
		// don't bother
	} else if s.IsEscapeCode() {
		switch s.EscapeCode() {
		case 65: // Up
		case 66: // Down
		case 67: // Right
		case 68: // Left
		}
	} else {
		r := s.Rune()
		// Exit if Ctrl+C or Esc are pressed
		if r == 3 || r == 27 {
			e.framebuffer.Clear()
			e.framebuffer.Flush()
			termo.Stop()
			os.Exit(0)
		}
	}
}

func padString(str []rune, pad rune, width int) []rune {
	for len([]rune(str)) < width {
		str = append(str, pad)
	}
	return str
}
