package ry

import (
	syp "github.com/kiasaki/syp-lang/interpreter"
)

// RegisterToEnv adds all ry's go primitives realted to the Editor, Buffers,
// Points, Marks, Windows and rendering to a Lang environment for use by the
// runtime
func RegisterToEnv(env *syp.Lang) {
	// Editor
	env.AddFunction("start-editor", StartEditorFunction)
	env.AddFunction("stop-editor", StopEditorFunction)
	env.AddFunction("editor-keypresses-chan", EditorKeypressesChan)
	env.AddFunction("clear-editor", ClearEditorFunction)
	env.AddFunction("set-cell", SetCellFunction)

	// Buffer
	env.AddFunction("make-buffer", MakeBufferFunction)
}
