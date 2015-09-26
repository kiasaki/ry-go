package main

import (
	"fmt"

	"github.com/kiasaki/ry/lang"
	"github.com/tiborvass/uniline"
)

func main() {
	prompt := "ry> "
	scanner := uniline.DefaultScanner()

	for scanner.Scan(prompt) {
		line := scanner.Text()
		if len(line) > 0 {
			scanner.AddToHistory(line)
			pr(eval(line))
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func eval(code string) string {
	result, err := lang.Parse([]byte(code))
	if err != nil {
		panic(err)
	}
	return result
}

func pr(result string) {
	fmt.Println(result)
}
