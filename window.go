package ry

import (
	"fmt"
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

type SexpWindow struct {
	Top         *SexpWindow
	Bottom      *SexpWindow
	Left        *SexpWindow
	Right       *SexpWindow
	Split       WindowSplit
	Leaf        *SexpBuffer
	Editor      *Editor
	Sizing      WindowSizing
	SizingScale float32
	SizingFixed int
}

func (w SexpWindow) SexpString() string {
	return fmt.Sprintf("#<window %x>", &w)
}

func NewWindow(e *Editor) *SexpWindow {
	return &SexpWindow{
		Editor:      e,
		Split:       WindowSplitLeaf,
		Leaf:        nil,
		Sizing:      WindowSizingScale,
		SizingScale: 0.5,
	}
}
