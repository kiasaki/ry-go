package lang

type Env struct {
	Parent *Env
	Values map[string]Value
}

func (e *Env) Get(key string) Value {
	val, ok := e.Values[key]
	if ok {
		return val
	} else if !ok && e.Parent != nil {
		return e.Parent.Get(key)
	}
	return nil
}

func (e *Env) Set(key string, val Value) {
	e.Values[key] = val
}

func (e *Env) SetOnRoot(key string, val Value) {
	if e.Parent != nil {
		e.Parent.SetOnRoot(key, val)
	} else {
		e.Set(key, val)
	}
}

func NewRootEnv() *Env {
	return &Env{nil, map[string]Value{}}
}

func NewBuiltinFilledEnv() *Env {
	env := NewRootEnv()

	env.Set("+", builtinAdd)
	env.Set("-", builtinSubtract)
	env.Set("*", builtinMultiply)
	env.Set("/", builtinDivide)
	env.Set("<", builtinSmallerThan)
	env.Set(">", builtinBiggerThan)
	env.Set("<=", builtinSmallerThanEq)
	env.Set(">=", builtinBiggerThanEq)
	env.Set("floor", builtinFloor)

	env.Set("type", builtinType)
	env.Set("eq?", builtinEq)
	env.Set("car", builtinCar)
	env.Set("cdr", builtinCdr)
	env.Set("cons", builtinCons)
	env.Set("append", builtinAppend)
	// reverse
	env.Set("length", builtinLength)
	env.Set("list-ref", builtinListRef)
	env.Set("list-set!", builtinListSet)
	env.Set("string-regexp-split", builtinStringRegexpSplit)
	env.Set("string-length", builtinStringLength)

	env.Set("list", builtinList)
	env.Set("string", builtinString)

	env.Set("integer->char", builtinIntegerToChar)
	env.Set("char->integer", builtinCharToInteger)
	env.Set("integer->string", builtinIntegerToString)
	env.Set("string->integer", builtinStringToInteger)

	env.Set("quote", builtinQuote)
	env.Set("unquote", builtinUnquote)
	env.Set("eval", builtinEval)
	env.Set("define", builtinDefine)
	env.Set("set!", builtinSet)
	env.Set("let", builtinLet)
	env.Set("let*", builtinLetStar)
	env.Set("letrec", builtinLetRec)
	env.Set("letrec*", builtinLetRecStar)
	env.Set("lambda", builtinLambda)
	env.Set("defmacro", builtinDefmacro)

	env.Set("if", builtinIf)
	env.Set("cond", builtinCond)
	env.Set("begin", builtinBegin)

	env.Set("include", builtinInclude)
	env.Set("error", builtinError)
	env.Set("read", builtinRead)
	env.Set("write", builtinWrite)

	return env
}
