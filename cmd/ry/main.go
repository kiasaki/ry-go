package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kiasaki/glisp/extensions"
	"github.com/kiasaki/glisp/interpreter"
	ry "github.com/kiasaki/ry"
)

const LOCAL_DOT_FILE = ".init.ryl"

func quitOnErr(message string, err error) {
	if err != nil {
		fmt.Println(message)
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	env := glisp.NewGlisp()
	env.ImportEval()
	glispext.ImportRandom(env)
	glispext.ImportTime(env)
	glispext.ImportChannels(env)
	glispext.ImportCoroutines(env)
	glispext.ImportRegex(env)
	ry.RegisterToEnv(env)

	args := flag.Args()
	if len(args) > 0 {
		// TODO open buffers with those files
		//runScript(env, args[0])
	}

	// load ry's lisp runtime
	for _, fileName := range ry.AssetNames() {
		err := env.LoadString(string(ry.MustAsset(fileName)))
		quitOnErr("Failed to load runtime", err)
	}

	// load current folder's .init.ryl
	file, err := os.Open(LOCAL_DOT_FILE)
	if err != nil && !os.IsNotExist(err) {
		quitOnErr("Error reading local "+LOCAL_DOT_FILE, err)
	} else if !os.IsNotExist(err) {
		defer file.Close()
		err = env.LoadFile(file)
		quitOnErr("Error parsing local "+LOCAL_DOT_FILE, err)
	}

	// hang till exit as last expr
	// coroutines will make this hanging in the main thred work
	err = env.LoadString("(<! quit-chan)")
	quitOnErr("Error evaluating quit signal", err)

	_, err = env.Run()
	quitOnErr(env.GetStackTrace(err), err)
}
