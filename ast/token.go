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
