package ast

import "fmt"

type Environment struct {
	Parent *Environment
	Scope  map[string]interface{}
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{parent, make(map[string]interface{})}
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

func (e *Environment) Declare(variable Variable, expr Expr) error {
	e.Scope[variable.Lexeme] = expr
	return nil
}
