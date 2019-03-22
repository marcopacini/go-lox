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

import (
	"fmt"
)

type Interpreter struct {
	Literal
	Locals map[string]int
	*Environment
}

type ReturnValue struct {
	Literal
}

func (r ReturnValue) Error() string {
	return r.Literal.String()
}

func (i *Interpreter) Run(stmts []Stmt) error {
	r := Resolver{}

	if err := r.Resolve(stmts); err != nil {
		return err
	}

	i.Locals = r.Locals

	i.Environment = NewEnvironment(nil)
	i.Environment.Set("clock", Clock{})

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

func (i *Interpreter) visitAssign(a Assign) error {
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

	return nil
}

func (i *Interpreter) visitBinary(b Binary) error {
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

	return nil
}

func (i *Interpreter) visitBlock(b Block) error {
	i.Environment = NewEnvironment(i.Environment)

	for _, stmt := range b.Stmts {
		if err := stmt.Accept(i); err != nil {
			return err
		}
	}

	i.Environment = i.Environment.Parent

	return nil
}

func (i *Interpreter) visitCall(c Call) error {
	callee, err := i.Evaluate(c.Callee)
	if err != nil {
		return err
	}

	var arguments []Expr
	for _, argument := range c.Arguments {
		value, err := i.Evaluate(argument)
		if err != nil {
			return err
		}

		arguments = append(arguments, value)
	}

	if f, ok := callee.Value.(Callable); ok {
		if f.Arity() != len(arguments) {
			return fmt.Errorf("expected %d arguments but got %d", f.Arity(), len(arguments))
		}

		l, err := f.Call(i, arguments)
		if err != nil {
			return err
		}

		i.Literal = l
	}

	return nil
}

func (i *Interpreter) visitClassStmt(c ClassStmt) error {
	return i.Environment.Declare(Variable{c.Name}, Literal{c})
}

func (i *Interpreter) visitDeclaration(d Declaration) error {
	i.Literal = Literal{nil}

	if d.Expr != nil {
		if _, err := i.Evaluate(d.Expr); err != nil {
			return err
		}
	}

	if err := i.Environment.Declare(Variable{d.Token}, i.Literal); err != nil {
		return err
	}

	return nil
}

func (i *Interpreter) visitExprStmt(e ExprStmt) error {
	return e.Expr.Accept(i)
}

func (i *Interpreter) visitForStmt(f ForStmt) error {
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

	return nil
}

func (i *Interpreter) visitFunction(f Function) error {
	f.Closure = i.Environment
	if err := i.Environment.Declare(Variable{f.Name}, Literal{f}); err != nil {
		return err
	}

	return nil
}

func (i *Interpreter) visitIfStmt(s IfStmt) error {
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

	return nil
}

func (i *Interpreter) visitPrintStmt(p PrintStmt) error {
	expr, err := i.Evaluate(p.Expr)
	if err != nil {
		return err
	}

	fmt.Println(expr)

	return nil
}

func (i *Interpreter) visitGet(g Get) error {
	l, err := i.Evaluate(g.Object)
	if err != nil {
		return nil
	}

	if obj, ok := l.Value.(ClassInstance); ok {
		i.Literal = obj.Get(g.Name)
	} else {
		return fmt.Errorf("error at line %d: invalid property: %v", g.Name.Line, g.Name.Lexeme)
	}

	return nil
}

func (i *Interpreter) visitGrouping(g Grouping) error {
	return g.Expr.Accept(i)
}

func (i *Interpreter) visitLiteral(l Literal) error {
	i.Literal = l
	return nil
}

func (i *Interpreter) visitLogical(l Logical) error {
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

	return nil
}

func (i *Interpreter) visitReturnStmt(r ReturnStmt) error {
	if err := r.Expr.Accept(i); err != nil {
		return err
	} else {
		return ReturnValue{i.Literal}
	}
}

func (i *Interpreter) visitSet(s Set) error {
	l, err := i.Evaluate(s.Object)
	if err != nil {
		return nil
	}

	if obj, ok := l.Value.(ClassInstance); ok {
		l, err := i.Evaluate(s.Value)
		if err != nil {
			return err
		}

		obj.Set(s.Name, l)
	}

	return nil
}

func (i *Interpreter) visitUnary(u Unary) error {
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

	return nil
}

func (i *Interpreter) visitVariable(v Variable) error {
	distance, _ := i.Locals[fmt.Sprintf("%v%v", v.Lexeme, v.Line)]

	e, err := i.Environment.Get(v, distance)
	if err != nil {
		return err
	}

	if l, ok := e.(Literal); ok {
		i.Literal = l
	}

	return nil
}

func (i *Interpreter) visitWhileStmt(w WhileStmt) error {
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

	return nil
}
