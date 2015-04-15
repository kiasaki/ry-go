package main

import (
	"io/ioutil"
)

type Buffer struct {
	Modes    Modes
	Name     string
	Filepath string
	Contents []byte
	Changed  bool
}

type Buffers []*Buffer

func NewBuffer(name, filepath string) *Buffer {
	basicMode := NewBasicMode()
	return &Buffer{
		Modes:    Modes{&basicMode},
		Name:     name,
		Filepath: filepath,
		Changed:  false,
	}
}

func (b *Buffer) Saveable() bool {
	return b.Filepath != ""
}

func (b *Buffer) ReadFromDisk() error {
	bytes, err := ioutil.ReadFile(b.Filepath)
	b.Contents = bytes
	return err
}
