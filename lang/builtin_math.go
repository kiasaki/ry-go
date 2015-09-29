package lang

import (
	"errors"
)

func builtinOp(op string, args []Value) (Value, error) {
	if len(args) == 0 {
		return nil, errors.New("Builtin '" + op + "' can't be called without params")
	}
	if len(args) != 2 {
		return nil, errors.New("Builtin '" + op + "' can't needs exactly 2 params, " + string(len(args)) + " given")
	}
	left := args[0]
	right := args[1]

	if (left.Type() != V_INTEGER && left.Type() != V_FLOAT) || (right.Type() != V_INTEGER && right.Type() != V_FLOAT) {
		// type check
		return nil, errors.New("Arguments passed to '" + op + "' arn't of type int or float. " + left.String() + " and " + right.String() + " passed")
	}

	if left.Type() == V_FLOAT || right.Type() == V_FLOAT {
		// float op
		var leftVal, rightVal float64
		// make sure to case the integer if occuring
		if left.Type() == V_INTEGER {
			leftVal = float64(left.(IntegerValue).Value)
		} else {
			leftVal = left.(FloatValue).Value
		}
		if right.Type() == V_INTEGER {
			rightVal = float64(right.(IntegerValue).Value)
		} else {
			rightVal = right.(FloatValue).Value
		}

		switch op {
		case "+":
			return FloatValue{leftVal + rightVal}, nil
		case "-":
			return FloatValue{leftVal - rightVal}, nil
		case "*":
			return FloatValue{leftVal * rightVal}, nil
		case "/":
			if rightVal == 0 {
				return nil, errors.New("Can't divide " + left.String() + " by 0")
			}
			return FloatValue{leftVal / rightVal}, nil
		}
	} else {
		// int op
		leftVal := left.(IntegerValue).Value
		rightVal := right.(IntegerValue).Value

		switch op {
		case "+":
			return IntegerValue{leftVal + rightVal}, nil
		case "-":
			return IntegerValue{leftVal - rightVal}, nil
		case "*":
			return IntegerValue{leftVal * rightVal}, nil
		case "/":
			if rightVal == 0 {
				return nil, errors.New("Can't divide " + left.String() + " by 0")
			}
			return IntegerValue{leftVal / rightVal}, nil
		}
	}
	return nil, errors.New("Bad call to builtin_math.go:builtinOp, reached return")
}

func buildBuiltinOp(op string) *FuncValue {
	return &FuncValue{
		Name:     op,
		ArgNames: []string{"left", "right"},
		Fn: func(env *Env, args []Value) (Value, error) {
			return builtinOp(op, args)
		},
	}
}

var builtinAdd = buildBuiltinOp("+")
var builtinSubtract = buildBuiltinOp("-")
var builtinMultiply = buildBuiltinOp("*")
var builtinDivide = buildBuiltinOp("/")
