package lang

type Env struct {
	Parent *Env
	Values map[string]Value
}

func NewRootEnv() *Env {
	return &Env{nil, map[string]Value{}}
}

func (e *Env) Get(key string) {
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
