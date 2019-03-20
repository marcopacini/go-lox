package ast

import (
	"fmt"
)

type Interpreter struct {
	Literal
	*Environment
}

func (i *Interpreter) Run(stmts []Stmt) error {
	i.Environment = NewEnvironment(nil)

	for _, stmt := range stmts {
		if err := stmt.Accept(i); err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) Evaluate(expr Expr) (Literal, error) {
	err := expr.Accept(i)
	return i.Literal, err
}

func (i *Interpreter) visitAssign(expr Expr) error {
	if a, ok := expr.(Assign); ok {
		l, err := i.Evaluate(a.Expr)
		if err != nil {
			return err
		}

		if err := i.Environment.Assign(a.Variable, l); err != nil {
			return err
		}

		if l, ok := a.Expr.(Literal); ok {
			i.Literal = l
		}
	}

	return nil
}

func (i *Interpreter) visitBinary(expr Expr) error {
	if b, ok := expr.(Binary); ok {
		left, err := i.Evaluate(b.Left)
		if err != nil {
			return err
		}

		right, err := i.Evaluate(b.Right)
		if err != nil {
			return err
		}

		invalidOperand := func(left interface{}, right interface{}) error {
			return fmt.Errorf("error at line %d: invalid operands for binary %s: %T, %T", b.Operator.Line, b.Operator.Lexeme, left, right)
		}

		switch b.Operator.TokenType {
		case Plus:
			{
				if l, ok := left.Value.(float64); ok {
					if r, ok := right.Value.(float64); ok {
						// Sum of numbers
						i.Literal = Literal{l + r}
					}
				} else if l, ok := left.Value.(string); ok {
					if r, ok := right.Value.(string); ok {
						// String concatenation
						i.Literal = Literal{l + r}
					}
				} else {
					// Invalids operands
					return invalidOperand(left.Value, right.Value)
				}
			}
		case Minus:
			{
				if l, ok := left.Value.(float64); ok {
					if r, ok := right.Value.(float64); ok {
						i.Literal = Literal{l - r}
					}
				} else {
					return invalidOperand(left.Value, right.Value)
				}
			}
		case Star:
			{
				if l, ok := left.Value.(float64); ok {
					if r, ok := right.Value.(float64); ok {
						i.Literal = Literal{l * r}
					}
				} else {
					return invalidOperand(left.Value, right.Value)
				}
			}
		case Slash:
			{
				if l, ok := left.Value.(float64); ok {
					if r, ok := right.Value.(float64); ok {
						i.Literal = Literal{l / r}
					}
				} else {
					return invalidOperand(left.Value, right.Value)
				}
			}
		case EqualEqual:
			{
				i.Literal = Literal{left.Value == right.Value}
			}
		case NotEqual:
			{
				i.Literal = Literal{left.Value != right.Value}
			}
		case Greater:
			{
				if l, ok := left.Value.(float64); ok {
					if r, ok := right.Value.(float64); ok {
						i.Literal = Literal{l > r}
					} else {
						return invalidOperand(left.Value, right.Value)
					}
				}
			}
		case GreaterEqual:
			{
				if l, ok := left.Value.(float64); ok {
					if r, ok := right.Value.(float64); ok {
						i.Literal = Literal{l >= r}
					} else {
						return invalidOperand(left.Value, right.Value)
					}
				}
			}
		case Less:
			{
				if l, ok := left.Value.(float64); ok {
					if r, ok := right.Value.(float64); ok {
						i.Literal = Literal{l < r}
					} else {
						return invalidOperand(left.Value, right.Value)
					}
				}
			}
		case LessEqual:
			{
				if l, ok := left.Value.(float64); ok {
					if r, ok := right.Value.(float64); ok {
						i.Literal = Literal{l <= r}
					} else {
						return invalidOperand(left.Value, right.Value)
					}
				}
			}
		}
	}

	return nil
}

func (i *Interpreter) visitBlock(stmt Stmt) error {
	if b, ok := stmt.(Block); ok {
		i.Environment = NewEnvironment(i.Environment)

		for _, stmt := range b.Stmts {
			if err := stmt.Accept(i); err != nil {
				return err
			}
		}

		i.Environment = i.Environment.Parent
	}

	return nil
}

func (i *Interpreter) visitDeclaration(stmt Stmt) error {
	if d, ok := stmt.(Declaration); ok {
		i.Literal = Literal{nil}

		if d.Expr != nil {
			if _, err := i.Evaluate(d.Expr); err != nil {
				return err
			}
		}

		if err := i.Environment.Declare(Variable{d.Token}, i.Literal); err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) visitExprStmt(stmt Stmt) error {
	if s, ok := stmt.(ExprStmt); ok {
		return s.Expr.Accept(i)
	}

	return nil
}

func (i *Interpreter) visitForStmt(stmt Stmt) error {
	if f, ok := stmt.(ForStmt); ok {
		if f.Init != nil {
			if err := f.Init.Accept(i); err != nil {
				return err
			}
		}

		for true {
			l, err := i.Evaluate(f.Condition)
			if err != nil {
				return err
			}

			if !l.Bool() {
				return nil
			}

			if err := f.Body.Accept(i); err != nil {
				return err
			}

			if err := f.Increment.Accept(i); err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *Interpreter) visitIfStmt(stmt Stmt) error {
	if s, ok := stmt.(IfStmt); ok {
		l, err := i.Evaluate(s.Condition)
		if err != nil {
			return err
		}

		if l.Bool() {
			if err := s.Then.Accept(i); err != nil {
				return err
			}
		} else {
			if s.Else != nil {
				if err := s.Else.Accept(i); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (i *Interpreter) visitPrintStmt(stmt Stmt) error {
	if p, ok := stmt.(PrintStmt); ok {
		expr, err := i.Evaluate(p.Expr)
		if err != nil {
			return err
		}

		fmt.Println(expr)
	}

	return nil
}

func (i *Interpreter) visitGrouping(expr Expr) error {
	if g, ok := expr.(Grouping); ok {
		if _, err := i.Evaluate(g.Expr); err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) visitLiteral(expr Expr) error {
	if l, ok := expr.(Literal); ok {
		i.Literal = l
	}

	return nil
}

func (i *Interpreter) visitLogical(expr Expr) error {
	if l, ok := expr.(Logical); ok {
		left, err := i.Evaluate(l.Left)
		if err != nil {
			return err
		}

		switch l.Operator.TokenType {
		case Or:
			{
				if !left.Bool() {
					right, err := i.Evaluate(l.Right)
					if err != nil {
						return err
					}

					i.Literal = Literal{right.Bool()}
				} else {
					i.Literal = Literal{true}
				}

				break
			}
		case And:
			{
				if left.Bool() {
					right, err := i.Evaluate(l.Right)
					if err != nil {
						return err
					}

					i.Literal = Literal{right.Bool()}
				} else {
					i.Literal = Literal{false}
				}

				break
			}
		}
	}

	return nil
}

func (i *Interpreter) visitUnary(expr Expr) error {
	if u, ok := expr.(Unary); ok {
		if _, err := i.Evaluate(u.Right); err != nil {
			return err
		}

		invalidOperand := func(operand interface{}) error {
			return fmt.Errorf("error at line %d: bad operand for unary %s: %T", u.Operator.Line, u.Operator.Lexeme, operand)
		}

		switch u.Operator.TokenType {
		case Not:
			{
				i.Literal = Literal{!i.Literal.Bool()}
			}
		case Minus:
			{
				if f, ok := i.Literal.Value.(float64); ok {
					i.Literal = Literal{-f}
				} else {
					return invalidOperand(i.Literal.Value)
				}
			}
		}
	}

	return nil
}

func (i *Interpreter) visitVariable(expr Expr) error {
	if v, ok := expr.(Variable); ok {
		e, err := i.Environment.Get(v, 0)
		if err != nil {
			return err
		}

		if l, ok := e.(Literal); ok {
			i.Literal = l
		}
	}

	return nil
}

func (i *Interpreter) visitWhileStmt(stmt Stmt) error {
	if w, ok := stmt.(WhileStmt); ok {
		for true {
			l, err := i.Evaluate(w.Condition)
			if err != nil {
				return err
			}

			if !l.Bool() {
				return nil
			}

			if err := w.Body.Accept(i); err != nil {
				return err
			}
		}
	}

	return nil
}
