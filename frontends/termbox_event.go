package frontends

import (
	"github.com/nsf/termbox-go"
)

type TermboxEvent struct {
	tbEvent termbox.Event
}

var _ = Event(TermboxEvent{})

func NewTermboxEvent(event termbox.Event) TermboxEvent {
	return TermboxEvent{tbEvent: event}
}

func (e TermboxEvent) Type() EventType {
	switch e.tbEvent.Type {
	case termbox.EventKey:
		return EventKey
	case termbox.EventResize:
		return EventResize
	case termbox.EventMouse:
		return EventMouse
	default:
		return EventNone
	}
}

func (e TermboxEvent) Height() int {
	return e.tbEvent.Height
}

func (e TermboxEvent) Width() int {
	return e.tbEvent.Width
}
