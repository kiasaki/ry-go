package ry

import (
	"errors"
	"strconv"
	"time"

	"github.com/kiasaki/ry/frontends"
	sypext "github.com/kiasaki/syp-lang/extensions"
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
	globalEditor.EvChan = make(chan frontends.Event, 50)
	globalEditor.KeyChan = make(chan syp.Sexp, 50)
	globalEditor.Running = true

	// goroutine for feeding the event channel
	go processPollEvents(globalEditor.EvChan)

	// goroutine for handling events
	go processHandleEvents(globalEditor.EvChan)

	// goroutine flushing updates to terminal
	go processFlushUpdates()

	return syp.SexpNull, nil
}

func processPollEvents(c chan frontends.Event) {
	for globalEditor.Running {
		ev, err := globalEditor.Frontend.PollEvent()
		if err != nil {
			globalEditor.Frontend.Close()
			panic(err)
		}
		c <- ev
	}
}

func processHandleEvents(c chan frontends.Event) {
	for globalEditor.Running {
		select {
		case event := <-c:
			// TODO optimize and fetch as much event as possible before going to redraw
			switch event.Type() {
			case frontends.EventResize:
				globalEditor.Height = event.Height()
				globalEditor.Width = event.Width()
			case frontends.EventKey:
				globalEditor.KeyChan <- syp.SexpChar(event.Character())
			}
		case <-time.After(50 * time.Millisecond):
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func processFlushUpdates() {
	for globalEditor.Running {
		globalEditor.Frontend.Flush()
		time.Sleep(50 * time.Millisecond)
	}
}

func StopEditorFunction(env *syp.Lang, fnname string,
	args []syp.Sexp) (syp.Sexp, error) {
	if globalEditor.Running {
		globalEditor.Running = false
		return syp.SexpNull, globalEditor.Frontend.Close()
	}
	return syp.SexpNull, nil
}

func EditorKeypressesChan(env *syp.Lang, fnname string,
	args []syp.Sexp) (syp.Sexp, error) {
	if globalEditor.KeyChan == nil {
		return syp.SexpNull, errors.New("Cant return keypresses channel, editor is not started")
	}
	return sypext.SexpChannel(globalEditor.KeyChan), nil
}

func ClearEditorFunction(env *syp.Lang, fnname string,
	args []syp.Sexp) (syp.Sexp, error) {
	if len(args) != 2 {
		return syp.SexpNull, syp.WrongNargs
	}

	var fg syp.SexpInt
	switch expr := args[0].(type) {
	case syp.SexpInt:
		fg = expr
	default:
		return syp.SexpNull, errors.New("clear-editor fg is not an int")
	}

	var bg syp.SexpInt
	switch expr := args[1].(type) {
	case syp.SexpInt:
		bg = expr
	default:
		return syp.SexpNull, errors.New("clear-editor bg is not an int")
	}

	globalEditor.Frontend.Clear(frontends.Attribute(fg), frontends.Attribute(bg))
	return syp.SexpNull, nil
}

func SetCellFunction(env *syp.Lang, fnname string,
	args []syp.Sexp) (syp.Sexp, error) {
	if len(args) != 5 {
		return syp.SexpNull, syp.WrongNargs
	}
	// args 0, 1, 3 and 4 need to be ints
	castedArgs := make([]syp.SexpInt, 5)
	for i := 0; i < 5 && i != 2; i++ {
		switch expr := args[i].(type) {
		case syp.SexpInt:
			castedArgs[i] = expr
		default:
			return syp.SexpNull, errors.New("set-cell arg #" + strconv.Itoa(i+1) + " is not an int")
		}
	}

	var ch syp.SexpChar
	switch expr := args[2].(type) {
	case syp.SexpChar:
		ch = expr
	default:
		return syp.SexpNull, errors.New("set-cell arg #3 is not an char")
	}

	globalEditor.Frontend.SetCell(
		int(castedArgs[0]),
		int(castedArgs[1]),
		rune(ch),
		frontends.Attribute(castedArgs[3]),
		frontends.Attribute(castedArgs[4]),
	)

	return syp.SexpNull, nil
}
