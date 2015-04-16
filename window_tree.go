package main

import (
	"strconv"
	"strings"

	"github.com/kiasaki/ry/frontends"
)

type WindowSplit uint8

const (
	WindowSplitH WindowSplit = iota
	WindowSplitV
	WindowSplitLeaf
)

type WindowSizing uint8

const (
	WindowSizingFixed WindowSizing = iota
	WindowSizingScale
)

type WindowTree struct {
	Top         *WindowTree
	Bottom      *WindowTree
	Left        *WindowTree
	Right       *WindowTree
	Split       WindowSplit
	Leaf        *Window
	Editor      *Editor
	Sizing      WindowSizing
	SizingScale float32
	SizingFixed int
}

func NewWindowTree(e *Editor) *WindowTree {
	return &WindowTree{
		Editor:      e,
		Split:       WindowSplitLeaf,
		Leaf:        nil,
		Sizing:      WindowSizingScale,
		SizingScale: 0.5,
	}
}

func (wt *WindowTree) TopLeftMostWindow() *Window {
	switch wt.Split {
	case WindowSplitLeaf:
		return wt.Leaf
	case WindowSplitH:
		return wt.Top.TopLeftMostWindow()
	case WindowSplitV:
		return wt.Left.TopLeftMostWindow()
	}
	return nil // should never happen
}

func (wt *WindowTree) Draw(f frontends.Frontend, r Rect) {
	// TODO do not asume we only have a leaf at root (multiple windows will come)
	window := wt.Leaf

	lineCount := window.Buffer.LineCount()
	numberGutterWidth := len(strconv.Itoa(lineCount))

	// display cursor if this is the active window
	if wt.Editor.ActiveWindow == window {
		f.SetCursor(window.CursorX+numberGutterWidth+1, window.CursorY)
	}

	// draw lines (numbers and contents)
	offset := r.X + numberGutterWidth + 1
	for i := 0; i < window.Buffer.LineCount()-1 && i < r.Height; i++ {
		// line number
		drawBytes(f, NewRect(r.X, r.Y+i, numberGutterWidth, 1), []byte(strconv.Itoa(i+1)),
			frontends.ColorYellow, frontends.ColorDefault, TextAlignRight)
		// line contents
		line := window.Buffer.Lines[i]
		drawBytes(f, NewRect(offset, r.Y+i, r.Width-offset, 1), line.Contents,
			frontends.ColorWhite, frontends.ColorDefault, TextAlignLeft)
	}

	// draw place holders for non lines
	// TODO this if is a bit naive (chack cases with scrolling and scroll padding)
	for y := r.Y + window.Buffer.LineCount() - 1; y < r.Y2()-1; y++ {
		f.SetCell(r.X, y, '~', frontends.ColorYellow, frontends.ColorDefault)
	}

	// draw footer
	fillRect(f, NewRect(r.X, r.Y2()-1, r.Width, 1), ' ', frontends.ColorWhite, frontends.ColorBlue)

	// draw footer txt
	bufferName := []byte(window.Buffer.Name)
	statusLine := []byte("(" + strings.Join(window.Buffer.ModeNames(), ", ") +
		") (" + strconv.Itoa(window.CursorX) + ", " + strconv.Itoa(window.CursorY) + ")")
	drawBytes(f, r.SetY(r.Y2()-1).SetHeight(1), bufferName, frontends.ColorWhite, frontends.ColorBlue, TextAlignLeft)
	drawBytes(f, r.SetY(r.Y2()-1).SetHeight(1), statusLine, frontends.ColorWhite, frontends.ColorBlue, TextAlignRight)
}
