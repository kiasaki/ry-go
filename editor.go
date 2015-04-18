package ry

import (
	"time"

	"github.com/kiasaki/ry/frontends"
	syp "github.com/kiasaki/syp-lang/interpreter"
)

type Editor struct {
	Running bool

	Buffers      []*SexpBuffer
	WindowTree   *SexpWindow
	ActiveWindow *SexpWindow
	Frontend     frontends.Frontend

	EvChan  chan frontends.Event
	KeyChan chan syp.Sexp

	Height int
	Width  int
}

var globalEditor *Editor

func init() {
	globalEditor = &Editor{
		Running:      false,
		Buffers:      []*SexpBuffer{},
		ActiveWindow: nil,
		Frontend:     frontends.TermboxFrontend{},
	}
	globalEditor.WindowTree = NewWindow(globalEditor)
}

func StartEditorFunction(env *syp.Lang, fnname string,
	args []syp.Sexp) (syp.Sexp, error) {

	err := globalEditor.Frontend.Init()
	if err != nil {
		return syp.SexpNull, err
	}
	globalEditor.Width, globalEditor.Height = globalEditor.Frontend.Size()
	globalEditor.EvChan = make(chan frontends.Event, 100)
	globalEditor.Running = true

	// goroutine for feeding the event channel
	go func(c chan frontends.Event) {
		for globalEditor.Running {
			ev, err := globalEditor.Frontend.PollEvent()
			if err != nil {
				globalEditor.Frontend.Close()
				panic(err)
			}
			c <- ev
		}
	}(globalEditor.EvChan)

	// goroutine for handling events
	go func(c chan frontends.Event) {
		for globalEditor.Running {
			select {
			case event := <-c:
				// TODO optimize and fetch as much event as possible before going to redraw
				switch event.Type() {
				case frontends.EventResize:
					globalEditor.Height = event.Height()
					globalEditor.Width = event.Width()
					break
				case frontends.EventKey:
					globalEditor.KeyChan <- syp.SexpChar(event.Character())
					break
				}
			case <-time.After(50 * time.Millisecond):
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
	}(globalEditor.EvChan)

	// goroutine flushing updates to terminal
	go func() {
		for globalEditor.Running {
			globalEditor.Frontend.Flush()
			time.Sleep(50 * time.Millisecond)
		}
	}()

	return syp.SexpNull, nil
}

func StopEditorFunction(env *syp.Lang, fnname string,
	args []syp.Sexp) (syp.Sexp, error) {
	if globalEditor.Running {
		globalEditor.Running = false
		return syp.SexpNull, globalEditor.Frontend.Close()
	}
	return syp.SexpNull, nil
}
