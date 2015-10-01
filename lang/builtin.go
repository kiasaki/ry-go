package lang

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

var builtinType *FuncValue
var builtinEq *FuncValue
var builtinCar *FuncValue
var builtinCdr *FuncValue
var builtinCons *FuncValue
var builtinAppend *FuncValue
var builtinLength *FuncValue
var builtinListRef *FuncValue
var builtinListSet *FuncValue
var builtinStringRegexpSplit *FuncValue
var builtinStringLength *FuncValue

var builtinList *FuncValue
var builtinString *FuncValue

var builtinQuote *MacroValue
var builtinUnquote *MacroValue
var builtinEval *MacroValue
var builtinDefine *MacroValue
var builtinSet *MacroValue
var builtinLet *MacroValue
var builtinLetStar *MacroValue
var builtinLetRec *MacroValue
var builtinLetRecStar *MacroValue
var builtinLambda *MacroValue
var builtinDefmacro *MacroValue

var builtinIf *MacroValue
var builtinCond *MacroValue
var builtinBegin *FuncValue

var builtinError *FuncValue
var builtinRead *FuncValue
var builtinWrite *FuncValue

func buildLetFunc(letName string, preComputedValues, preBoundValues bool) func(*Env, []Value) (Value, error) {
	return func(env *Env, args []Value) (Value, error) {
		name := "Builtin '" + letName + "'"
		if err := AssertArgsMinCount(name, args, 2); err != nil {
			return nil, err
		}

		letEnv := NewRootEnv()
		letEnv.Parent = env
		var defs []Value
		var body []Value

		// handle named let
		if args[0].Type() == V_SYMBOL {
			if err := AssertArgsMinCount(name+" (in named form)", args, 3); err != nil {
				return nil, err
			}
			body = args[2:]
			if err := AssertType(name, "defs", args[1], V_LIST); err != nil {
				return nil, err
			}
			defs = args[1].(ListValue).Childs

			letEnv.Set(args[0].(SymbolValue).Value, &FuncValue{
				Name:     args[0].(SymbolValue).Value,
				ArgNames: []string{".", "vals"},
				Fn: func(e *Env, args []Value) (Value, error) {
					for i, arg := range args {
						if i < len(defs) {
							// we can asume it's a symbol as there is a check
							// down there and this func wont be called before it
							e.Set(defs[i].(SymbolValue).Value, arg)
						}
					}

					var lastResult Value
					var err error
					for _, expr := range body {
						if lastResult, err = Eval(expr, letEnv); err != nil {
							return nil, err
						}
					}
					return lastResult, nil
				},
			})
		} else {
			body = args[1:]
			if err := AssertType(name, "defs", args[0], V_LIST); err != nil {
				return nil, err
			}
			defs = args[0].(ListValue).Childs
		}

		// extract defs and body
		for i, def := range defs {
			if err := AssertType(name, "defs["+strconv.FormatInt(int64(i), 10)+"]", def, V_LIST); err != nil {
				return nil, err
			}
			if err := AssertType(name, "defs["+strconv.FormatInt(int64(i), 10)+"][0]", def.(ListValue).Childs[0], V_SYMBOL); err != nil {
				return nil, err
			}
		}

		// pre bind symbols for recursive use (rec)
		if preBoundValues {
			for _, def := range defs {
				castedDef := def.(ListValue).Childs
				letEnv.Set(castedDef[0].(SymbolValue).Value, NewEmptyListValue())
			}
		}

		// eval assignments
		evaledResults := map[string]Value{}
		for _, def := range defs {
			castedDef := def.(ListValue).Childs
			castedDefName := castedDef[0].(SymbolValue).Value

			result, err := Eval(castedDef[1], letEnv)
			if err != nil {
				return nil, err
			}

			if preComputedValues {
				evaledResults[castedDefName] = result
			} else {
				letEnv.Set(castedDefName, result)
			}
		}

		// assign definition values if not done during evaluation
		if !preComputedValues {
			for key, val := range evaledResults {
				letEnv.Set(key, val)
			}
		}

		// eval body
		var lastResult Value
		var err error
		for _, expr := range body {
			if lastResult, err = Eval(expr, letEnv); err != nil {
				return nil, err
			}
		}
		return lastResult, nil
	}
}

