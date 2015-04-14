package main

import (
	"unicode/utf8"

	"github.com/kiasaki/ry/frontends"
)

func drawNFirstRunes(fe frontends.Frontend, y, off, n int, fg, bg frontends.Attribute, text []byte) {
	for n > 0 {
		r, size := utf8.DecodeRune(text)
		fe.SetCell(off, y, r, fg, bg)
		text = text[size:]
		off++
		n--
	}
}
