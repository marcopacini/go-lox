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

func (p *Parser) declaration() (Stmt, error) {
	if p.match(Var) {
		return p.variable()
	}

	return p.statement()
}

func (p *Parser) variable() (Stmt, error) {
	token, err := p.consume(Identifier)
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.match(Equal) {
		if initializer, err = p.expression(); err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(Semicolon); err != nil {
		return nil, err
	}

	return Declaration{token, initializer}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(If) {
		if _, err := p.consume(LeftParenthesis); err != nil {
			return nil, err
		}

		condition, err := p.expression()
		if err != nil {
			return nil, err
		}

		if _, err := p.consume(RightParenthesis); err != nil {
			return nil, err
		}

		thenBranch, err := p.statement()
		if err != nil {
			return nil, err
		}

		var elseBranch Stmt
		if p.match(Else) {
			if elseBranch, err = p.statement(); err != nil {
				return nil, err
			}
		}

		return IfStmt{condition, thenBranch, elseBranch}, nil
	}

	if p.match(For) {
		if _, err := p.consume(LeftParenthesis); err != nil {
			return nil, err
		}

		var init Stmt
		var err error

		if p.match(Semicolon) {
			init = nil
		} else if p.match(Var) {
			init, err = p.variable()
			if err != nil {
				return nil, err
			}
		} else {
			expr, err := p.expression()
			if err != nil {
				return nil, err
			}

			init = ExprStmt{expr}
		}

		var condition Expr

		if p.peek().TokenType != Semicolon {
			condition, err = p.expression()
			if err != nil {
				return nil, err
			}
		}

		if _, err := p.consume(Semicolon); err != nil {
			return nil, err
		}

		var increment Expr

		if p.peek().TokenType != RightParenthesis {
			increment, err = p.expression()
			if err != nil {
				return nil, err
			}
		}

		if _, err := p.consume(RightParenthesis); err != nil {
			return nil, err
		}

		body, err := p.statement()
		if err != nil {
			return nil, err
		}

		return ForStmt{init, condition, increment, body}, nil
	}

	if p.match(Fun) {
		name, err := p.consume(Identifier)
		if err != nil {
			return nil, err
		}

		if _, err := p.consume(LeftParenthesis); err != nil {
			return nil, err
		}

		var arguments []Token
		if p.peek().TokenType != RightParenthesis {
			for true {
				token, err := p.consume(Identifier)
				if err != nil {
					return nil, err
				}

				arguments = append(arguments, token)

				if !p.match(Comma) {
					break
				}
			}
		}

		if _, err := p.consume(RightParenthesis); err != nil {
			return nil, err
		}

		if _, err := p.consume(LeftSquare); err != nil {
			return nil, err
		}

		body, err := p.block()
		if err != nil {
			return nil, err
		}

		return Function{name, nil, arguments, body}, nil
	}

	if p.match(Print) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if _, err := p.consume(Semicolon); err != nil {
			return nil, err
		}

		return PrintStmt{expr}, nil
	}

	if p.match(Return) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if _, err := p.consume(Semicolon); err != nil {
			return nil, err
		}

		return ReturnStmt{expr}, nil
	}

	if p.match(While) {
		if _, err := p.consume(LeftParenthesis); err != nil {
			return nil, err
		}

		condition, err := p.expression()
		if err != nil {
			return nil, err
		}

		if _, err := p.consume(RightParenthesis); err != nil {
			return nil, err
		}

		body, err := p.statement()
		if err != nil {
			return nil, err
		}

		return WhileStmt{condition, body}, nil
	}

	if p.match(LeftSquare) {
		b, err := p.block()
		if err != nil {
			return nil, err
		}

		return Block{b}, nil
	}

	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(Semicolon); err != nil {
		return nil, err
	}

	return ExprStmt{expr}, nil
}

func (p *Parser) block() ([]Stmt, error) {
	var stmts []Stmt

	for !(p.peek().TokenType == RightSquare) && !p.isEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}

		stmts = append(stmts, stmt)
	}

	if _, err := p.consume(RightSquare); err != nil {
		return nil, err
	}

	return stmts, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(Equal) {
		if t, ok := p.previous(); ok {
			value, err := p.assignment()
			if err != nil {
				return nil, err
			}

			if v, ok := expr.(Variable); ok {
				return Assign{v, t, value}, nil
			}

			return nil, fmt.Errorf("error at line %d: invalid assignment target", t.Line)
		}
	}

	return expr, nil
}

func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(Or) {
		if operator, ok := p.previous(); ok {
			right, err := p.and()
			if err != nil {
				return nil, err
			}

			expr = Logical{expr, operator, right}
		}
	}

	return expr, nil
}

func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(And) {
		if operator, ok := p.previous(); ok {
			right, err := p.equality()
			if err != nil {
				return nil, err
			}

			expr = Logical{expr, operator, right}
		}
	}

	return expr, nil
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

	return p.call()
}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for true {
		if p.match(LeftParenthesis) {
			var arguments []Expr
			if p.peek().TokenType != RightParenthesis {
				for true {
					argument, err := p.expression()
					if err != nil {
						return nil, err
					}

					arguments = append(arguments, argument)

					if !p.match(Comma) {
						break
					}
				}
			}

			if _, err := p.consume(RightParenthesis); err != nil {
				return nil, err
			}

			return Call{expr, arguments}, nil
		} else {
			break
		}
	}

	return expr, nil
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

	if p.match(Identifier) {
		if token, ok := p.previous(); ok {
			return Variable{token}, nil
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

func (p *Parser) Parse() ([]Stmt, error) {
	var stmts []Stmt

	for !p.isEnd() {
		if stmt, err := p.declaration(); err != nil {
			return nil, err
		} else {
			stmts = append(stmts, stmt)
		}
	}

	return stmts, nil
}
