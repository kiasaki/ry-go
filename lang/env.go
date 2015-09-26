package lang

type Env struct {
	Parent *Env
	Values map[string]Value
}
