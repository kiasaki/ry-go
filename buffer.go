package ry

import ()

type Buffer struct {
	Name     string
	Filename string
	Point    uint32
	Contents []byte
}

func NewBuffer(name string, filename string) *Buffer {
	return &Buffer{
		Name:     name,
		Filename: filename,
		Point:    0,
		Contents: []byte{},
	}
}
