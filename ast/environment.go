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

type Environment struct {
	Parent *Environment
	Scope  map[string]interface{}
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{parent, make(map[string]interface{})}
}

func (e *Environment) Assign(variable Variable, expr Expr) error {
	if _, ok := e.Scope[variable.Lexeme]; ok {
		e.Scope[variable.Lexeme] = expr
		return nil
	}

	if e.Parent != nil {
		return e.Parent.Assign(variable, expr)
	}

	return fmt.Errorf("error at line %d: undefined variable %v", variable.Line, variable.Lexeme)
}

func (e Environment) Contains(variable Variable) bool {
	if _, ok := e.Scope[variable.Lexeme]; ok {
		return true
	}

	if e.Parent != nil {
		return e.Parent.Contains(variable)
	}

	return false
}

func (e *Environment) Declare(variable Variable, expr Expr) error {
	e.Scope[variable.Lexeme] = expr
	return nil
}

func (e Environment) Get(variable Variable, distance int) (interface{}, error) {
	local := &e

	for i := 0; i < distance; i++ {
		local = local.Parent
	}

	if expr, ok := local.Scope[variable.Lexeme]; ok {
		return expr, nil
	}

	if local.Parent != nil {
		return local.Parent.Get(variable, 0)
	}

	return nil, fmt.Errorf("error at line %d: undefined variable %v", variable.Line, variable.Lexeme)
}

func (e *Environment) Set(name string, callable Callable) {
	e.Scope[name] = Literal{callable}
}
