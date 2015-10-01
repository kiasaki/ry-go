package lang

import (
	"errors"
	"strconv"
)

func AssertArgsCount(what string, args []Value, expectedCount int) error {
	if len(args) != expectedCount {
		expectedStr := strconv.FormatInt(int64(expectedCount), 10)
		gotStr := strconv.FormatInt(int64(len(args)), 10)
		return errors.New(what + " takes exactly " + expectedStr + " params. " + gotStr + " given")
	}
	return nil
}

func AssertArgsMinCount(what string, args []Value, expectedCount int) error {
	if len(args) < expectedCount {
		expectedStr := strconv.FormatInt(int64(expectedCount), 10)
		gotStr := strconv.FormatInt(int64(len(args)), 10)
		return errors.New(what + " takes minimum " + expectedStr + " params. " + gotStr + " given")
	}
	return nil
}

func AssertType(what string, argName string, arg Value, expectedType ValueType) error {
	if arg.Type() != expectedType {
		return errors.New(what + " expected argument " + argName + " to be of type " + TypeName(expectedType) + ". " + TypeName(arg.Type()) + " '" + arg.String() + "' given")
	}
	return nil
}

func TypeName(typ ValueType) string {
	switch typ {
	case V_SYMBOL:
		return "symbol"
	case V_STRING:
		return "string"
	case V_INTEGER:
		return "integer"
	case V_FLOAT:
		return "float"
	case V_CHAR:
		return "char"
	case V_BOOL:
		return "bool"
	case V_LIST:
		return "list"
	case V_FUNC:
		return "func"
	case V_MACRO:
		return "macro"
	default:
		return "unknown"
	}
}
