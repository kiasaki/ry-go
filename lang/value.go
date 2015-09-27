package lang

import (
	"errors"
	"strconv"
)

type ValueType int

const (
	V_SYMBOL  ValueType = iota
	V_STRING            = iota
	V_INTEGER           = iota
	V_FLOAT             = iota
	V_CHAR              = iota
	V_BOOL              = iota
	V_LIST              = iota
	V_FUNC              = iota
)

type Value interface {
	Type() ValueType
	String() string
}

/* SYMBOL */
type SymbolValue struct {
	value string
}

func (SymbolValue) Type() ValueType {
	return V_SYMBOL
}
func (v SymbolValue) String() string {
	return v.value
}

/* STRING */
type StringValue struct {
	value string
}

func (StringValue) Type() ValueType {
	return V_STRING
}
func (v StringValue) String() string {
	return "\"" + v.value + "\""
}

/* INTEGER */
type IntegerValue struct {
	value int64
}

func (IntegerValue) Type() ValueType {
	return V_INTEGER
}
func (v IntegerValue) String() string {
	return strconv.FormatInt(v.value, 10)
}

/* FLOAT */
type FloatValue struct {
	value float64
}

func (FloatValue) Type() ValueType {
	return V_FLOAT
}
func (v FloatValue) String() string {
	return strconv.FormatFloat(v.value, 'f', -1, 64)
}

/* CHAR */
type CharValue struct {
	value rune
}

func (CharValue) Type() ValueType {
	return V_CHAR
}
func (v CharValue) String() string {
	return "'" + string(v.value) + "'"
}

/* BOOL */
type BoolValue struct {
	value bool
}

func (BoolValue) Type() ValueType {
	return V_BOOL
}
func (v BoolValue) String() string {
	if v.value {
		return "#t"
	}
	return "#f"
}

/* LIST */
type ListValue struct {
	childs []Value
}

func (ListValue) Type() ValueType {
	return V_LIST
}
func (v ListValue) String() string {
	str := "("
	for i, child := range v.childs {
		if i != 0 {
			str = str + " "
		}
		str = str + child.String()
	}
	return str + ")"
}

func NewEmptyListValue() ListValue {
	return ListValue{childs: []Value{}}
}

/* FUNC */
type FuncValue struct {
	name string
	fn   func([]Value) (Value, error)
	env  *Env
}

func (FuncValue) Type() ValueType {
	return V_FUNC
}
func (v FuncValue) String() string {
	return "<fn:" + v.name + ">"
}

func (v *FuncValue) Call(args []Value, parentEnv *Env) (Value, error) {

}

func NewFunction(name string, fn func([]Value) (Value, error), argsNames Value) (*FuncValue, error) {
	if name == "" {
		name = "*lambda*"
	}

	env := NewRootEnv()
	for _, arg := range argsNames {
		if arg.Type() == V_SYMBOL {
			env.Set()
		} else {
			return nil, errors.New("Can't define function '" + name + "' because param " + arg.String() + " is not a symbol")
		}
	}

	return &FuncValue{
		name: name,
		fn:   fn,
		env:  env,
	}, nil
}
