package rpn

import (
	"errors"
	gtoken "go/token"
	"strings"
	"unicode"
)

var ErrNotSupportOperation = errors.New("not support opertation")

var operationPrecedence = map[gtoken.Token]int{
	gtoken.OR:  0, // |
	gtoken.AND: 0, // &

	gtoken.EQL: 1, // ==
	gtoken.LSS: 1, // <
	gtoken.GTR: 1, // >
	gtoken.NEQ: 1, // !=
	gtoken.LEQ: 1, // <=
	gtoken.GEQ: 1, // >=

	gtoken.ADD: 2, // +
	gtoken.SUB: 2, // -

	gtoken.MUL: 3, // *
	gtoken.QUO: 4, // /
	gtoken.XOR: 5, // ^
}

// left - false
// right - true
var operationAssociativity = map[gtoken.Token]bool{
	gtoken.XOR: true, // ^
}

type Parser struct {
	lexer *Lexer
}

func deleteSpace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func NewParser(src string) *Parser {
	return &Parser{
		lexer: NewLexer([]byte(deleteSpace(src))),
	}

}

func (p *Parser) Parse() ([]*Token, error) {
	stack := NewStack[*Token]()
	rpn := make([]*Token, 0, p.lexer.Size)

	for {
		t := p.lexer.Consume()
		if t.Kind == gtoken.EOF || t.Kind == gtoken.SEMICOLON {
			break
		}
		if t.Kind.IsLiteral() {
			rpn = append(rpn, &t)
			continue
		}
		if t.Kind.IsOperator() {
			switch t.Kind {
			case gtoken.LPAREN:
				stack.Push(&t)
			case gtoken.RPAREN:
				for stack.Count() > 0 && stack.Peek().Kind != gtoken.LPAREN {
					rpn = append(rpn, stack.Pop())
				}
				stack.Pop()
			default:
				if priority, ok := operationPrecedence[t.Kind]; ok {
					for stack.Count() > 0 && stack.Peek().Kind != gtoken.LPAREN && (operationPrecedence[stack.Peek().Kind] > priority || (operationPrecedence[stack.Peek().Kind] == priority && !operationAssociativity[t.Kind])) {
						rpn = append(rpn, stack.Pop())
					}
					stack.Push(&t)
				}
			}
		}
	}
	for stack.Count() > 0 {
		rpn = append(rpn, stack.Pop())
	}
	return rpn, nil
}
