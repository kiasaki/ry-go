package lang

import (
	"strconv"
	"strings"

	"github.com/kiasaki/sexpr"
)

func ParseFile(file string) ([]Value, error) {
	ast := &sexpr.AST{}
	err := sexpr.ParseFile(ast, file, buildSyntaxParser())
	return Read(ast.Root.Children), err
}

func Parse(code []byte) ([]Value, error) {
	ast := &sexpr.AST{}
	err := sexpr.Parse(ast, code, buildSyntaxParser())
	return Read(ast.Root.Children), err
}

func Read(nodes []*sexpr.Node) []Value {
	return readNodes(nodes)
}

func readNodes(nodes []*sexpr.Node) []Value {
	var values = []Value{}

	for _, node := range nodes {
		value := readASTNode(node)
		if value != nil {
			values = append(values, value)
		}
	}

	return values
}

func readASTNode(node *sexpr.Node) Value {
	nodeValue := string(node.Data)
	switch node.Type {
	case sexpr.TokListOpen:
		if nodeValue[0] == '\'' {
			return ListValue{append([]Value{SymbolValue{"quote"}}, ListValue{readNodes(node.Children)})}
		} else {
			return ListValue{readNodes(node.Children)}
		}
	case sexpr.TokIdent:
		return SymbolValue{nodeValue}
	case sexpr.TokString:
		return StringValue{nodeValue}
	case sexpr.TokRawString:
		return StringValue{strconv.Quote(nodeValue)}
	case sexpr.TokChar:
		return CharValue{[]rune(nodeValue)[0]}
	case sexpr.TokNumber:
		if strings.Contains(nodeValue, ".") {
			if f, err := strconv.ParseFloat(nodeValue, 64); err == nil {
				return FloatValue{f}
			}
		} else {
			if i, err := strconv.ParseInt(nodeValue, 0, 64); err == nil {
				return IntegerValue{i}
			}
		}
	case sexpr.TokBoolean:
		if nodeValue == "#t" {
			return BoolValue{true}
		} else {
			return BoolValue{false}
		}
	}
	return nil
}

func buildSyntaxParser() *sexpr.Syntax {
	s := new(sexpr.Syntax)

	s.StringLit = []string{"\"", "\""}
	s.RawStringLit = []string{"`", "`"}
	s.CharLit = []string{"'", "'"}
	s.Delimiters = [][2]string{{"(", ")"}}
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
