package frontends

type Frontend interface {
	Init() error
	Close() error

	Clear(fg, bg Attribute) error
	SetCursor(x, y int) // -1, -1 being hide cursor
	SetCell(x, y int, ch rune, fg, bg Attribute)

	PollEvent() (Event, error)
	CancelPollEvent()

	Size() (int, int)
	Flush() error
}

type EventType uint8

const (
	EventKey EventType = iota
	EventResize
	EventMouse
	EventNone
)

type Event interface {
	Type() EventType
}
