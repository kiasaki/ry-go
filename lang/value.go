package lang

type ValueType int

const (
	V_SYMBOL  ValueType = iota
	V_STRING            = iota
	V_INTEGER           = iota
	V_FLOAT             = iota
	V_CHAR              = iota
	V_BOOL              = iota
	V_LIST              = iota
)

type Value interface {
	Type() ValueType
	String() string
}

/* SYMBOL */
type SymbolValue struct {
	value string
}

func (SymbolValue) Type() {
	return V_SYMBOL
}
func (v *SymbolValue) String() {
	return v.value
}

/* STRING */
type StringValue struct {
	value string
}

func (StringValue) Type() {
	return V_STRING
}
func (v *StringValue) String() {
	return "\"" + v.value + "\""
}

/* INTEGER */
type IntegerValue struct {
	value int64
}

func (IntegerValue) Type() {
	return V_INTEGER
}
func (v *IntegerValue) String() {
	return string(v.value)
}

/* FLOAT */
type FloatValue struct {
	value float64
}

func (FloatValue) Type() {
	return V_FLOAT
}
func (v *FloatValue) String() {
	return string(v.value)
}

/* CHAR */
type CharValue struct {
	value rune
}

func (CharValue) Type() {
	return V_CHAR
}
func (v *CharValue) String() {
	return "'" + string(v.value) + "'"
}

/* BOOL */
type BoolValue struct {
	value bool
}

func (BoolValue) Type() {
	return V_BOOL
}
func (v *CharValue) String() {
	if v.value {
		return "#t"
	}
	return "#f"
}

/* LIST */
type ListValue struct {
	childs []Value
}

func (ListValue) Type() {
	return V_BOOL
}
func (v *ListValue) String() {
	str := "("
	for i, c := range v.childs {
		if i != 0 {
			str = str + " "
		}
		str = str + c.String()
	}
	return str + ")"
}
