package ast

type ClassInstance struct {
	ClassStmt
	Fields map[string]Literal
}

func (c ClassInstance) Get(t Token) Literal {
	return c.Fields[t.Lexeme]
}

func (c ClassInstance) Set(t Token, l Literal) {
	c.Fields[t.Lexeme] = l
}

func (c ClassInstance) String() string {
	return c.Name.Lexeme
}
