package ast

import (
	"time"
)

type Callable interface {
	Arity() int
	Call(interpreter *Interpreter, arguments []Expr) (Literal, error)
}

func (f Function) Arity() int {
	return len(f.Arguments)
}

func (f Function) Call(i *Interpreter, arguments []Expr) (Literal, error) {
	i.Environment = NewEnvironment(f.Closure)

	for j, argument := range arguments {
		expr, err := i.Evaluate(argument)
		if err != nil {
			return Literal{}, err
		}

		if err := i.Environment.Declare(Variable{f.Arguments[j]}, expr); err != nil {
			return Literal{}, err
		}
	}

	for _, stmt := range f.Body {
		if err := stmt.Accept(i); err != nil {
			if r, ok := err.(ReturnValue); ok {
				i.Environment = i.Environment.Parent
				return r.Literal, nil
			}

			return Literal{}, err
		}
	}

	i.Environment = i.Environment.Parent

	return Literal{}, nil // void
}

type Clock struct{}

func (c Clock) Arity() int {
	return 0
}

func (c Clock) Call(interpreter *Interpreter, arguments []Expr) (Literal, error) {
	return Literal{time.Now().Unix()}, nil
}
