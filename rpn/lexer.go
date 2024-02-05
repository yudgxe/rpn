package rpn

import (
	"go/scanner"
	gtoken "go/token"
)

type Token struct {
	Kind  gtoken.Token
	Value string
}

type Lexer struct {
	scanner scanner.Scanner

	peeked bool
	tok    Token

	Size int
}

func NewLexer(src []byte) *Lexer {
	fset := gtoken.NewFileSet()
	file := fset.AddFile("<expr.go>", fset.Base(), len(src))

	var s scanner.Scanner
	s.Init(file, src, nil, scanner.ScanComments)

	return &Lexer{
		scanner: s,
		Size:    len(src),
	}
}

func (l *Lexer) Consume() Token {
	if l.peeked {
		l.peeked = false
		return l.tok
	}
	_, kind, value := l.scanner.Scan()
	return Token{
		Kind:  kind,
		Value: value,
	}
}

func (l *Lexer) Peek() Token {
	l.tok = l.Consume()
	l.peeked = true
	return l.tok
}
