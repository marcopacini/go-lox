package ast

import (
	"fmt"
	"strconv"
)

type Parser struct {
	Tokens  []Token
	current int
}

func (p Parser) peek() Token {
	return p.Tokens[p.current]
}

func (p *Parser) previous() (Token, bool) {
	if p.current-1 < 0 {
		return Token{}, false
	}

	return p.Tokens[p.current-1], true
}

func (p Parser) isEnd() bool {
	return p.peek().TokenType == Eof
}

func (p *Parser) advance() (Token, bool) {
	if p.isEnd() {
		return Token{}, false
	} else {
		p.current++
	}

	return p.previous()
}

func (p *Parser) match(tt ...TokenType) bool {
	for _, t := range tt {
		if !p.isEnd() && p.peek().TokenType == t {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(t TokenType) (Token, error) {
	token := p.peek()

	if token.TokenType != t {
		return Token{}, fmt.Errorf("error at line %d: expected '%v'", token.Line, t.String())
	}

	p.advance()

	return token, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(NotEqual, EqualEqual) {
		if operator, ok := p.previous(); ok {
			right, err := p.comparison()
			if err != nil {
				return nil, err
			}

			expr = Binary{expr, operator, right}
		}
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.addition()
	if err != nil {
		return nil, err
	}

	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		if operator, ok := p.previous(); ok {
			right, err := p.addition()
			if err != nil {
				return nil, err
			}

			expr = Binary{expr, operator, right}
		}
	}

	return expr, nil
}

func (p *Parser) addition() (Expr, error) {
	expr, err := p.multiplication()
	if err != nil {
		return nil, err
	}

	for p.match(Minus, Plus) {
		if operator, ok := p.previous(); ok {
			right, err := p.multiplication()
			if err != nil {
				return nil, err
			}

			expr = Binary{expr, operator, right}
		}
	}

	return expr, nil
}

func (p *Parser) multiplication() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(Slash, Star) {
		if operator, ok := p.previous(); ok {
			right, err := p.unary()
			if err != nil {
				return nil, err
			}

			expr = Binary{expr, operator, right}
		}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(Not, Minus) {
		if operator, ok := p.previous(); ok {
			right, err := p.unary()
			if err != nil {
				return nil, err
			}

			return Unary{operator, right}, nil
		}
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(True) {
		return Literal{true}, nil
	}

	if p.match(False) {
		return Literal{false}, nil
	}

	if p.match(Nil) {
		return Literal{nil}, nil
	}

	if p.match(Number) {
		if token, ok := p.previous(); ok {
			value, err := strconv.ParseFloat(token.Literal, 64)
			if err != nil {
				return nil, err
			}

			return Literal{value}, nil
		}
	}

	if p.match(String) {
		if token, ok := p.previous(); ok {
			return Literal{token.Literal}, nil
		}
	}

	if p.match(LeftParenthesis) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if _, err := p.consume(RightParenthesis); err != nil {
			return nil, err
		}

		return Grouping{expr}, nil
	}

	return nil, fmt.Errorf("error at line %d: unknown token '%s'", p.peek().Line, p.peek().Literal)
}

func (p *Parser) Parse() (Expr, error) {
	return p.expression()
}
