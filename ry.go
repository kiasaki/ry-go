package main

import (
	"log"
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

	// debuging
	go func() {
		time.Sleep(10 * time.Second)
		// close before writing to stdout
		editor.Close()
		log.Fatal("Timeout")
	}()

	for editor.Running {
		editor.Update()
		editor.Draw()
	}
}
