package ry

import (
	"fmt"
	"os"
	"time"

	"github.com/jonvaldes/termo"
	"github.com/kiasaki/ry/lang"
)

const RUNTIME_FILE = "runtime.ryl"

var editor *Editor

type Editor struct {
	width       int
	height      int
	framebuffer *termo.Framebuffer

	env *lang.Env

	keyChan chan (termo.ScanCode)
	errChan chan (error)

	ticker <-chan (time.Time)
}

func (e *Editor) Start() {
	if err := termo.Init(); err != nil {
		panic(err)
	}
	termo.ShowCursor()

	defer func() {
		termo.Stop()
		if err := recover(); err != nil {
			// debug panic(err)
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	// Get size
	var err error
	e.width, e.height, err = termo.Size()
	if err != nil {
		panic(err)
	}

	// Create frame buffer
	e.framebuffer = termo.NewFramebuffer(e.width, e.height)

	e.startLispRuntime()

	e.startKeyReadLoop()

	e.MainLoop()
}

func (e *Editor) startKeyReadLoop() {
	// Read keys
	e.keyChan = make(chan termo.ScanCode, 100)
	e.errChan = make(chan error)
	termo.StartKeyReadLoop(e.keyChan, e.errChan)
	e.ticker = time.Tick(100 * time.Millisecond)
}

func (e *Editor) startLispRuntime() {
	// Create environment
	e.env = lang.NewBuiltinFilledEnv()

	// Load runtime
	if exprs, err := lang.ParseFile(RUNTIME_FILE); err != nil {
		die(err)
	} else {
		e.EvalLispExpressions(exprs)
	}

	// Register native Go functions
	e.registerLispFunctions()

	// Call init
	e.EvalLisp("(editor-initialize)")
}

func (e *Editor) EvalLisp(code string) {
	if exprs, err := lang.Parse([]byte(code)); err != nil {
		die(err)
	} else {
		e.EvalLispExpressions(exprs)
	}
}

func (e *Editor) EvalLispExpressions(exprs []lang.Value) {
	for _, expr := range exprs {
		if _, err := lang.Eval(expr, e.env); err != nil {
			die(err)
		}
	}
}

func (e *Editor) defineLispFunc(name string, argNames []string, fn func(*lang.Env, []lang.Value) (lang.Value, error)) {
	e.EvalLispExpressions([]lang.Value{lang.ListValue{[]lang.Value{
		lang.SymbolValue{"define"},
		lang.SymbolValue{name},
		lang.FuncValue{
			Name:     name,
			ArgNames: argNames,
			Fn:       fn,
		},
	}}})
}

func (e *Editor) MainLoop() {
	for {
		// Check for terminal resize
		if _w, _h, _ := termo.Size(); e.width != _w || e.height != _h {
			e.width = _w
			e.height = _h
			e.framebuffer = termo.NewFramebuffer(e.width, e.height)
		}

		// Clear framebuffer
		e.framebuffer.Clear()

		e.EvalLisp("(editor-render)")

		// Read keyboard
		select {
		case <-e.ticker:
			// Periodically flush framebuffer to screen
			e.framebuffer.Flush()
		case s := <-e.keyChan:
			e.handleScanCode(s)
		case err := <-e.errChan:
			die(err)
		}
	}
}

func (e *Editor) Render() {
	e.framebuffer.AttribRect(0, 0, e.width, e.height-2, termo.CellState{
		termo.AttrNone,
		termo.ColorGray,
		termo.ColorBlack,
	})
}

func (e *Editor) handleScanCode(s termo.ScanCode) {
	if s.IsMouseMoveEvent() || s.IsMouseDownEvent() || s.IsMouseUpEvent() {
		// don't bother
	} else if s.IsEscapeCode() {
		switch s.EscapeCode() {
		case 65: // Up
		case 66: // Down
		case 67: // Right
		case 68: // Left
		}
	} else {
		r := s.Rune()

		if r == 3 {
			// Exit if Ctrl+C is pressed
			e.Quit()
		} else {
			e.EvalLispExpressions([]lang.Value{lang.ListValue{[]lang.Value{
				lang.SymbolValue{"editor-handle-keypress"},
				lang.CharValue{r},
			}}})
		}
	}
}

func (e *Editor) Quit() {
	e.framebuffer.Clear()
	e.framebuffer.Flush()
	termo.Stop()
	os.Exit(0)
}

func StartMainEditor() {
	editor = &Editor{}
	editor.Start()
}
