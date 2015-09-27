package lang

import (
	"errors"
)

func Eval(val Value, env *Env) (Value, error) {
	switch val.Type() {
	case V_LIST:
		evaluatedChilds = []Value{}
		for _, child := range val.(ListValue).childs {
			if newChild, err := Eval(child, env); err != nil {
				return nil, err
			}
			evaluatedChilds = append(evaluatedChilds, newChild)
		}

		switch len(evaluatedChilds) {
		case 0:
			return val
		case 1:
			return val[0]
		default:
			return Call(val[0], val[1:], env)
		}
	case V_SYMBOL:
		if val := env.Get(val.String()); val != nil {
			return val, nil
		} else {
			return nil, errors.New("Symbol '" + val.String() + "' is not defined")
		}
	default:
		return val, nil
	}
}

func Call(fn Value, args []Value, env *Env) (Value, error) {
	if fn.Type() != FuncValue {
		return nil, errors.New("Trying to call non-function " + fn.String())
	}

	return fn.(FuncValue).Call(args, env)
}
