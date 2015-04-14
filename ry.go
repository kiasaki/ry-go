package main

import (
	"time"

	"github.com/kiasaki/ry/frontends"
)

func main() {
	frontend := frontends.TermboxFrontend{}

	err := frontend.Init()
	if err != nil {
		panic(err)
	}
	defer frontend.Close()

	i := 9
	for i > 0 {
		frontend.SetCell(1, 1, ':', frontends.ColorMagenta, frontends.ColorDefault)
		frontend.SetCell(2, 1, ')', frontends.ColorMagenta, frontends.ColorDefault)
		frontend.SetCell(1, 4, rune(i+64), frontends.ColorCyan, frontends.ColorWhite)
		frontend.Flush()
		time.Sleep(1 * time.Second)
		i -= 1
	}
}
