package main

import (
	"fmt"

	"github.com/kiasaki/ry/lang"
	"github.com/tiborvass/uniline"
)

func main() {
	prompt := "ry> "
	scanner := uniline.DefaultScanner()
	env := lang.NewBuiltinFilledEnv()

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

func read(code string) []lang.Value {
	result, err := lang.Parse([]byte(code))
	if err != nil {
		return []lang.Value{lang.NewSymbolValue("ERROR: " + err.Error())}
	}
	return result
}

func eval(env *lang.Env, exprs []lang.Value) string {
	var lastResult lang.Value
	var err error

	if len(exprs) == 0 {
		return ""
	}

	for _, expr := range exprs {
		lastResult, err = lang.Eval(expr, env)
	}
	if err != nil {
		return "ERROR: " + err.Error()
	}

	return lastResult.String()
}

func pr(result string) {
	fmt.Println(result)
}
