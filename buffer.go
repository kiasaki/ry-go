package main

import (
	"io/ioutil"
	"strings"
)

type Line struct {
	Contents []byte
	Buffer   *Buffer
}

type Buffer struct {
	Modes       Modes
	Name        string
	Filepath    string
	Lines       []*Line
	CurrentLine int
	Changed     bool
}

type Buffers []*Buffer

func NewBuffer(name, filepath string) *Buffer {
	return &Buffer{
		Modes:    Modes{NewBasicMode()},
		Name:     name,
		Filepath: filepath,
		Lines:    []*Line{},
		Changed:  false,
	}
}

func (b *Buffer) Saveable() bool {
	return b.Filepath != ""
}

func (b *Buffer) ReadFromDisk() error {
	bytes, err := ioutil.ReadFile(b.Filepath)
	if err != nil {
		return err
	}

	b.Lines = []*Line{}
	for _, lineContents := range strings.Split(string(bytes), "\n") {
		b.Lines = append(b.Lines, &Line{
			Contents: []byte(lineContents),
			Buffer:   b,
		})
	}

	return nil
}

func (b *Buffer) ModeNames() (names []string) {
	for _, mode := range b.Modes {
		names = append(names, mode.Name())
	}
	return
}

func (b *Buffer) LineCount() int {
	return len(b.Lines)
}