func init() {
	builtinType = &FuncValue{
		Name:     "type",
		ArgNames: []string{"value"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsCount("Builtin 'type'", args, 1); err != nil {
				return nil, err
			}
			return SymbolValue{TypeName(args[0].Type())}, nil
		},
	}
	builtinEq = &FuncValue{
		Name:     "eq?",
		ArgNames: []string{"value"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsCount("Builtin 'eq?'", args, 2); err != nil {
				return nil, err
			}
			if args[0].Type() != args[1].Type() {
				return BoolValueFalse, nil
			}

			switch args[0].Type() {
			case V_SYMBOL:
				return BoolValue{args[0].(SymbolValue).Value == args[1].(SymbolValue).Value}, nil
			case V_STRING:
				return BoolValue{args[0].(StringValue).Value == args[1].(StringValue).Value}, nil
			case V_INTEGER:
				return BoolValue{args[0].(IntegerValue).Value == args[1].(IntegerValue).Value}, nil
			case V_FLOAT:
				return BoolValue{args[0].(FloatValue).Value == args[1].(FloatValue).Value}, nil
			case V_CHAR:
				return BoolValue{args[0].(CharValue).Value == args[1].(CharValue).Value}, nil
			case V_BOOL:
				return BoolValue{args[0].(BoolValue).Value == args[1].(BoolValue).Value}, nil
			case V_LIST:
				if len(args[0].(ListValue).Childs) != len(args[1].(ListValue).Childs) {
					return BoolValueFalse, nil
				}
				equal := true
				for _, leftVal := range args[0].(ListValue).Childs {
					for _, rightVal := range args[1].(ListValue).Childs {
						if val, err := builtinEq.Fn(env, []Value{leftVal, rightVal}); err != nil {
							return nil, err
						} else {
							if !val.(BoolValue).Value {
								equal = false
							}
						}
					}
				}
				return BoolValue{equal}, nil
			case V_FUNC:
				return BoolValue{args[0].(*FuncValue).Name == args[1].(*FuncValue).Name}, nil
			case V_MACRO:
				return BoolValue{args[0].(*MacroValue).Name == args[1].(*MacroValue).Name}, nil
			}

			return BoolValueFalse, nil
		},
	}
	builtinCar = &FuncValue{
		Name:     "car",
		ArgNames: []string{"lst"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsCount("Builtin 'car'", args, 1); err != nil {
				return nil, err
			}
			if err := AssertType("Builtin 'car'", "1", args[0], V_LIST); err != nil {
				return nil, err
			}
			return args[0].(ListValue).Childs[0], nil
		},
	}
	builtinCdr = &FuncValue{
		Name:     "cdr",
		ArgNames: []string{"lst"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsCount("Builtin 'cdr'", args, 1); err != nil {
				return nil, err
			}
			if err := AssertType("Builtin 'cdr'", "1", args[0], V_LIST); err != nil {
				return nil, err
			}
			return ListValue{args[0].(ListValue).Childs[1:]}, nil
		},
	}
	builtinCons = &FuncValue{
		Name:     "cons",
		ArgNames: []string{"value", "lst"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsCount("Builtin 'cons'", args, 2); err != nil {
				return nil, err
			}
			if err := AssertType("Builtin 'cons'", "2", args[1], V_LIST); err != nil {
				return nil, err
			}

			return ListValue{append([]Value{args[0]}, args[1].(ListValue).Childs...)}, nil
		},
	}
	builtinAppend = &FuncValue{
		Name:     "append",
		ArgNames: []string{"value", "lst"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsMinCount("Builtin 'append'", args, 2); err != nil {
				return nil, err
			}

			newList := []Value{}

			for i, arg := range args {
				if err := AssertType("Builtin 'append'", strconv.FormatInt(int64(i), 10), arg, V_LIST); err != nil {
					return nil, err
				}

				newList = append(newList, arg.(ListValue).Childs...)
			}

			return ListValue{newList}, nil
		},
	}
	builtinLength = &FuncValue{
		Name:     "length",
		ArgNames: []string{"lst"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsCount("Builtin 'length'", args, 1); err != nil {
				return nil, err
			}
			if err := AssertType("Builtin 'length'", "1", args[0], V_LIST); err != nil {
				return nil, err
			}

			return IntegerValue{int64(len(args[0].(ListValue).Childs))}, nil
		},
	}
	builtinListRef = &FuncValue{
		Name:     "list-ref",
		ArgNames: []string{"lst", "i"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsCount("Builtin 'list-ref'", args, 2); err != nil {
				return nil, err
			}
			if err := AssertType("Builtin 'list-ref'", "1", args[0], V_LIST); err != nil {
				return nil, err
			}
			if err := AssertType("Builtin 'list-ref'", "2", args[1], V_INTEGER); err != nil {
				return nil, err
			}

			i := args[1].(IntegerValue).Value
			return args[0].(ListValue).Childs[i], nil
		},
	}
	builtinListSet = &FuncValue{
		Name:     "list-set!",
		ArgNames: []string{"lst", "i", "val"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsCount("Builtin 'list-set!'", args, 3); err != nil {
				return nil, err
			}
			if err := AssertType("Builtin 'list-set!'", "1", args[0], V_LIST); err != nil {
				return nil, err
			}
			if err := AssertType("Builtin 'list-set!'", "2", args[1], V_INTEGER); err != nil {
				return nil, err
			}

			list := args[0].(ListValue)
			i := args[1].(IntegerValue).Value
			list.Childs[i] = args[2]
			return list, nil
		},
	}
	builtinStringRegexpSplit = &FuncValue{
		Name:     "string-regexp-split",
		ArgNames: []string{"regexp", "str"},
		Fn: func(env *Env, args []Value) (Value, error) {
			name := "Builtin 'string-regexp-split'"
			if err := AssertArgsCount(name, args, 2); err != nil {
				return nil, err
			}
			if err := AssertType(name, "regexp", args[0], V_STRING); err != nil {
				return nil, err
			}
			if err := AssertType(name, "str", args[1], V_STRING); err != nil {
				return nil, err
			}

			regexpStr := args[0].(StringValue).Value
			str := args[1].(StringValue).Value
			if r, err := regexp.Compile(regexpStr); err != nil {
				return nil, err
			} else {
				parts := r.Split(str, -1)
				lispParts := []Value{}
				for _, part := range parts {
					lispParts = append(lispParts, StringValue{part})
				}
				return ListValue{lispParts}, nil
			}
		},
	}
	builtinStringLength = &FuncValue{
		Name:     "string-length",
		ArgNames: []string{"str"},
		Fn: func(env *Env, args []Value) (Value, error) {
			name := "Builtin 'string-length'"
			if err := AssertArgsCount(name, args, 1); err != nil {
				return nil, err
			}
			if err := AssertType(name, "str", args[0], V_STRING); err != nil {
				return nil, err
			}

			return IntegerValue{int64(len([]rune(args[0].(StringValue).Value)))}, nil
		},
	}

	builtinList = &FuncValue{
		Name:     "list",
		ArgNames: []string{".", "values"},
		Fn: func(env *Env, args []Value) (Value, error) {
			return ListValue{args}, nil
		},
	}

	// Create a new string, can be passed strings or chars and will join them
	builtinString = &FuncValue{
		Name:     "string",
		ArgNames: []string{".", "values"},
		Fn: func(env *Env, args []Value) (Value, error) {
			result := ""
			for _, val := range args {
				if val.Type() == V_STRING {
					result = result + val.(StringValue).Value
				} else if val.Type() == V_CHAR {
					result = result + string(val.(CharValue).Value)
				} else {
					return nil, errors.New("Builtin 'string' called with non string or char argument")
				}
			}
			return StringValue{result}, nil
		},
	}

	builtinQuote = &MacroValue{
		Name:     "quote",
		ArgNames: []string{"value"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsCount("Builtin 'quote'", args, 1); err != nil {
				return nil, err
			}

			return args[0], nil
		},
	}
	builtinUnquote = &MacroValue{
		Name:     "unquote",
		ArgNames: []string{"value"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsCount("Builtin 'unquote'", args, 1); err != nil {
				return nil, err
			}

			return Eval(args[0], env)
		},
	}
	builtinEval = &MacroValue{
		Name:     "eval",
		ArgNames: []string{"value"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsCount("Builtin 'eval'", args, 1); err != nil {
				return nil, err
			}

			return Eval(args[0], env)
		},
	}
	builtinDefine = &MacroValue{
		Name:     "define",
		ArgNames: []string{"name", ".", "values"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsMinCount("Builtin 'define'", args, 1); err != nil {
				return nil, err
			}

			// Check if params #1 is list, if so, define function instead
			if args[0].Type() == V_LIST {
				childs := args[0].(ListValue).Childs
				if len(childs) < 1 {
					return nil, errors.New("Builtin 'define' parameter '1' needs at least one item in it")
				}
				for i, child := range childs {
					if err := AssertType("Builtin 'define'", "formals["+strconv.FormatInt(int64(i), 10)+"]", child, V_SYMBOL); err != nil {
						return nil, err
					}
				}

				return Eval(ListValue{[]Value{
					SymbolValue{"define"},
					childs[0],
					ListValue{append([]Value{
						SymbolValue{"lambda"},
						ListValue{childs[1:]},
					}, args[1:]...)},
				}}, env)
			}

			// Else assign value in root env
			if err := AssertType("Builtin 'define'", "1", args[0], V_SYMBOL); err != nil {
				return nil, err
			}

			if evaledValue, err := Eval(args[1], env); err != nil {
				return nil, err
			} else {
				env.SetOnRoot(args[0].(SymbolValue).Value, evaledValue)
				return NewEmptyListValue(), nil
			}
		},
	}
	builtinSet = &MacroValue{
		Name:     "set!",
		ArgNames: []string{"name", "value"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsCount("Builtin 'set!'", args, 2); err != nil {
				return nil, err
			}
			if err := AssertType("Builtin 'set!'", "1", args[0], V_SYMBOL); err != nil {
				return nil, err
			}

			if evaledValue, err := Eval(args[1], env); err != nil {
				return nil, err
			} else {
				env.SetOnRoot(args[0].(SymbolValue).Value, evaledValue)
				return NewEmptyListValue(), nil
			}
		},
	}
	builtinLet = &MacroValue{
		Name:     "let",
		ArgNames: []string{"defs", ".", "body"},
		Fn:       buildLetFunc("let", false, false),
	}
	builtinLetStar = &MacroValue{
		Name:     "let*",
		ArgNames: []string{"defs", ".", "body"},
		Fn:       buildLetFunc("let*", true, false),
	}
	builtinLetRec = &MacroValue{
		Name:     "letrec",
		ArgNames: []string{"defs", ".", "body"},
		Fn:       buildLetFunc("letrec", false, true),
	}
	builtinLetRecStar = &MacroValue{
		Name:     "letrec*",
		ArgNames: []string{"defs", ".", "body"},
		Fn:       buildLetFunc("letrec*", true, true),
	}
	builtinLambda = &MacroValue{
		Name:     "lambda",
		ArgNames: []string{"args", ".", "body"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsMinCount("Builtin 'lambda'", args, 1); err != nil {
				return nil, err
			}
			if err := AssertType("Builtin 'lambda'", "1", args[0], V_LIST); err != nil {
				return nil, err
			}

			body := args[1:]

			return NewFunction("", func(env *Env, args []Value) (Value, error) {
				var lastResult Value
				var err error
				// TODO handle tail recursion?
				for _, expr := range body {
					if lastResult, err = Eval(expr, env); err != nil {
						return nil, err
					}
				}
				return lastResult, nil
			}, args[0])
		},
	}
	builtinDefmacro = &MacroValue{
		Name:     "defmacro",
		ArgNames: []string{"args", ".", "body"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsMinCount("Builtin 'defmacro'", args, 1); err != nil {
				return nil, err
			}
			if err := AssertType("Builtin 'defmacro'", "1", args[0], V_LIST); err != nil {
				return nil, err
			}

			definition := args[0].(ListValue).Childs
			body := args[1:]

			if err := AssertType("Builtin 'defmacro'", "name", definition[0], V_SYMBOL); err != nil {
				return nil, err
			}

			macro, err := NewMacro(definition[0].(SymbolValue).Value, func(env *Env, args []Value) (Value, error) {
				var lastResult Value
				var err error
				// TODO handle tail recursion?
				for _, expr := range body {
					if lastResult, err = Eval(expr, env); err != nil {
						return nil, err
					}
				}
				return lastResult, nil
			}, ListValue{definition[1:]})
			if err == nil {
				env.SetOnRoot(definition[0].(SymbolValue).Value, macro)
			}
			return macro, err
		},
	}

	builtinIf = &MacroValue{
		Name:     "if",
		ArgNames: []string{"cond", ".", "body"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsMinCount("Builtin 'if'", args, 2); err != nil {
				return nil, err
			}

			condResult, err := Eval(args[0], env)
			if err != nil {
				return nil, err
			}
			if err = AssertType("Builtin 'if'", "1", condResult, V_BOOL); err != nil {
				return nil, err
			}

			if condResult.(BoolValue).Value {
				return Eval(args[1], env)
			} else if len(args) > 2 {
				return Eval(args[2], env)
			} else {
				return NewEmptyListValue(), nil
			}
		},
	}
	builtinCond = &MacroValue{
		Name:     "cond",
		ArgNames: []string{".", "conds"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsMinCount("Builtin 'cond'", args, 1); err != nil {
				return nil, err
			}

			for _, arg := range args {
				if err := AssertType("Builtin 'cond'", "cond", arg, V_LIST); err != nil {
					return nil, err
				}
				condResult, err := Eval(arg.(ListValue).Childs[0], env)
				if err != nil {
					return nil, err
				}
				if err = AssertType("Builtin 'cond'", "cond", condResult, V_BOOL); err != nil {
					return nil, err
				}
				if condResult.(BoolValue).Value {
					return Eval(ListValue{
						append([]Value{SymbolValue{"begin"}}, arg.(ListValue).Childs[1:]...),
					}, env)
				}
			}

			return NewEmptyListValue(), nil
		},
	}
	builtinBegin = &FuncValue{
		Name:     "begin",
		ArgNames: []string{".", "exprs"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsMinCount("Builtin 'begin'", args, 1); err != nil {
				return nil, err
			}

			return args[len(args)-1], nil
		},
	}

	builtinError = &FuncValue{
		Name:     "error",
		ArgNames: []string{"message"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsCount("Builtin 'error'", args, 1); err != nil {
				return nil, err
			}
			if err := AssertType("Builtin 'error'", "1", args[0], V_STRING); err != nil {
				return nil, err
			}

			return nil, errors.New(args[0].(StringValue).Value)
		},
	}
	builtinRead = &FuncValue{
		Name:     "read",
		ArgNames: []string{"filename"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsCount("Builtin 'read'", args, 1); err != nil {
				return nil, err
			}
			if err := AssertType("Builtin 'read'", "1", args[0], V_STRING); err != nil {
				return nil, err
			}

			filename := args[0].(StringValue).Value
			if contents, err := ioutil.ReadFile(filename); err != nil {
				return NewEmptyListValue(), nil
			} else {
				return StringValue{string(contents)}, nil
			}
		},
	}
	builtinWrite = &FuncValue{
		Name:     "write",
		ArgNames: []string{"filename", "contents"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsCount("Builtin 'write'", args, 3); err != nil {
				return nil, err
			}
			if err := AssertType("Builtin 'write'", "1", args[0], V_STRING); err != nil {
				return nil, err
			}
			if err := AssertType("Builtin 'write'", "2", args[1], V_STRING); err != nil {
				return nil, err
			}
			if err := AssertType("Builtin 'write'", "3", args[2], V_INTEGER); err != nil {
				return nil, err
			}

			filename := args[0].(StringValue).Value
			contents := args[1].(StringValue).Value
			perm := os.FileMode(args[2].(IntegerValue).Value)
			if err := ioutil.WriteFile(filename, []byte(contents), perm); err != nil {
				return BoolValueFalse, nil
			} else {
				return BoolValueTrue, nil
			}
		},
	}
}
