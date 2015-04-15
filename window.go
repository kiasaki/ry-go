package main

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

type Window struct {
	Buffer  *Buffer
	CursorX int
	CursorY int
	Editor  *Editor
}

func NewWindow(e *Editor, b *Buffer) *Window {
	return &Window{
		Buffer:  b,
		CursorX: 0,
		CursorY: 0,
		Editor:  e,
	}
}
