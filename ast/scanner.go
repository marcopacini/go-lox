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

type Scanner struct {
	Text string
}

func (s *Scanner) Scan() ([]Token, error) {
	runes := []rune(s.Text)

	start := 0
	current := 0
	line := 1

	tokens := make([]Token, 0)

	isEnd := func() bool {
		return current >= len(runes)
	}

	peek := func() rune {
		if isEnd() {
			return '\x00'
		}

		return runes[current]
	}

	peekNext := func() rune {
		if current+1 >= len(runes) {
			return '\x00'
		}

		return runes[current+1]
	}

	advance := func() rune {
		current++
		return runes[current-1]
	}

	isNext := func(r rune) bool {
		if isEnd() || runes[current] != r {
			return false
		}

		current++

		return true
	}

	isDigit := func(r rune) bool {
		if r >= '0' && r <= '9' {
			return true
		}

		return false
	}

	isLetter := func(r rune) bool {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return true
		}

		return false
	}

	addToken := func(tokenType TokenType) {
		tokens = append(tokens, Token{tokenType, string(runes[start:current]), "", line})
	}

	scanToken := func() error {
		r := advance()

		switch r {
		case ' ':
		case '\t':
		case '\r':
			{
				break // ignore withe space
			}

		case '\n':
			{
				line++
				break
			}

		// Single-character lexeme: '(', ')', '[', ']', '.', '-', '+', ';'
		case '(':
			{
				addToken(LeftParenthesis)
				break
			}

		case ')':
			{
				addToken(RightParenthesis)
				break
			}

		case '{':
			{
				addToken(LeftSquare)
				break
			}

		case '}':
			{
				addToken(RightSquare)
				break
			}

		case '.':
			{
				addToken(Dot)
				break
			}

		case '-':
			{
				addToken(Minus)
				break
			}

		case '+':
			{
				addToken(Plus)
				break
			}

		case '*':
			{
				addToken(Star)
				break
			}

		case ',':
			{
				addToken(Comma)
				break
			}

		case ';':
			{
				addToken(Semicolon)
				break
			}

		// Multi-character lexeme (potentially): '/', '!', '=', '<', '>', '!=', '==', '<=', '>=', '//'
		case '!':
			{
				if isNext('=') {
					addToken(NotEqual)
				} else {
					addToken(Not)
				}

				break
			}

		case '=':
			{
				if isNext('=') {
					addToken(EqualEqual)
				} else {
					addToken(Equal)
				}

				break
			}

		case '>':
			{
				if isNext('=') {
					addToken(GreaterEqual)
				} else {
					addToken(Greater)
				}

				break
			}

		case '<':
			{
				if isNext('=') {
					addToken(LessEqual)
				} else {
					addToken(Less)
				}
			}

		case '/':
			{
				if isNext('/') {
					for peek() != '\n' && !isEnd() {
						advance()
					}
				} else {
					addToken(Slash)
				}
			}

		case '"':
			{
				for peek() != '"' && !isEnd() {
					if peek() == '\n' {
						line++
					}

					advance()
				}

				// unterminated string
				if isEnd() {
					return fmt.Errorf("error at line %d: unterminated string", line)
				}

				advance()

				lexeme := string(runes[start:current])
				literal := string(runes[start+1 : current-1]) // remove double quotes

				tokens = append(tokens, Token{String, lexeme, literal, line})
			}

		default:
			{
				if isDigit(r) {
					for isDigit(peek()) {
						advance()
					}

					if peek() == '.' && isDigit(peekNext()) {
						advance()

						for isDigit(peek()) {
							advance()
						}
					}

					number := string(runes[start:current])
					tokens = append(tokens, Token{Number, number, number, line})
				} else if isLetter(r) {
					for isLetter(peek()) || isDigit(peek()) {
						advance()
					}

					if t, ok := keywords[string(runes[start:current])]; ok {
						tokens = append(tokens, Token{t, string(runes[start:current]), "", line})
					} else {
						tokens = append(tokens, Token{Identifier, string(runes[start:current]), "", line})
					}
				} else {
					return fmt.Errorf("unknown character '%v' at line %d", string(r), line)
				}
			}
		}

		return nil
	}

	for !isEnd() {
		start = current
		if err := scanToken(); err != nil {
			return nil, err
		}
	}

	// cannot use addToken because lexeme will get the last character
	tokens = append(tokens, Token{Eof, "", "", line})

	return tokens, nil
}
