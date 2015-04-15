package main

import (
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/kiasaki/ry/frontends"
)

type TextAlign uint8

const (
	TextAlignLeft TextAlign = iota
	TextAlignRight
	TextAlignCenter
)

type Rect struct {
	X      int
	Y      int
	Height int
	Width  int
}

func NewRect(x, y, width, height int) Rect {
	return Rect{X: x, Y: y, Width: width, Height: height}
}

func (r Rect) X2() int {
	return r.X + r.Width
}

func (r Rect) Y2() int {
	return r.Y + r.Height
}

func (r Rect) SetY(y int) Rect {
	return NewRect(r.X, y, r.Width, r.Height)
}

func (r Rect) SetHeight(height int) Rect {
	return NewRect(r.X, r.Y, r.Width, height)
}

func drawNFirstRunes(fe frontends.Frontend, off, y, n int, text []byte, fg, bg frontends.Attribute) {
	for n > 0 {
		r, size := utf8.DecodeRune(text)
		fe.SetCell(off, y, r, fg, bg)
		text = text[size:]
		off++
		n--
	}
}

func drawBytes(fe frontends.Frontend, r Rect, text []byte, fg, bg frontends.Attribute, align TextAlign) {
	textlen := utf8.RuneCount(text)
	if textlen > r.Width {
		textlen = r.Width
	}
	switch align {
	case TextAlignLeft:
		drawNFirstRunes(fe, r.X, r.Y, textlen, text, fg, bg)
		break
	case TextAlignRight:
		drawNFirstRunes(fe, r.X2()-textlen, r.Y, textlen, text, fg, bg)
		break
	case TextAlignCenter:
		drawNFirstRunes(fe, (r.X2()/2)-(textlen/2), r.Y, textlen, text, fg, bg)
		break
	}
}

func fillRect(fe frontends.Frontend, r Rect, ru rune, fg, bg frontends.Attribute) {
	for x := r.X; x < r.X2(); x++ {
		for y := r.Y; y < r.Y2(); y++ {
			fe.SetCell(x, y, ru, fg, bg)
		}
	}
}

func substituteHome(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}
	home := os.Getenv("HOME")
	if home == "" {
		panic("HOME is not set")
	}
	return filepath.Join(home, path[1:])
}
