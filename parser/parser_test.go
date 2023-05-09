package parser

import (
	"fmt"
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

func TestReturnStatement(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 993322;
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)
	assert.NotNil(t, program, "ParseProgram() returned nil")
	assert.Equal(t, len(program.Statements), 3, "program.Statements does not contain 3 statements")
	testReturnStatement(t, program.Statements[0], "5")
	testReturnStatement(t, program.Statements[0], "10")
	testReturnStatement(t, program.Statements[0], "993322")
}

func testLetStatement(t *testing.T, s ast.Statement, name string) {
	t.Helper()
	assert.Equal(t, s.TokenLiteral(), "let")
	letStmt, ok := s.(*ast.LetStatement)
	assert.True(t, ok)
	assert.Equal(t, name, letStmt.Name.Value)
	assert.Equal(t, name, letStmt.Name.TokenLiteral())
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	assert.NotNil(t, program, "ParseProgram() returned nil")
	assert.Equal(t, len(program.Statements), 1, "program.Statements does not contain 1 statement")
	testIdentifierExpression(t, program.Statements[0], "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	assert.NotNil(t, program, "ParseProgram() returned nil")
	assert.Equal(t, len(program.Statements), 1, "program.Statements does not contain 1 statement")
	testIntegerLiteralExpression(t, program.Statements[0], "5", 5)
}

func TestPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		integer  int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		assert.NotNil(t, program, "ParseProgram() returned nil")
		assert.Equal(t, len(program.Statements), 1, "program.Statements does not contain 1 statement")
		testPrefixExpression(t, program.Statements[0], tt.operator, tt.integer)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		assert.NotNil(t, program, "ParseProgram() returned nil")
		assert.Equal(t, len(program.Statements), 1, "program.Statements does not contain 1 statement")
		testInfixExpression(t, program.Statements[0], tt.leftValue, tt.operator, tt.rightValue)
	}
}

func testReturnStatement(t *testing.T, s ast.Statement, name string) {
	t.Helper()
	assert.Equal(t, s.TokenLiteral(), "return")
	_, ok := s.(*ast.ReturnStatement)
	assert.True(t, ok)
}

func testIdentifierExpression(t *testing.T, s ast.Statement, name string) {
	t.Helper()
	stmt, ok := s.(*ast.ExpressionStatement)
	assert.True(t, ok)
	ident, ok := stmt.Expression.(*ast.Identifier)
	assert.True(t, ok)
	assert.Equal(t, name, ident.Value)
	assert.Equal(t, name, ident.TokenLiteral())
}

func testIntegerLiteralExpression(t *testing.T, s ast.Statement, name string, value int64) {
	t.Helper()
	stmt, ok := s.(*ast.ExpressionStatement)
	assert.True(t, ok)
	testIntegerLiteral(t, stmt.Expression, value)
}

func testIntegerLiteral(t *testing.T, possiblyAnInteger ast.Expression, value int64) {
	integer, ok := possiblyAnInteger.(*ast.IntegerLiteral)
	assert.True(t, ok)
	assert.Equal(t, value, integer.Value)
	assert.Equal(t, fmt.Sprintf("%d", value), integer.TokenLiteral())
}

func testPrefixExpression(t *testing.T, s ast.Statement, operator string, integer int64) {
	t.Helper()
	stmt, ok := s.(*ast.ExpressionStatement)
	assert.True(t, ok)
	exp, ok := stmt.Expression.(*ast.PrefixExpression)
	assert.True(t, ok)
	assert.Equal(t, operator, exp.Operator)
	testIntegerLiteral(t, exp.Right, integer)
}

func testInfixExpression(t *testing.T, s ast.Statement, leftValue int64, operator string, rightValue int64) {
	t.Helper()
	stmt, ok := s.(*ast.ExpressionStatement)
	assert.True(t, ok)
	exp, ok := stmt.Expression.(*ast.InfixExpression)
	assert.True(t, ok)
	testIntegerLiteral(t, exp.Left, leftValue)
	assert.Equal(t, operator, exp.Operator)
	testIntegerLiteral(t, exp.Right, rightValue)
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
