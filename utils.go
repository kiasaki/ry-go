package ry

import (
	"errors"
)

func padString(str []rune, pad rune, width int) []rune {
	for len([]rune(str)) < width {
		str = append(str, pad)
	}
	return str
}

func die(err error) {
	panic(errors.New("ERROR: " + err.Error()))
}
