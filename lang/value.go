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
	V_MACRO             = iota
)

type Value interface {
	Type() ValueType
	String() string
}

/* SYMBOL */
type SymbolValue struct {
	Value string
}

func (SymbolValue) Type() ValueType {
	return V_SYMBOL
}
func (v SymbolValue) String() string {
	return v.Value
}

func NewSymbolValue(value string) SymbolValue {
	return SymbolValue{value}
}

/* STRING */
type StringValue struct {
	Value string
}

func (StringValue) Type() ValueType {
	return V_STRING
}
func (v StringValue) String() string {
	return "\"" + v.Value + "\""
}

func NewStringValue(value string) StringValue {
	return StringValue{value}
}

/* INTEGER */
type IntegerValue struct {
	Value int64
}

func (IntegerValue) Type() ValueType {
	return V_INTEGER
}
func (v IntegerValue) String() string {
	return strconv.FormatInt(v.Value, 10)
}

/* FLOAT */
type FloatValue struct {
	Value float64
}

func (FloatValue) Type() ValueType {
	return V_FLOAT
}
func (v FloatValue) String() string {
	return strconv.FormatFloat(v.Value, 'f', -1, 64)
}

/* CHAR */
type CharValue struct {
	Value rune
}

func (CharValue) Type() ValueType {
	return V_CHAR
}
func (v CharValue) String() string {
	return "'" + string(v.Value) + "'"
}

/* BOOL */
type BoolValue struct {
	Value bool
}

func (BoolValue) Type() ValueType {
	return V_BOOL
}
func (v BoolValue) String() string {
	if v.Value {
		return "#t"
	}
	return "#f"
}

var BoolValueTrue BoolValue = BoolValue{true}
var BoolValueFalse BoolValue = BoolValue{false}

/* LIST */
type ListValue struct {
	Childs []Value
}

func (ListValue) Type() ValueType {
	return V_LIST
}
func (v ListValue) String() string {
	str := "("
	for i, child := range v.Childs {
		if i != 0 {
			str = str + " "
		}
		str = str + child.String()
	}
	return str + ")"
}

func NewEmptyListValue() ListValue {
	return ListValue{Childs: []Value{}}
}

/* FUNC */
type FuncValue struct {
	Name     string
	ArgNames []string
	Fn       func(*Env, []Value) (Value, error)
}

func (FuncValue) Type() ValueType {
	return V_FUNC
}
func (v FuncValue) String() string {
	return "<fn:" + v.Name + ">"
}

func (v *FuncValue) Call(args []Value, parentEnv *Env) (Value, error) {
	if env, err := createEnvAndAssignArgs(v.Name, parentEnv, args, v.ArgNames); err != nil {
		return nil, err
	} else {
		return v.Fn(env, args)
	}
}

func NewFunction(name string, fn func(*Env, []Value) (Value, error), argNames Value) (*FuncValue, error) {
	if name == "" {
		name = "*lambda*"
	}

	if argNames.Type() != V_LIST {
		return nil, errors.New("Can't create a function, passed a non-list for argument names")
	}

	argNamesString := []string{}
	for _, arg := range argNames.(ListValue).Childs {
		// TODO handle "." as rest arg
		if arg.Type() == V_SYMBOL {
			argNamesString = append(argNamesString, arg.(SymbolValue).Value)
		} else {
			return nil, errors.New("Can't define function '" + name + "' because param " + arg.String() + " is not a symbol")
		}
	}

	return &FuncValue{
		Name:     name,
		Fn:       fn,
		ArgNames: argNamesString,
	}, nil
}

/* MACRO */
type MacroValue struct {
	Name     string
	ArgNames []string
	Fn       func(*Env, []Value) (Value, error)
}

func (MacroValue) Type() ValueType {
	return V_MACRO
}
func (v MacroValue) String() string {
	return "<macro:" + v.Name + ">"
}

func (v *MacroValue) Call(args []Value, parentEnv *Env) (Value, error) {
	if env, err := createEnvAndAssignArgs(v.Name, parentEnv, args, v.ArgNames); err != nil {
		return nil, err
	} else {
		return v.Fn(env, args)
	}
}

func NewMacro(name string, fn func(*Env, []Value) (Value, error), argNames Value) (*MacroValue, error) {
	if name == "" {
		name = "*macro*"
	}

	if argNames.Type() != V_LIST {
		return nil, errors.New("Can't create a macro, passed a non-list for argument names")
	}

	argNamesString := []string{}
	for _, arg := range argNames.(ListValue).Childs {
		// TODO handle "." as rest arg
		if arg.Type() == V_SYMBOL {
			argNamesString = append(argNamesString, arg.(SymbolValue).Value)
		} else {
			return nil, errors.New("Can't define macro '" + name + "' because param " + arg.String() + " is not a symbol")
		}
	}

	return &MacroValue{
		Name:     name,
		Fn:       fn,
		ArgNames: argNamesString,
	}, nil
}

/* utils */
func createEnvAndAssignArgs(fnName string, parentEnv *Env, args []Value, argNames []string) (*Env, error) {
	env := NewRootEnv()
	env.Parent = parentEnv
	for i, argName := range argNames {
		if argName == "." {
			if len(argNames) > i+1 {
				if len(args) >= i {
					env.Set(argNames[i+1], ListValue{args[i:]})
				} else {
					env.Set(argNames[i+1], NewEmptyListValue())
				}
				return env, nil
			} else {
				return nil, errors.New("Function '" + fnName + "' argument list can't end with '.'. Add a rest argument name")
			}
		}
		if i < len(args) {
			env.Set(argName, args[i])
		} else {
			return nil, errors.New("Function '" + fnName + "' argument named '" + argName + "' wasn't passed value. Expected argument count " + strconv.FormatInt(int64(len(argNames)), 10) + ". Got " + strconv.FormatInt(int64(len(args)), 10))
		}
	}
	return env, nil
}
