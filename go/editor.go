package main

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/kiasaki/ry/frontends"
)

func DEBUG(o interface{}) {
	logFile, _ := os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	spew.Fdump(logFile, o)
}

type Editor struct {
	Running bool

	Buffers      Buffers
	WindowTree   *WindowTree
	ActiveWindow *Window
	Frontend     frontends.Frontend

	StatusLeft       []byte
	StatusLeftStyle  frontends.Attribute
	StatusRight      []byte
	StatusRightStyle frontends.Attribute

	EvChan chan (frontends.Event)

	Height int
	Width  int
}

func NewEditor(f frontends.Frontend) *Editor {
	e := &Editor{
		Running:      false,
		Buffers:      Buffers{},
		ActiveWindow: nil,
		Frontend:     f,
	}
	e.WindowTree = NewWindowTree(e)
	return e
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

///////////////////////////////////////////////////////////////////////////////
// Update / Handle / Draw loop
///////////////////////////////////////////////////////////////////////////////

func (e *Editor) handleEvent(event frontends.Event) {
	switch event.Type() {
	case frontends.EventResize:
		e.Height = event.Height()
		e.Width = event.Width()
		break
	case frontends.EventKey:
		e.handleKey(event)
		break
	}
}

func (e *Editor) handleKey(event frontends.Event) {
	e.ActiveWindow.Buffer.MajorMode().HandleKey(e, event)
}

func (e *Editor) Update() {
	// if no buffer left, exit
	if len(e.Buffers) == 0 {
		e.Running = false
		return
	}

	// if no window left, select first buffer
	if e.ActiveWindow == nil {
		topLeftWindow := e.WindowTree.TopLeftMostWindow()
		if topLeftWindow != nil {
			// set as active
			e.ActiveWindow = topLeftWindow
		} else {
			// all windows closed, open first buffer
			e.WindowTree.Split = WindowSplitLeaf
			e.WindowTree.Leaf = NewWindow(e, e.Buffers[0])
			e.ActiveWindow = e.WindowTree.Leaf
		}
	}

	// event pooling from frontend
	select {
	case ev := <-e.EvChan:
		// TODO optimize and fetch as much event as possible before going to redraw
		e.handleEvent(ev)
	case <-time.After(50 * time.Millisecond):
		break
	}
	// updates
}

func (e *Editor) Draw() {
	// TODO remove debugging
	DEBUG(e)

	// status bar
	statusBarRect := NewRect(0, e.Height-1, e.Width/2, 1)
	drawBytes(e.Frontend, statusBarRect, e.StatusLeft,
		e.StatusLeftStyle, frontends.ColorDefault, TextAlignLeft)
	drawBytes(e.Frontend, statusBarRect.SetX(e.Width/2), e.StatusRight,
		e.StatusRightStyle, frontends.ColorDefault, TextAlignRight)

	// windows
	r := NewRect(0, 0, e.Width, e.Height-1)
	e.WindowTree.Draw(e.Frontend, r)

	e.Frontend.Flush()
}

///////////////////////////////////////////////////////////////////////////////
// Status bar
///////////////////////////////////////////////////////////////////////////////

func (e *Editor) SetStatusLeft(b []byte, fg frontends.Attribute) {
	e.StatusLeft = b
	e.StatusLeftStyle = fg
}

func (e *Editor) SetStatusRight(b []byte, fg frontends.Attribute) {
	e.StatusRight = b
	e.StatusRightStyle = fg
}

///////////////////////////////////////////////////////////////////////////////
// Creatings buffers
///////////////////////////////////////////////////////////////////////////////

func (e *Editor) NewFileBuffer(path string) (*Buffer, error) {
	path = substituteHome(path)
	absolute, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	basename := filepath.Base(path)

	_, err = os.Stat(path)
	if err != nil {
		// assume non-exiting file / to be saved
		return NewBuffer(e.UniqueName(basename+" (new)"), absolute), nil
	} else {
		buffer := NewBuffer(e.UniqueName(basename), absolute)
		err = buffer.ReadFromDisk()
		return buffer, err
	}
}

func (e *Editor) NewEmptyBuffer() *Buffer {
	return NewBuffer(e.UniqueName("(scratch)"), "")
}

func (e *Editor) FindBuffer(name string) (*Buffer, bool) {
	for _, b := range e.Buffers {
		if b.Name == name {
			return b, true
		}
	}
	return nil, false
}

func (e *Editor) AppendBuffer(b *Buffer) {
	e.Buffers = append(e.Buffers, b)
}

func (e *Editor) UniqueName(name string) string {
	if _, found := e.FindBuffer(name); !found {
		return name
	}
	for i := 2; i < 9001; i++ {
		numberedName := name + " <" + strconv.Itoa(i) + ">"
		if _, found := e.FindBuffer(numberedName); !found {
			return numberedName
		}
	}
	panic("Too many buffer with the same name")
}