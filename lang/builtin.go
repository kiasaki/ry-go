package lang

import (
	"errors"
)

var builtinType *FuncValue = &FuncValue{
	Name:     "type",
	ArgNames: []string{"value"},
	Fn: func(env *Env, args []Value) (Value, error) {
		if err := AssertArgsCount("Builtin 'type'", args, 1); err != nil {
			return nil, err
		}
		return SymbolValue{TypeName(args[0].Type())}, nil
	},
}

var builtinEq *FuncValue
var builtinList *FuncValue
var builtinString *FuncValue
var builtinDefine *MacroValue
var builtinSet *MacroValue
var builtinLambda *MacroValue

func init() {
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

	builtinDefine = &MacroValue{
		Name:     "define",
		ArgNames: []string{"name", "value"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsCount("Builtin 'define'", args, 2); err != nil {
				return nil, err
			}
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
		Name:     "set",
		ArgNames: []string{"name", "value"},
		Fn: func(env *Env, args []Value) (Value, error) {
			if err := AssertArgsCount("Builtin 'set'", args, 2); err != nil {
				return nil, err
			}
			if err := AssertType("Builtin 'set'", "1", args[0], V_SYMBOL); err != nil {
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
}
