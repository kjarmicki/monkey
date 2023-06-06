package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"kjarmicki.github.com/monkey/ast"
	"kjarmicki.github.com/monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)
		assert.NotNil(t, program, "ParseProgram() returned nil")
		assert.Equal(t, len(program.Statements), 1, "program.Statements does not contain 1 statement")
		stmt := program.Statements[0]
		testLetStatement(t, stmt, tt.expectedIdentifier)
		val := stmt.(*ast.LetStatement).Value
		testLiteralExpression(t, val, tt.expectedValue)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue any
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return y;", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)
		assert.NotNil(t, program, "ParseProgram() returned nil")
		assert.Equal(t, len(program.Statements), 1, "program.Statements does not contain 1 statement")
		stmt := program.Statements[0]
		testReturnStatement(t, stmt)
		val := stmt.(*ast.ReturnStatement).ReturnValue
		testLiteralExpression(t, val, tt.expectedValue)
	}
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

func TestBooleanExpression(t *testing.T) {
	input := "true;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	assert.NotNil(t, program, "ParseProgram() returned nil")
	assert.Equal(t, len(program.Statements), 1, "program.Statements does not contain 1 statement")
	testBooleanExpression(t, program.Statements[0], true)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	assert.NotNil(t, program, "ParseProgram() returned nil")
	assert.Equal(t, len(program.Statements), 1, "program.Statements does not contain 1 statement")
	testStringLiteralExpression(t, program.Statements[0], "hello world")
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
		leftValue  any
		operator   string
		rightValue any
	}{
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
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
		testInfixStatement(t, program.Statements[0], tt.leftValue, tt.operator, tt.rightValue)
	}
}

func testReturnStatement(t *testing.T, s ast.Statement) {
	t.Helper()
	assert.Equal(t, s.TokenLiteral(), "return")
	_, ok := s.(*ast.ReturnStatement)
	assert.True(t, ok)
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		// grouping
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		assert.Equal(t, tt.expected, actual)
	}
}

func TestIfExpression(t *testing.T) {
	input := "if (x < y) { x }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	assert.NotNil(t, program, "ParseProgram() returned nil")
	assert.Equal(t, len(program.Statements), 1, "program.Statements does not contain 1 statement")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	exp, ok := stmt.Expression.(*ast.IfExpression)
	assert.True(t, ok)
	testInfixExpression(t, exp.Condition, "x", "<", "y")
	assert.Equal(t, len(exp.Consequence.Statements), 1)
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	testIdentifier(t, consequence.Expression, "x")
	assert.Nil(t, exp.Alternative)
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	assert.NotNil(t, program, "ParseProgram() returned nil")
	assert.Equal(t, len(program.Statements), 1, "program.Statements does not contain 1 statement")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	exp, ok := stmt.Expression.(*ast.IfExpression)
	assert.True(t, ok)
	testInfixExpression(t, exp.Condition, "x", "<", "y")

	assert.Equal(t, len(exp.Consequence.Statements), 1)
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	testIdentifier(t, consequence.Expression, "x")

	assert.Equal(t, len(exp.Alternative.Statements), 1)
	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	testIdentifier(t, alternative.Expression, "y")
}

func TestFunctionLiteral(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	assert.NotNil(t, program, "ParseProgram() returned nil")
	assert.Equal(t, len(program.Statements), 1, "program.Statements does not contain 1 statement")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	assert.True(t, ok)
	assert.Equal(t, len(function.Parameters), 2)
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")
	assert.Equal(t, len(function.Body.Statements), 1)
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	assert.NotNil(t, program, "ParseProgram() returned nil")
	assert.Equal(t, len(program.Statements), 1, "program.Statements does not contain 1 statement")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	exp, ok := stmt.Expression.(*ast.CallExpression)
	assert.True(t, ok)
	testIdentifier(t, exp.Function, "add")
	assert.Equal(t, len(exp.Arguments), 3)
	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 3, 4 + 5]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
	exp, ok := stmt.Expression.(*ast.ArrayLiteral)
	assert.True(t, ok)
	assert.Equal(t, 3, len(exp.Elements))
	testIntegerLiteral(t, exp.Elements[0], 1)
	testInfixExpression(t, exp.Elements[1], 2, "*", 3)
	testInfixExpression(t, exp.Elements[2], 4, "+", 5)
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

func testBooleanExpression(t *testing.T, s ast.Statement, value bool) {
	t.Helper()
	stmt, ok := s.(*ast.ExpressionStatement)
	assert.True(t, ok)
	testBoolean(t, stmt.Expression, value)
}

func testBoolean(t *testing.T, possiblyABoolean ast.Expression, value bool) {
	boolean, ok := possiblyABoolean.(*ast.Boolean)
	assert.True(t, ok)
	assert.Equal(t, value, boolean.Value)
}

func testStringLiteralExpression(t *testing.T, s ast.Statement, value string) {
	t.Helper()
	stmt, ok := s.(*ast.ExpressionStatement)
	assert.True(t, ok)
	testString(t, stmt.Expression, value)
}

func testString(t *testing.T, possiblyAString ast.Expression, value string) {
	str, ok := possiblyAString.(*ast.StringLiteral)
	assert.True(t, ok)
	assert.Equal(t, value, str.Value)
}

func testIntegerLiteral(t *testing.T, possiblyAnInteger ast.Expression, value int64) {
	integer, ok := possiblyAnInteger.(*ast.IntegerLiteral)
	assert.True(t, ok)
	assert.Equal(t, value, integer.Value)
	assert.Equal(t, fmt.Sprintf("%d", value), integer.TokenLiteral())
}

func testPrefixExpression(t *testing.T, s ast.Statement, operator string, value any) {
	t.Helper()
	stmt, ok := s.(*ast.ExpressionStatement)
	assert.True(t, ok)
	exp, ok := stmt.Expression.(*ast.PrefixExpression)
	assert.True(t, ok)
	assert.Equal(t, operator, exp.Operator)
	testLiteralExpression(t, exp.Right, value)
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) {
	ident, ok := exp.(*ast.Identifier)
	assert.True(t, ok)
	assert.Equal(t, value, ident.Value)
	assert.Equal(t, value, ident.Token.Literal)
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) {
	switch v := expected.(type) {
	case int:
		testIntegerLiteral(t, exp, int64(v))
		return
	case int64:
		testIntegerLiteral(t, exp, v)
		return
	case string:
		testIdentifier(t, exp, v)
		return
	case bool:
		testBoolean(t, exp, v)
		return
	}
	t.Errorf("type of exp not handled, got %T", exp)
}

func testInfixStatement(t *testing.T, s ast.Statement, leftValue any, operator string, rightValue any) {
	t.Helper()
	stmt, ok := s.(*ast.ExpressionStatement)
	assert.True(t, ok)
	testInfixExpression(t, stmt.Expression, leftValue, operator, rightValue)
}

func testInfixExpression(t *testing.T, ex ast.Expression, leftValue any, operator string, rightValue any) {
	exp, ok := ex.(*ast.InfixExpression)
	assert.True(t, ok)
	testLiteralExpression(t, exp.Left, leftValue)
	assert.Equal(t, operator, exp.Operator)
	testLiteralExpression(t, exp.Right, rightValue)
}

func checkParserErrors(t *testing.T, p *Parser) {
	t.Helper()
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
