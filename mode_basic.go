package main

type BasicMode struct {
}

func NewBasicMode() Mode {
	return Mode(BasicMode{})
}

func (BasicMode) Major() bool {
	return true
}

func (BasicMode) Name() string {
	return "basic"
}
