package ry

import (
	"errors"
	"fmt"

	glisp "github.com/kiasaki/glisp/interpreter"
)

type SexpBuffer struct {
	Name     string
	Filename string
}

func (b SexpBuffer) SexpString() string {
	return fmt.Sprintf(`#<buffer %s>`, b.Name)
}

func NewSexpBuffer(name string, filename string) SexpBuffer {
	return SexpBuffer{
		Name:     name,
		Filename: filename,
	}
}

func MakeBufferFunction(env *glisp.Glisp, fnname string,
	args []glisp.Sexp) (glisp.Sexp, error) {
	var name string
	var filename string

	if len(args) >= 1 {
		switch expr := args[0].(type) {
		case glisp.SexpStr:
			name = string(expr)
		default:
			return glisp.SexpNull, errors.New("make-buffer passed a non-string parameter 1")
		}
	}

	if len(args) >= 2 {
		switch expr := args[1].(type) {
		case glisp.SexpStr:
			filename = string(expr)
		default:
			return glisp.SexpNull, errors.New("make-buffer passed a non-string parameter 2")
		}
	}

	return NewSexpBuffer(name, filename), nil
}
