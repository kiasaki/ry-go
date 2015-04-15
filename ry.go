package main

import (
	"log"
	"os"
	"time"

	"github.com/kiasaki/ry/frontends"
)

func main() {
	editor := NewEditor(frontends.TermboxFrontend{})

	err := editor.Init()
	if err != nil {
		panic(err)
	}
	defer editor.Close()

	for _, path := range os.Args[1:] {
		buffer, err := editor.NewFileBuffer(path)
		if err != nil {
			editor.Close()
			log.Fatal(err)
		}
		editor.AppendBuffer(buffer)
	}

	// debuging
	go func() {
		time.Sleep(10 * time.Second)
		// close before writing to stdout
		editor.Close()
		log.Fatal("Timeout")
	}()

	// event pooling
	go func(c chan frontends.Event) {
		for {
			ev, err := editor.Frontend.PollEvent()
			if err != nil {
				panic(err)
			}
			c <- ev
		}
	}(editor.EvChan)

	// main loop
	for editor.Running {
		editor.Update()
		editor.Draw()
	}
}