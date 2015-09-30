package ry

import (
	"errors"

	"github.com/jonvaldes/termo"
)

func stringToTerminalAttr(name string) termo.Attribute {
	switch name {
	case "bold":
		return termo.AttrBold
	case "dim":
		return termo.AttrDim
	case "under":
		return termo.AttrUnder
	case "blink":
		return termo.AttrBlink
	case "rev":
		return termo.AttrRev
	case "hid":
		return termo.AttrHid
	default:
		return termo.AttrNone
	}
}

func stringToTerminalColor(name string) termo.Color {
	switch name {
	case "black":
		return termo.ColorBlack
	case "light-black":
		return termo.ColorBlack.Light()
	case "red":
		return termo.ColorRed
	case "light-red":
		return termo.ColorRed.Light()
	case "green":
		return termo.ColorGreen
	case "light-green":
		return termo.ColorGreen.Light()
	case "yellow":
		return termo.ColorYellow
	case "light-yellow":
		return termo.ColorYellow.Light()
	case "blue":
		return termo.ColorBlue
	case "light-blue":
		return termo.ColorBlue.Light()
	case "magenta":
		return termo.ColorMagenta
	case "light-magenta":
		return termo.ColorMagenta.Light()
	case "cyan":
		return termo.ColorCyan
	case "light-cyan":
		return termo.ColorCyan.Light()
	case "gray":
		return termo.ColorGray
	case "light-gray":
		return termo.ColorGray.Light()
	default:
		return termo.ColorDefault
	}
}

func die(err error) {
	panic(errors.New("ERROR: " + err.Error()))
}
