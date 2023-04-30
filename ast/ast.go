package ast

import "kjarmicki.github.com/monkey/token"

/*
 * AST is an internal data representation of the source code. AST is the result of parsing (sytactic analysis).
 */

type Node interface {
	TokenLiteral() string
}

// expression is a value, or anything that executes and in the end produces a value (e.g. 3 + 5)
type Expression interface {
	Node
	expressionNode()
}

// statements are built from expressions (e.g. let x = 3 + 5)
type Statement interface {
	Node
	statementNode()
}

// identifier is an expression even though it doesn't produce a value to keep things simple
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
