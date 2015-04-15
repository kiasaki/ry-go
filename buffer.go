package main

import (
	"io/ioutil"
)

type Buffer struct {
	Name     string
	Filepath string
	Contents []byte
	Changed  bool
}

type Buffers []*Buffer

func NewBuffer(name, filepath string) *Buffer {
	return &Buffer{
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
