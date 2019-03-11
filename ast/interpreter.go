package ast

import (
	"fmt"
)

type ExprVisitor interface {
	visitBinary(Expr) error
	visitGrouping(Expr) error
	visitLiteral(Expr) error
	visitUnary(Expr) error
}

type Interpreter struct {
	Literal
}

func (i *Interpreter) Evaluate(expr Expr) (Literal, error) {
	err := expr.Accept(i)
	return i.Literal, err
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
						i.Literal = Literal{l > r}
					} else {
						return invalidOperand(left.Value, right.Value)
					}
				}
			}
		}
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
