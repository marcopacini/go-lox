package ast

import "fmt"

type Scope map[string]bool

func NewScope() Scope {
	return make(map[string]bool, 0)
}

type Stack struct {
	stack []Scope
}

func NewStack() Stack {
	return Stack{make([]Scope, 0)}
}

func (s *Stack) Head() (Scope, bool) {
	if len(s.stack) <= 0 {
		return nil, false
	}

	return s.stack[len(s.stack)-1], true
}

func (s *Stack) Push(scope Scope) {
	s.stack = append(s.stack, scope)
}

func (s *Stack) Pop() (Scope, bool) {
	if len(s.stack) <= 0 {
		return nil, false
	}

	last := len(s.stack) - 1

	scope := s.stack[last]
	s.stack = s.stack[:last]

	return scope, false
}

func (s *Stack) Declare(name string) {
	if s, ok := s.Head(); ok {
		s[name] = false
	}
}

func (s *Stack) Define(name string) {
	if s, ok := s.Head(); ok {
		s[name] = true
	}
}

type Resolver struct {
	Stack
	Locals map[string]int
}

func (r *Resolver) Resolve(stmts []Stmt) error {
	r.Stack = NewStack()
	r.Stack.Push(NewScope())
	r.Locals = make(map[string]int, 0)

	for _, stmt := range stmts {
		if err := stmt.Accept(r); err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) beginScope() {
	r.Stack.Push(NewScope())
}

func (r *Resolver) endScope() {
	r.Stack.Pop()
}

func (r *Resolver) visitAssign(a Assign) error {
	if err := a.Variable.Accept(r); err != nil {
		return err
	}

	if err := a.Expr.Accept(r); err != nil {
		return err
	}

	return nil
}

func (r *Resolver) visitBinary(b Binary) error {
	if err := b.Left.Accept(r); err != nil {
		return err
	}

	if err := b.Right.Accept(r); err != nil {
		return err
	}

	return nil
}

func (r *Resolver) visitBlock(b Block) error {
	r.beginScope()
	for _, stmt := range b.Stmts {
		if err := stmt.Accept(r); err != nil {
			return err
		}
	}
	r.endScope()

	return nil
}

func (r *Resolver) visitCall(c Call) error {
	if err := c.Callee.Accept(r); err != nil {
		return nil
	}

	for _, expr := range c.Arguments {
		if err := expr.Accept(r); err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) visitDeclaration(d Declaration) error {
	r.Stack.Declare(d.Lexeme)
	if d.Expr != nil {
		if err := d.Expr.Accept(r); err != nil {
			return err
		}
	}
	r.Stack.Define(d.Lexeme)

	return nil
}

func (r *Resolver) visitExprStmt(e ExprStmt) error {
	if err := e.Expr.Accept(r); err != nil {
		return err
	}

	return nil
}

func (r *Resolver) visitForStmt(f ForStmt) error {
	if f.Init != nil {
		if err := f.Init.Accept(r); err != nil {
			return err
		}
	}

	if f.Condition != nil {
		if err := f.Condition.Accept(r); err != nil {
			return err
		}
	}

	if f.Increment != nil {
		if err := f.Increment.Accept(r); err != nil {
			return err
		}
	}

	if err := f.Body.Accept(r); err != nil {
		return err
	}

	return nil
}

func (r *Resolver) visitFunction(f Function) error {
	r.beginScope()
	for _, argument := range f.Arguments {
		r.Stack.Declare(argument.Lexeme)
		r.Stack.Define(argument.Lexeme)
	}

	for _, stmt := range f.Body {
		if err := stmt.Accept(r); err != nil {
			return err
		}
	}
	r.endScope()

	return nil
}

func (r *Resolver) visitGrouping(g Grouping) error {
	if err := g.Expr.Accept(r); err != nil {
		return err
	}

	return nil
}

func (r *Resolver) visitIfStmt(i IfStmt) error {
	if err := i.Condition.Accept(r); err != nil {
		return err
	}

	if err := i.Then.Accept(r); err != nil {
		return err
	}

	if i.Else != nil {
		if err := i.Else.Accept(r); err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) visitLiteral(l Literal) error {
	return nil
}

func (r *Resolver) visitLogical(l Logical) error {
	if err := l.Left.Accept(r); err != nil {
		return err
	}

	if err := l.Right.Accept(r); err != nil {
		return err
	}

	return nil
}

func (r *Resolver) visitPrintStmt(p PrintStmt) error {
	if err := p.Expr.Accept(r); err != nil {
		return err
	}

	return nil
}

func (r *Resolver) visitReturnStmt(s ReturnStmt) error {
	if err := s.Expr.Accept(r); err != nil {
		return err
	}

	return nil
}

func (r *Resolver) visitUnary(u Unary) error {
	if err := u.Right.Accept(r); err != nil {
		return err
	}

	return nil
}

func (r *Resolver) visitVariable(v Variable) error {
	if s, ok := r.Stack.Head(); ok {
		if b, ok := s[v.Lexeme]; ok && !b {
			return fmt.Errorf("error at line %d: cannot read local variable in its own initializer\n", v.Line)
		}
	}

	for i := len(r.stack) - 1; i >= 0; i-- {
		if _, ok := r.stack[i][v.Lexeme]; ok {
			r.Locals[fmt.Sprintf("%v%v", v.Lexeme, v.Line)] = len(r.stack) - 1 - i
			break
		}
	}

	return nil
}

func (r *Resolver) visitWhileStmt(w WhileStmt) error {
	if err := w.Condition.Accept(r); err != nil {
		return err
	}

	if err := w.Body.Accept(r); err != nil {
		return err
	}

	return nil
}
