package ast

import "testing"

func TestScanner_Scan(t *testing.T) {
	table := []struct {
		in  string
		out []TokenType
	}{
		{"(){}", []TokenType{LeftParenthesis, RightParenthesis, LeftSquare, RightSquare, Eof}},
		{"+ - * / , ; ! > <", []TokenType{Plus, Minus, Star, Slash, Comma, Semicolon, Not, Greater, Less, Eof}},
		{"== != >= <=", []TokenType{EqualEqual, NotEqual, GreaterEqual, LessEqual, Eof}},
		{"// This text have to be ignored", []TokenType{Eof}},
		{"\"This is a string!\"", []TokenType{String, Eof}},
		{"1 12 12.3", []TokenType{Number, Number, Number, Eof}},
		{"and or true false", []TokenType{And, Or, True, False, Eof}},
		{"if else for while", []TokenType{If, Else, For, While, Eof}},
		{"fun return", []TokenType{Fun, Return, Eof}},
		{"class var nil", []TokenType{Class, Var, Nil, Eof}},
		{"print x", []TokenType{Print, Identifier, Eof}},
	}

	for _, test := range table {
		t.Run(test.in, func(t *testing.T) {
			scanner := Scanner{test.in}
			tokens, _ := scanner.Scan()

			if len(tokens) != len(test.out) {
				t.Errorf("want %d tokens, got %d: %v != %v", len(tokens), len(test.out), tokens, test.out)
			} else {
				for i := range tokens {
					if tokens[i].TokenType != test.out[i] {
						t.Errorf("want %v, got %v: %v != %v", tokens[i].TokenType, test.out[i], tokens, test.out)
					}
				}
			}
		})
	}
}
