package main

import (
	"fmt"

	"github.com/kiasaki/ry/lang"
	"github.com/tiborvass/uniline"
)

func main() {
	prompt := "ry> "
	scanner := uniline.DefaultScanner()
	env := lang.NewRootEnv()

	for scanner.Scan(prompt) {
		line := scanner.Text()
		if len(line) > 0 {
			scanner.AddToHistory(line)
			pr(eval(env, read((line))))
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func read(code string) []Value {
	result, err := lang.Parse([]byte(code))
	if err != nil {
		panic(err)
	}
	return result
}

func eval(env *lang.Env, exprs []Value) string {
	var lastResult Value

	for _, expr := range exprs {
		lastResult, err = lang.Eval(env, expr)
	}
	if err != nil {
		panic(err)
	}

	return lastResult.String()
}

func pr(result string) {
	fmt.Println(result)
}
