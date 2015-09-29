package lang

import (
	"errors"
)

func Eval(val Value, env *Env) (Value, error) {
	switch val.Type() {
	case V_LIST:
		childExprs := val.(ListValue).Childs
		evaluatedChilds := []Value{}

		if len(childExprs) == 0 {
			return val, nil
		}

		if evaluatedChild, err := Eval(childExprs[0], env); err != nil {
			return nil, err
		} else {
			evaluatedChilds = append(evaluatedChilds, evaluatedChild)
		}

		if evaluatedChilds[0].Type() == V_MACRO {
			return Call(evaluatedChilds[0], childExprs[1:], env)
		}

		// Now macros are out of the way, let's evaluate what's left
		for _, child := range childExprs[1:] {
			if evaluatedChild, err := Eval(child, env); err != nil {
				return nil, err
			} else {
				evaluatedChilds = append(evaluatedChilds, evaluatedChild)
			}
		}
		if evaluatedChilds[0].Type() == V_FUNC {
			return Call(evaluatedChilds[0], evaluatedChilds[1:], env)
		} else if len(evaluatedChilds) == 1 {
			return evaluatedChilds[0], nil
		} else {
			return ListValue{evaluatedChilds}, nil
		}
	case V_SYMBOL:
		if evVal := env.Get(val.String()); evVal != nil {
			return evVal, nil
		} else {
			return nil, errors.New("Symbol '" + val.String() + "' is not defined")
		}
	default:
		return val, nil
	}
}

func Call(fn Value, args []Value, env *Env) (Value, error) {
	if fn.Type() != V_FUNC && fn.Type() != V_MACRO {
		return nil, errors.New("Trying to call non-function " + fn.String())
	}

	if fn.Type() == V_FUNC {
		callee := fn.(*FuncValue)
		return callee.Call(args, env)
	} else {
		callee := fn.(*MacroValue)
		return callee.Call(args, env)
	}
}
