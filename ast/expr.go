package ast

import (
	"fmt"
	"math"
)

type Expr interface {
	Accept(ExprVisitor) error
}

type ExprVisitor interface {
	visitAssign(Expr) error
	visitBinary(Expr) error
	visitGrouping(Expr) error
	visitLiteral(Expr) error
	visitLogical(Expr) error
	visitUnary(Expr) error
	visitVariable(Expr) error
}

type Assign struct {
	Variable
	Token
	Expr
}

func (a Assign) Accept(visitor ExprVisitor) error {
	return visitor.visitAssign(a)
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b Binary) Accept(visitor ExprVisitor) error {
	return visitor.visitBinary(b)
}

type Grouping struct {
	Expr
}

func (g Grouping) Accept(visitor ExprVisitor) error {
	return visitor.visitGrouping(g)
}

type Literal struct {
	Value interface{}
}

func (l Literal) Bool() bool {
	if b, ok := l.Value.(bool); ok {
		return b
	} else {
		if l.Value == nil {
			return false
		} else {
			return true
		}
	}
}

func (l Literal) String() string {
	if s, ok := l.Value.(string); ok {
		return s
	}

	if f, ok := l.Value.(float64); ok {
		if f == math.Trunc(f) {
			return fmt.Sprintf("%d", int64(f))
		}

		return fmt.Sprintf("%f", f)
	}

	if b, ok := l.Value.(bool); ok {
		return fmt.Sprintf("%t", b)
	}

	if l.Value == nil {
		return "nil"
	}

	return fmt.Sprintf("%v", l.Value)
}

func (l Literal) Accept(visitor ExprVisitor) error {
	return visitor.visitLiteral(l)
}

type Logical struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (l Logical) Accept(visitor ExprVisitor) error {
	return visitor.visitLogical(l)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (u Unary) Accept(visitor ExprVisitor) error {
	return visitor.visitUnary(u)
}

type Variable struct {
	Token
}

func (v Variable) Accept(visitor ExprVisitor) error {
	return visitor.visitVariable(v)
}
