package lang

import (
	"github.com/kiasaki/sexpr"
)

func ParseFile(file string) (string, error) {
	env := &Env{nil, map[string]Value{}}
	ast := &sexpr.AST{}
	err := sexpr.ParseFile(ast, file, buildSyntaxParser())
	return Read(env, ast).String(), err
}

func Parse(code []byte) (string, error) {
	env := &Env{nil, map[string]Value{}}
	ast := &sexpr.AST{}
	err := sexpr.Parse(ast, code, buildSyntaxParser())
	return Read(env, ast).String(), err
}

func Read(env *Env, ast *sexpr.AST) Value {
	var values = []Value{}

	for _, node := range ast.Root.Children {
		values = append(values, readASTNode(node))
	}

	// TODO return all
	return values[0]
}

func readASTNode(node *sexpr.Node) Value {
	nodeValue := string(node.Data)
	switch node.Type {
	case sexpr.TokListOpen:
	case sexpr.TokIdent:
	case sexpr.TokString:
	case sexpr.TokRawString:
	case sexpr.TokChar:
	case sexpr.TokNumber:
	case sexpr.TokBool:
		if nodeValue == "#t" {
			return BoolValue{true}
		} else {
			return BoolValue{false}
		}
	default:
	}
}

func buildSyntaxParser() *sexpr.Syntax {
	s := new(sexpr.Syntax)

	s.StringLit = []string{"\"", "\""}
	s.RawStringLit = []string{"`", "`"}
	s.CharLit = []string{"'", "'"}
	s.Delimiters = [][2]string{{"(", ")"}, {"'(", ")"}}
	s.NumberFunc = sexpr.LexNumber
	s.BooleanFunc = func(l *sexpr.Lexer) int {
		if ret := l.AcceptLiteral("#t"); ret != 0 {
			return ret
		}
		return l.AcceptLiteral("#f")
	}
	s.SingleLineComment = ";"
	s.MultiLineComment = []string{"#|", "|#"}

	return s
}
