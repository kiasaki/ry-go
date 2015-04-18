package ry

import (
	"errors"
	"fmt"

	syp "github.com/kiasaki/syp-lang/interpreter"
)

type SexpBuffer struct {
	Buffer *Buffer
}

func (b SexpBuffer) SexpString() string {
	return fmt.Sprintf(`#<buffer %s>`, b.Buffer.Name)
}

func MakeBufferFunction(env *syp.Lang, fnname string,
	args []syp.Sexp) (syp.Sexp, error) {
	var name string
	var filename string

	if len(args) >= 1 {
		switch expr := args[0].(type) {
		case syp.SexpStr:
			name = string(expr)
		default:
			return syp.SexpNull, errors.New("make-buffer passed a non-string parameter 1")
		}
	}

	if len(args) >= 2 {
		switch expr := args[1].(type) {
		case syp.SexpStr:
			filename = string(expr)
		default:
			return syp.SexpNull, errors.New("make-buffer passed a non-string parameter 2")
		}
	}

	return SexpBuffer{NewBuffer(name, filename)}, nil
}
