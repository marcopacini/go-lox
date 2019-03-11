package ast

type Expr interface {
	Accept(ExprVisitor) error
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

func (l Literal) Accept(visitor ExprVisitor) error {
	return visitor.visitLiteral(l)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (u Unary) Accept(visitor ExprVisitor) error {
	return visitor.visitUnary(u)
}
