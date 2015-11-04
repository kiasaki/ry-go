package lang

import (
	"errors"
	"strconv"
	"strings"

	"github.com/kiasaki/sexpr"
)

func ParseFile(file string) ([]Value, error) {
	ast := &sexpr.AST{}
	err := sexpr.ParseFile(ast, file, buildSyntaxParser())
	if err != nil {
		return nil, err
	}
	return Read(ast.Root.Children)
}

func Parse(code []byte) ([]Value, error) {
	ast := &sexpr.AST{}
	err := sexpr.Parse(ast, code, buildSyntaxParser())
	if err != nil {
		return nil, err
	}
	return Read(ast.Root.Children)
}

func Read(nodes []*sexpr.Node) ([]Value, error) {
	return readNodes(nodes)
}

func readNodes(nodes []*sexpr.Node) ([]Value, error) {
	var values = []Value{}

	for _, node := range nodes {
		value, err := readASTNode(node)
		if err != nil {
			return nil, err
		}
		// In case we hit non-value tokens, like comments
		if value != nil {
			values = append(values, value)
		}
	}

	return values, nil
}

func readASTNode(node *sexpr.Node) (Value, error) {
	nodeValue := string(node.Data)
	switch node.Type {
	case sexpr.TokListOpen:
		if nodes, err := readNodes(node.Children); err != nil {
			return nil, err
		} else {
			if nodeValue[0] == '\'' {
				return ListValue{append([]Value{SymbolValue{"quote"}}, ListValue{nodes})}, nil
			} else {
				return ListValue{nodes}, nil
			}
		}
	case sexpr.TokIdent:
		if nodeValue[0] == '\'' {
			return ListValue{[]Value{SymbolValue{"quote"}, SymbolValue{nodeValue[1:]}}}, nil
		} else {
			return SymbolValue{nodeValue}, nil
		}
	case sexpr.TokString:
		return StringValue{nodeValue}, nil
	case sexpr.TokRawString:
		return StringValue{strconv.Quote(nodeValue)}, nil
	case sexpr.TokChar:
		if string(nodeValue) == "space" {
			nodeValue = " "
		}
		if len([]rune(nodeValue)) != 1 {
			return nil, errors.New("Tried reading a char literal of 0 or more than 1 characters")
		}
		return CharValue{[]rune(nodeValue)[0]}, nil
	case sexpr.TokNumber:
		if strings.Contains(nodeValue, ".") {
			if f, err := strconv.ParseFloat(nodeValue, 64); err == nil {
				return FloatValue{f}, nil
			}
		} else {
			if i, err := strconv.ParseInt(nodeValue, 0, 64); err == nil {
				return IntegerValue{i}, nil
			}
		}
	case sexpr.TokBoolean:
		if nodeValue == "#t" {
			return BoolValueTrue, nil
		} else {
			return BoolValueFalse, nil
		}

	case sexpr.TokComment:
		return nil, nil
	}

	return nil, errors.New("Couln't parse invalid token '" + nodeValue + "'")
}

func buildSyntaxParser() *sexpr.Syntax {
	s := new(sexpr.Syntax)

	s.StringLit = []string{"\"", "\""}
	s.RawStringLit = []string{"`", "`"}
	s.CharLit = []string{"#\\", " "}
	s.Delimiters = [][2]string{{"'(", ")"}, {"(", ")"}}
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
