package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	ry "github.com/kiasaki/ry"
	sypext "github.com/kiasaki/syp-lang/extensions"
	syp "github.com/kiasaki/syp-lang/interpreter"
)

var exitOnFailure = flag.Bool("exitonfail", false,
	"exit on failure instead of starting repl")

func getLine(reader *bufio.Reader) (string, error) {
	line := make([]byte, 0)
	for {
		linepart, hasMore, err := reader.ReadLine()
		if err != nil {
			return "", err
		}
		line = append(line, linepart...)
		if !hasMore {
			break
		}
	}
	return string(line), nil
}

func isBalanced(str string) bool {
	parens := 0
	squares := 0

	for _, c := range str {
		switch c {
		case '(':
			parens++
		case ')':
			parens--
		case '[':
			squares++
		case ']':
			squares--
		}
	}

	return parens == 0 && squares == 0
}

func getExpression(reader *bufio.Reader) (string, error) {
	fmt.Printf("> ")
	line, err := getLine(reader)
	if err != nil {
		return "", err
	}
	for !isBalanced(line) {
		fmt.Printf(">> ")
		nextline, err := getLine(reader)
		if err != nil {
			return "", err
		}
		line += "\n" + nextline
	}
	return line, nil
}

func processDumpCommand(env *syp.Lang, args []string) {
	if len(args) == 0 {
		env.DumpEnvironment()
	} else {
		err := env.DumpFunctionByName(args[0])
		if err != nil {
			fmt.Println(err)
		}
	}
}

func repl(env *syp.Lang) {
	fmt.Printf("ry version %s\n", ry.Version())
	fmt.Printf("syp version %s\n", syp.Version())
	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := getExpression(reader)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		parts := strings.Split(line, " ")
		if len(parts) == 0 {
			continue
		}

		if parts[0] == "quit" {
			break
		}

		if parts[0] == "dump" {
			processDumpCommand(env, parts[1:])
			continue
		}

		expr, err := env.EvalString(line)
		if err != nil {
			fmt.Print(env.GetStackTrace(err))
			env.Clear()
			continue
		}

		if expr != syp.SexpNull {
			fmt.Println(expr.SexpString())
		}
	}
}

func runScript(env *syp.Lang, fname string) {
	file, err := os.Open(fname)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer file.Close()

	err = env.LoadFile(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	_, err = env.Run()
	if err != nil {
		fmt.Print(env.GetStackTrace(err))
		if *exitOnFailure {
			os.Exit(-1)
		}
		repl(env)
	}
}

func main() {
	env := syp.NewLang()
	env.ImportEval()
	sypext.ImportRandom(env)
	sypext.ImportTime(env)
	sypext.ImportChannels(env)
	sypext.ImportCoroutines(env)
	sypext.ImportRegex(env)
	ry.RegisterToEnv(env)

	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		runScript(env, args[0])
	} else {
		repl(env)
	}
}
