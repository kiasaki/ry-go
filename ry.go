package ry

import (
	syp "github.com/kiasaki/syp-lang/interpreter"
)

// RegisterToEnv adds all ry's go primitives realted to the Editor, Buffers,
// Points, Marks, Windows and rendering to a Lang environment for use by the
// runtime
func RegisterToEnv(env *syp.Lang) {
	env.AddFunction("make-buffer", MakeBufferFunction)
}
