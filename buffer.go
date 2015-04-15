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
	return &Buffer{
		Modes:    Modes{NewBasicMode()},
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

func (b *Buffer) ModeNames() (names []string) {
	for _, mode := range b.Modes {
		names = append(names, mode.Name())
	}
	return
}
