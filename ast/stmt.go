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
