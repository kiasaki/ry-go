package main

import (
	"flag"
	"fmt"
	"os"

	ry "github.com/kiasaki/ry"
	sypext "github.com/kiasaki/syp-lang/extensions"
	syp "github.com/kiasaki/syp-lang/interpreter"
	"github.com/nsf/termbox-go"
)

const LOCAL_DOT_FILE = ".init.syp"

func quitOnErr(message string, err error) {
	if err != nil {
		termbox.Close()
		fmt.Println(message)
		fmt.Println(err)
		os.Exit(1)
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

	// load current folder's .init.syp
	file, err := os.Open(LOCAL_DOT_FILE)
	if err != nil && !os.IsNotExist(err) {
		quitOnErr("Error reading local "+LOCAL_DOT_FILE, err)
	} else if !os.IsNotExist(err) {
		defer file.Close()
		err = env.LoadFile(file)
		quitOnErr("Error parsing local "+LOCAL_DOT_FILE, err)
	}

	// hang till exit as last expr
	// coroutines will make this hang in the main thread util `quit` is called
	err = env.LoadString("(<! quit-chan)")
	quitOnErr("Error evaluating quit signal", err)

	_, err = env.Run()
	quitOnErr(env.GetStackTrace(err), err)
}
