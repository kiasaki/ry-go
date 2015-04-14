package frontends

import (
	"errors"

	"github.com/nsf/termbox-go"
)

var (
	ErrorTbInterrupt = errors.New("Event error: Interrupt")
	ErrorTbUnknown   = errors.New("Event error: Unknown")
)

type TermboxFrontend struct{}

func (TermboxFrontend) Init() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	termbox.SetInputMode(termbox.InputAlt)
	termbox.SetOutputMode(termbox.OutputNormal)
	return nil
}

func (TermboxFrontend) Close() error {
	termbox.Close()
	return nil
}

func (TermboxFrontend) Clear(fg, bg Attribute) error {
	return termbox.Clear(termbox.Attribute(fg), termbox.Attribute(bg))
}

func (TermboxFrontend) SetCursor(x, y int) {
	termbox.SetCursor(x, y)
}

func (TermboxFrontend) SetCell(x, y int, ch rune, fg, bg Attribute) {
	termbox.SetCell(x, y, ch, termbox.Attribute(fg), termbox.Attribute(bg))
}

func (TermboxFrontend) PollEvent() (Event, error) {
	event := termbox.PollEvent()
	switch event.Type {
	case termbox.EventKey:
	case termbox.EventResize:
	case termbox.EventMouse:
		return NewTermboxEvent(event), nil
	case termbox.EventError:
		return nil, event.Err
	case termbox.EventInterrupt:
		return nil, ErrorTbInterrupt
	}
	return nil, ErrorTbUnknown
}

func (TermboxFrontend) CancelPollEvent() {
	termbox.Interrupt()
}

func (TermboxFrontend) Size() (int, int) {
	return termbox.Size()
}

func (TermboxFrontend) Flush() error {
	return termbox.Flush()
}

var _ = Frontend(TermboxFrontend{})
