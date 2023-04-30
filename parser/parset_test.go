package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"kjarmicki.github.com/monkey/ast"
	"kjarmicki.github.com/monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	assert.NotNil(t, program, "ParseProgram() returned nil")
	assert.Equal(t, len(program.Statements), 3, "program.Statements does not contain 3 statements")

	testLetStatement(t, program.Statements[0], "x")
	testLetStatement(t, program.Statements[1], "y")
	testLetStatement(t, program.Statements[2], "foobar")
}

func testLetStatement(t *testing.T, s ast.Statement, name string) {
	t.Helper()
	assert.Equal(t, s.TokenLiteral(), "let")
	letStmt, ok := s.(*ast.LetStatement)
	assert.True(t, ok)
	assert.Equal(t, name, letStmt.Name.Value)
	assert.Equal(t, name, letStmt.Name.TokenLiteral())
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
