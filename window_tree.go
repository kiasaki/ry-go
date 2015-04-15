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
	window := wt.Leaf
	fillRect(f, NewRect(r.X, r.Y2()-1, r.Width, 1), ' ', frontends.ColorWhite, frontends.ColorBlue)

	bufferName := []byte(window.Buffer.Name)
	statusLine := []byte("(" + strings.Join(window.Buffer.ModeNames(), ", ") +
		") (" + strconv.Itoa(window.CursorX) + ", " + strconv.Itoa(window.CursorY) + ")")
	drawBytes(f, r.SetY(r.Y2()-1).SetHeight(1), bufferName, frontends.ColorWhite, frontends.ColorBlue, TextAlignLeft)
	drawBytes(f, r.SetY(r.Y2()-1).SetHeight(1), statusLine, frontends.ColorWhite, frontends.ColorBlue, TextAlignRight)
}
