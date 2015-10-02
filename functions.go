package ry

import (
	"errors"
	"strconv"

	"github.com/jonvaldes/termo"
	"github.com/kiasaki/ry/lang"
)

func (e *Editor) registerLispFunctions() {
	e.defineLispFunc(
		"die",
		[]string{"error"},
		func(env *lang.Env, args []lang.Value) (lang.Value, error) {
			if err := lang.AssertArgsCount("Function 'die'", args, 1); err != nil {
				return nil, err
			}
			if err := lang.AssertType("Function 'die'", "error", args[0], lang.V_STRING); err != nil {
				return nil, err
			}

			die(errors.New(args[0].(lang.StringValue).Value))
			return lang.NewEmptyListValue(), nil
		},
	)
	e.defineLispFunc(
		"quit",
		[]string{},
		func(env *lang.Env, args []lang.Value) (lang.Value, error) {
			e.Quit()
			return lang.NewEmptyListValue(), nil
		},
	)
	e.defineLispFunc(
		"editor-set-cursor",
		[]string{"x", "y"},
		func(env *lang.Env, args []lang.Value) (lang.Value, error) {
			name := "Function 'editor-set-cursor'"
			if err := lang.AssertArgsCount(name, args, 2); err != nil {
				return nil, err
			}
			if err := lang.AssertType(name, "x", args[0], lang.V_INTEGER); err != nil {
				return nil, err
			}
			if err := lang.AssertType(name, "y", args[1], lang.V_INTEGER); err != nil {
				return nil, err
			}

			x := int(args[0].(lang.IntegerValue).Value)
			y := int(args[1].(lang.IntegerValue).Value)
			termo.SetCursor(x, y)
			return lang.NewEmptyListValue(), nil
		},
	)
	e.defineLispFunc(
		"editor-height",
		[]string{},
		func(env *lang.Env, args []lang.Value) (lang.Value, error) {
			return lang.IntegerValue{int64(e.height)}, nil
		},
	)
	e.defineLispFunc(
		"editor-width",
		[]string{},
		func(env *lang.Env, args []lang.Value) (lang.Value, error) {
			return lang.IntegerValue{int64(e.width)}, nil
		},
	)
	e.defineLispFunc(
		"editor-draw-text",
		[]string{"x", "y", "style", "text"},
		func(env *lang.Env, args []lang.Value) (lang.Value, error) {
			if err := lang.AssertArgsCount("Function 'editor-draw-text'", args, 4); err != nil {
				return nil, err
			}
			if err := lang.AssertType("Function 'editor-draw-text'", "x", args[0], lang.V_INTEGER); err != nil {
				return nil, err
			}
			if err := lang.AssertType("Function 'editor-draw-text'", "y", args[1], lang.V_INTEGER); err != nil {
				return nil, err
			}
			if err := lang.AssertType("Function 'editor-draw-text'", "style", args[2], lang.V_LIST); err != nil {
				return nil, err
			}
			if err := lang.AssertType("Function 'editor-draw-text'", "text", args[3], lang.V_STRING); err != nil {
				return nil, err
			}
			cellStyleChilds := args[2].(lang.ListValue).Childs
			for i, child := range cellStyleChilds {
				if err := lang.AssertType("Function 'editor-draw-text'", "style["+strconv.FormatInt(int64(i), 10)+"]", child, lang.V_SYMBOL); err != nil {
					return nil, err
				}
			}

			x := int(args[0].(lang.IntegerValue).Value)
			y := int(args[1].(lang.IntegerValue).Value)
			e.framebuffer.AttribText(x, y, termo.CellState{
				stringToTerminalAttr(cellStyleChilds[0].(lang.SymbolValue).Value),
				stringToTerminalColor(cellStyleChilds[1].(lang.SymbolValue).Value),
				stringToTerminalColor(cellStyleChilds[2].(lang.SymbolValue).Value),
			}, args[3].(lang.StringValue).Value)

			return lang.NewEmptyListValue(), nil
		},
	)
	e.defineLispFunc(
		"editor-draw-attribute-rect",
		[]string{"x", "y", "x2", "y2", "style"},
		func(env *lang.Env, args []lang.Value) (lang.Value, error) {
			if err := lang.AssertArgsCount("Function 'editor-draw-attribute-rect'", args, 5); err != nil {
				return nil, err
			}
			for i := 0; i < 4; i++ {
				if err := lang.AssertType("Function 'editor-draw-attribute-rect'", strconv.FormatInt(int64(i+1), 10), args[i], lang.V_INTEGER); err != nil {
					return nil, err
				}
			}
			if err := lang.AssertType("Function 'editor-draw-attribute-rect'", "style", args[4], lang.V_LIST); err != nil {
				return nil, err
			}
			cellStyleChilds := args[4].(lang.ListValue).Childs
			for i, child := range cellStyleChilds {
				if err := lang.AssertType("Function 'editor-draw-attribute-rect'", "style["+strconv.FormatInt(int64(i), 10)+"]", child, lang.V_SYMBOL); err != nil {
					return nil, err
				}
			}

			x := int(args[0].(lang.IntegerValue).Value)
			y := int(args[1].(lang.IntegerValue).Value)
			x2 := int(args[2].(lang.IntegerValue).Value)
			y2 := int(args[3].(lang.IntegerValue).Value)
			e.framebuffer.AttribRect(x, y, x2, y2, termo.CellState{
				stringToTerminalAttr(cellStyleChilds[0].(lang.SymbolValue).Value),
				stringToTerminalColor(cellStyleChilds[1].(lang.SymbolValue).Value),
				stringToTerminalColor(cellStyleChilds[2].(lang.SymbolValue).Value),
			})

			return lang.NewEmptyListValue(), nil
		},
	)
}
