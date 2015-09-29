package lang

import (
	"errors"
)

var builtinType *FuncValue = &FuncValue{
	name:     "type",
	argNames: []string{"value"},
	fn: func(env *Env, args []Value) (Value, error) {
		return nil, nil
	},
}

var builtinEq *FuncValue

func init() {
	builtinEq = &FuncValue{
		name:     "eq?",
		argNames: []string{"value"},
		fn: func(env *Env, args []Value) (Value, error) {
			if len(args) != 2 {
				return nil, errors.New("Builtin 'eq?' must be called with 2 params, " + string(len(args)) + " passed")
			}
			if args[0].Type() != args[1].Type() {
				return BoolValueFalse, nil
			}

			switch args[0].Type() {
			case V_SYMBOL:
				return BoolValue{args[0].(SymbolValue).value == args[1].(SymbolValue).value}, nil
			case V_STRING:
				return BoolValue{args[0].(StringValue).value == args[1].(StringValue).value}, nil
			case V_INTEGER:
				return BoolValue{args[0].(IntegerValue).value == args[1].(IntegerValue).value}, nil
			case V_FLOAT:
				return BoolValue{args[0].(FloatValue).value == args[1].(FloatValue).value}, nil
			case V_CHAR:
				return BoolValue{args[0].(CharValue).value == args[1].(CharValue).value}, nil
			case V_BOOL:
				return BoolValue{args[0].(BoolValue).value == args[1].(BoolValue).value}, nil
			case V_LIST:
				if len(args[0].(ListValue).childs) != len(args[1].(ListValue).childs) {
					return BoolValueFalse, nil
				}
				equal := true
				for _, leftVal := range args[0].(ListValue).childs {
					for _, rightVal := range args[1].(ListValue).childs {
						if val, err := builtinEq.fn(env, []Value{leftVal, rightVal}); err != nil {
							return nil, err
						} else {
							if !val.(BoolValue).value {
								equal = false
							}
						}
					}
				}
				return BoolValue{equal}, nil
			case V_FUNC:
				return BoolValue{args[0].(FuncValue).name == args[1].(FuncValue).name}, nil
			}

			return BoolValueFalse, nil
		},
	}
}
