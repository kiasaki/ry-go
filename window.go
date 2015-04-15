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
	top         *WindowTree
	bottom      *WindowTree
	left        *WindowTree
	right       *WindowTree
	split       WindowSplit
	leaf        *Window
	Editor      *Editor
	Sizing      WindowSizing
	SizingScale float32
	SizingFixed int
}

type Window struct {
	Title    string
	Filepath string
	Modes    Modes
	Buffer   *Buffer
	CursorX  int
	CursorY  int
	Editor   *Editor
}
