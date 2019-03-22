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

type Stmt interface {
	Accept(StmtVisitor) error
}

type StmtVisitor interface {
	visitBlock(Block) error
	visitClassStmt(ClassStmt) error
	visitDeclaration(Declaration) error
	visitForStmt(ForStmt) error
	visitFunction(Function) error
	visitIfStmt(IfStmt) error
	visitExprStmt(ExprStmt) error
	visitPrintStmt(PrintStmt) error
	visitReturnStmt(ReturnStmt) error
	visitWhileStmt(WhileStmt) error
}

type Block struct {
	Stmts []Stmt
}

func (b Block) Accept(visitor StmtVisitor) error {
	return visitor.visitBlock(b)
}

type ClassStmt struct {
	Name    Token
	Methods []Function
}

func (c ClassStmt) Accept(visitor StmtVisitor) error {
	return visitor.visitClassStmt(c)
}

func (c ClassStmt) CreateInstance() Literal {
	return Literal{ClassInstance{c, make(map[string]Literal)}}
}

type Declaration struct {
	Token
	Expr
}

func (d Declaration) Accept(visitor StmtVisitor) error {
	return visitor.visitDeclaration(d)
}

type ForStmt struct {
	Init      Stmt
	Condition Expr
	Increment Expr
	Body      Stmt
}

func (f ForStmt) Accept(visitor StmtVisitor) error {
	return visitor.visitForStmt(f)
}

type IfStmt struct {
	Condition Expr
	Then      Stmt
	Else      Stmt
}

func (i IfStmt) Accept(visitor StmtVisitor) error {
	return visitor.visitIfStmt(i)
}

type ExprStmt struct {
	Expr
}

func (e ExprStmt) Accept(visitor StmtVisitor) error {
	return visitor.visitExprStmt(e)
}

type Function struct {
	Name      Token
	Closure   *Environment
	Arguments []Token
	Body      []Stmt
}

func (f Function) Accept(visitor StmtVisitor) error {
	return visitor.visitFunction(f)
}

type PrintStmt struct {
	Expr
}

func (p PrintStmt) Accept(visitor StmtVisitor) error {
	return visitor.visitPrintStmt(p)
}

type ReturnStmt struct {
	Expr
}

func (r ReturnStmt) Accept(visitor StmtVisitor) error {
	return visitor.visitReturnStmt(r)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (w WhileStmt) Accept(visitor StmtVisitor) error {
	return visitor.visitWhileStmt(w)
}
