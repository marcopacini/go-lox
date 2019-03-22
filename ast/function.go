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

func (c ClassStmt) Arity() int {
	return 0
}

func (c ClassStmt) Call(i *Interpreter, arguments []Expr) (Literal, error) {
	return c.CreateInstance(), nil
}

type Clock struct{}

func (c Clock) Arity() int {
	return 0
}

func (c Clock) Call(interpreter *Interpreter, arguments []Expr) (Literal, error) {
	return Literal{time.Now().Unix()}, nil
}
