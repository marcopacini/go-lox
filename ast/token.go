//  MIT License
//
//  Copyright (c) 2019 Marco Pacini
//
//  Permission is hereby granted, free of charge, to any person obtaining a copy
//  of this software and associated documentation files (the "Software"), to deal
//  in the Software without restriction, including without limitation the rights
//  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//  copies of the Software, and to permit persons to whom the Software is
//  furnished to do so, subject to the following conditions:
//
//  The above copyright notice and this permission notice shall be included in all
//  copies or substantial portions of the Software.
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//  SOFTWARE.

package ast

import "fmt"

type TokenType int

const (
	And TokenType = iota
	Class
	Comma
	Dot
	Else
	Eof
	Equal
	EqualEqual
	False
	For
	Fun
	Greater
	GreaterEqual
	Identifier
	If
	LeftParenthesis
	LeftSquare
	Less
	LessEqual
	Minus
	Nil
	Not
	NotEqual
	Number
	Or
	Plus
	Print
	Return
	RightParenthesis
	RightSquare
	Semicolon
	Slash
	Star
	String
	Super
	This
	True
	Var
	While
)

var keywords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"fun":    Fun,
	"for":    For,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

func (t TokenType) String() string {
	switch t {
	case LeftParenthesis:
		return "LEFT_PARENTHESIS"
	case RightParenthesis:
		return "RIGHT_PARENTHESIS"
	case LeftSquare:
		return "LEFT_SQUARE"
	case RightSquare:
		return "RIGHT_SQUARE"
	case Comma:
		return "COMMA"
	case Dot:
		return "DOT"
	case Minus:
		return "MINUS"
	case Plus:
		return "PLUS"
	case Semicolon:
		return "SEMICOLON"
	case Slash:
		return "SLASH"
	case Star:
		return "STAR"
	case Not:
		return "NOT"
	case Equal:
		return "EQUAL"
	case EqualEqual:
		return "EQUAL_EQUAL"
	case NotEqual:
		return "NOT_EQUAL"
	case Greater:
		return "GREATER"
	case GreaterEqual:
		return "GREATER_EQUAL"
	case Less:
		return "LESS"
	case LessEqual:
		return "LESS_EQUAL"
	case String:
		return "STRING"
	case Number:
		return "NUMBER"
	case Identifier:
		return "IDENTIFIER"
	case Eof:
		return "EOF"
	case And:
		return "AND"
	case Class:
		return "CLASS"
	case Var:
		return "VAR"
	case If:
		return "IF"
	case Else:
		return "ELSE"
	case For:
		return "FOR"
	case While:
		return "WHILE"
	case Print:
		return "PRINT"
	case True:
		return "TRUE"
	case False:
		return "FALSE"
	case Nil:
		return "NIL"
	}

	return "UNKNOWN"
}

type Token struct {
	TokenType
	Lexeme  string
	Literal string
	Line    int
}

func (t Token) String() string {
	return fmt.Sprintf("%v %v %v %d", t.TokenType, t.Lexeme, t.Literal, t.Line)
}
