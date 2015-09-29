package lang

import (
	"errors"
)

func Eval(val Value, env *Env) (Value, error) {
	switch val.Type() {
	case V_LIST:
		evaluatedChilds := []Value{}
		for _, child := range val.(ListValue).childs {
			if newChild, err := Eval(child, env); err != nil {
				return nil, err
			} else {
				evaluatedChilds = append(evaluatedChilds, newChild)
			}
		}

		switch len(evaluatedChilds) {
		case 0:
			return val, nil
		case 1:
			return evaluatedChilds[0], nil
		default:
			// TODO Check if callee is macro, if so, don't eval childs
			return Call(evaluatedChilds[0], evaluatedChilds[1:], env)
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
	if fn.Type() != V_FUNC {
		return nil, errors.New("Trying to call non-function " + fn.String())
	}

	callee := fn.(*FuncValue)
	return callee.Call(args, env)
}
