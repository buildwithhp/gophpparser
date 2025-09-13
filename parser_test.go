package gophpparser

import (
	"testing"
)

func TestParseSimpleAssignment(t *testing.T) {
	input := `<?php
$name = "John";
?>`

	l := New(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("parser has %d errors", len(p.Errors()))
		for _, err := range p.Errors() {
			t.Errorf("parser error: %q", err)
		}
		return
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T",
			program.Statements[0])
	}

	assignment, ok := stmt.Expression.(*AssignmentExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not AssignmentExpression. got=%T",
			stmt.Expression)
	}

	if assignment.Name.Name != "name" {
		t.Errorf("assignment.Name.Name not 'name'. got=%s", assignment.Name.Name)
	}

	stringLit, ok := assignment.Value.(*StringLiteral)
	if !ok {
		t.Fatalf("assignment.Value is not StringLiteral. got=%T", assignment.Value)
	}

	if stringLit.Value != "John" {
		t.Errorf("stringLit.Value not 'John'. got=%s", stringLit.Value)
	}
}

func TestParseFunctionDeclaration(t *testing.T) {
	input := `<?php
function add($a, $b) {
    return $a + $b;
}
?>`

	l := New(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("parser has %d errors", len(p.Errors()))
		for _, err := range p.Errors() {
			t.Errorf("parser error: %q", err)
		}
		return
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*FunctionDeclaration)
	if !ok {
		t.Fatalf("program.Statements[0] is not FunctionDeclaration. got=%T",
			program.Statements[0])
	}

	if stmt.Name.Value != "add" {
		t.Errorf("stmt.Name.Value not 'add'. got=%s", stmt.Name.Value)
	}

	if len(stmt.Parameters) != 2 {
		t.Fatalf("function parameters wrong. want 2, got=%d", len(stmt.Parameters))
	}

	if stmt.Parameters[0].Name != "a" {
		t.Errorf("stmt.Parameters[0].Name not 'a'. got=%s", stmt.Parameters[0].Name)
	}

	if stmt.Parameters[1].Name != "b" {
		t.Errorf("stmt.Parameters[1].Name not 'b'. got=%s", stmt.Parameters[1].Name)
	}
}

func TestParseEchoStatement(t *testing.T) {
	input := `<?php
echo "Hello World";
?>`

	l := New(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("parser has %d errors", len(p.Errors()))
		for _, err := range p.Errors() {
			t.Errorf("parser error: %q", err)
		}
		return
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*EchoStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not EchoStatement. got=%T",
			program.Statements[0])
	}

	if len(stmt.Values) != 1 {
		t.Fatalf("echo statement values wrong. want 1, got=%d", len(stmt.Values))
	}

	stringLit, ok := stmt.Values[0].(*StringLiteral)
	if !ok {
		t.Fatalf("stmt.Values[0] is not StringLiteral. got=%T", stmt.Values[0])
	}

	if stringLit.Value != "Hello World" {
		t.Errorf("stringLit.Value not 'Hello World'. got=%s", stringLit.Value)
	}
}

func TestParseWhileStatement(t *testing.T) {
	input := `<?php
while ($x < 10) {
    echo $x;
}
?>`

	l := New(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("parser has %d errors", len(p.Errors()))
		for _, err := range p.Errors() {
			t.Errorf("parser error: %q", err)
		}
		return
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*WhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *WhileStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Condition == nil {
		t.Fatalf("stmt.Condition is nil")
	}

	if stmt.Body == nil {
		t.Fatalf("stmt.Body is nil")
	}
}

func TestParseForeachStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`<?php
foreach ($arr as $value) {
    echo $value;
}
?>`,
			"foreach",
		},
		{
			`<?php
foreach ($arr as $key => $value) {
    echo $key . $value;
}
?>`,
			"foreach",
		},
	}

	for _, tt := range tests {
		l := New(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			t.Errorf("parser has %d errors", len(p.Errors()))
			for _, err := range p.Errors() {
				t.Errorf("parser error: %q", err)
			}
			continue
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ForeachStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ForeachStatement. got=%T",
				program.Statements[0])
		}

		if stmt.Array == nil {
			t.Fatalf("stmt.Array is nil")
		}

		if stmt.Value == nil {
			t.Fatalf("stmt.Value is nil")
		}

		if stmt.Body == nil {
			t.Fatalf("stmt.Body is nil")
		}
	}
}

func TestParseBreakContinueStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"<?php break; ?>", "break"},
		{"<?php continue; ?>", "continue"},
		{"<?php break 2; ?>", "break"},
		{"<?php continue 1; ?>", "continue"},
	}

	for _, tt := range tests {
		l := New(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			t.Errorf("parser has %d errors", len(p.Errors()))
			for _, err := range p.Errors() {
				t.Errorf("parser error: %q", err)
			}
			continue
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]

		switch tt.expected {
		case "break":
			if _, ok := stmt.(*BreakStatement); !ok {
				t.Fatalf("stmt is not *BreakStatement. got=%T", stmt)
			}
		case "continue":
			if _, ok := stmt.(*ContinueStatement); !ok {
				t.Fatalf("stmt is not *ContinueStatement. got=%T", stmt)
			}
		}
	}
}

func TestParseAssociativeArray(t *testing.T) {
	input := `<?php
$arr = ["name" => "John", "age" => 30];
?>`

	l := New(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("parser has %d errors", len(p.Errors()))
		for _, err := range p.Errors() {
			t.Errorf("parser error: %q", err)
		}
		return
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ExpressionStatement. got=%T",
			program.Statements[0])
	}

	assignExpr, ok := stmt.Expression.(*AssignmentExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *AssignmentExpression. got=%T",
			stmt.Expression)
	}

	assocArray, ok := assignExpr.Value.(*AssociativeArrayLiteral)
	if !ok {
		t.Fatalf("assignExpr.Value is not *AssociativeArrayLiteral. got=%T",
			assignExpr.Value)
	}

	if len(assocArray.Pairs) != 2 {
		t.Fatalf("assocArray.Pairs length is not 2. got=%d", len(assocArray.Pairs))
	}
}

func TestParseClassDeclaration(t *testing.T) {
	input := `<?php
class User {
    public $name;
    private $age;
    
    public function getName() {
        return $this->name;
    }
    
    private static function validate($data) {
        return true;
    }
}
?>`

	l := New(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("parser has %d errors", len(p.Errors()))
		for _, err := range p.Errors() {
			t.Errorf("parser error: %q", err)
		}
		return
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ClassDeclaration)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ClassDeclaration. got=%T",
			program.Statements[0])
	}

	if stmt.Name.Value != "User" {
		t.Errorf("class name not 'User'. got=%s", stmt.Name.Value)
	}

	if len(stmt.Properties) != 2 {
		t.Errorf("class properties length not 2. got=%d", len(stmt.Properties))
	}

	if len(stmt.Methods) != 2 {
		t.Errorf("class methods length not 2. got=%d", len(stmt.Methods))
	}
}

func TestParseClassInheritance(t *testing.T) {
	input := `<?php
class Admin extends User {
    public $permissions;
}
?>`

	l := New(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("parser has %d errors", len(p.Errors()))
		for _, err := range p.Errors() {
			t.Errorf("parser error: %q", err)
		}
		return
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ClassDeclaration)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ClassDeclaration. got=%T",
			program.Statements[0])
	}

	if stmt.SuperClass == nil {
		t.Fatalf("stmt.SuperClass is nil")
	}

	if stmt.SuperClass.Value != "User" {
		t.Errorf("superclass name not 'User'. got=%s", stmt.SuperClass.Value)
	}
}

func TestParseNewExpression(t *testing.T) {
	input := `<?php
$user = new User("John", 25);
?>`

	l := New(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("parser has %d errors", len(p.Errors()))
		for _, err := range p.Errors() {
			t.Errorf("parser error: %q", err)
		}
		return
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ExpressionStatement. got=%T",
			program.Statements[0])
	}

	assignExpr, ok := stmt.Expression.(*AssignmentExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *AssignmentExpression. got=%T",
			stmt.Expression)
	}

	newExpr, ok := assignExpr.Value.(*NewExpression)
	if !ok {
		t.Fatalf("assignExpr.Value is not *NewExpression. got=%T",
			assignExpr.Value)
	}

	if newExpr.ClassName.Value != "User" {
		t.Errorf("class name not 'User'. got=%s", newExpr.ClassName.Value)
	}

	if len(newExpr.Arguments) != 2 {
		t.Errorf("arguments length not 2. got=%d", len(newExpr.Arguments))
	}
}

func TestParseObjectAccess(t *testing.T) {
	input := `<?php
$name = $user->getName();
$age = $user->age;
?>`

	l := New(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("parser has %d errors", len(p.Errors()))
		for _, err := range p.Errors() {
			t.Errorf("parser error: %q", err)
		}
		return
	}

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d",
			len(program.Statements))
	}

	// Test method call
	stmt1, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ExpressionStatement. got=%T",
			program.Statements[0])
	}

	assignExpr1, ok := stmt1.Expression.(*AssignmentExpression)
	if !ok {
		t.Fatalf("stmt1.Expression is not *AssignmentExpression. got=%T",
			stmt1.Expression)
	}

	callExpr, ok := assignExpr1.Value.(*CallExpression)
	if !ok {
		t.Fatalf("assignExpr1.Value is not *CallExpression. got=%T",
			assignExpr1.Value)
	}

	objAccess, ok := callExpr.Function.(*ObjectAccessExpression)
	if !ok {
		t.Fatalf("callExpr.Function is not *ObjectAccessExpression. got=%T",
			callExpr.Function)
	}

	if objAccess.Object.(*Variable).Name != "user" {
		t.Errorf("object name not 'user'. got=%s", objAccess.Object.(*Variable).Name)
	}
}

func TestParseStaticAccess(t *testing.T) {
	input := `<?php
$result = User::validate($data);
?>`

	l := New(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("parser has %d errors", len(p.Errors()))
		for _, err := range p.Errors() {
			t.Errorf("parser error: %q", err)
		}
		return
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ExpressionStatement. got=%T",
			program.Statements[0])
	}

	assignExpr, ok := stmt.Expression.(*AssignmentExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *AssignmentExpression. got=%T",
			stmt.Expression)
	}

	callExpr, ok := assignExpr.Value.(*CallExpression)
	if !ok {
		t.Fatalf("assignExpr.Value is not *CallExpression. got=%T",
			assignExpr.Value)
	}

	staticAccess, ok := callExpr.Function.(*StaticAccessExpression)
	if !ok {
		t.Fatalf("callExpr.Function is not *StaticAccessExpression. got=%T",
			callExpr.Function)
	}

	if staticAccess.Class.(*Identifier).Value != "User" {
		t.Errorf("class name not 'User'. got=%s", staticAccess.Class.(*Identifier).Value)
	}
}

func TestParseNamespaceDeclaration(t *testing.T) {
	input := `<?php
namespace App;
?>`

	l := New(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("parser has %d errors", len(p.Errors()))
		for _, err := range p.Errors() {
			t.Errorf("parser error: %q", err)
		}
		return
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*NamespaceDeclaration)
	if !ok {
		t.Fatalf("program.Statements[0] is not *NamespaceDeclaration. got=%T",
			program.Statements[0])
	}

	if stmt.Name.Value != "App" { // Simple namespace for now
		t.Errorf("namespace name not 'App'. got=%s", stmt.Name.Value)
	}
}

func TestParseUseStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		alias    string
	}{
		{`<?php use Models; ?>`, "Models", ""},
		{`<?php use Models as UserModel; ?>`, "Models", "UserModel"},
	}

	for _, tt := range tests {
		l := New(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			t.Errorf("parser has %d errors", len(p.Errors()))
			for _, err := range p.Errors() {
				t.Errorf("parser error: %q", err)
			}
			continue
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*UseStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *UseStatement. got=%T",
				program.Statements[0])
		}

		if stmt.Namespace.Value != tt.expected {
			t.Errorf("namespace name not '%s'. got=%s", tt.expected, stmt.Namespace.Value)
		}

		if tt.alias != "" {
			if stmt.Alias == nil {
				t.Errorf("expected alias '%s' but got nil", tt.alias)
			} else if stmt.Alias.Value != tt.alias {
				t.Errorf("alias not '%s'. got=%s", tt.alias, stmt.Alias.Value)
			}
		}
	}
}

func TestParseTryStatement(t *testing.T) {
	input := `<?php
try {
    $result = risky_operation();
} catch (Exception $e) {
    echo $e->getMessage();
} finally {
    cleanup();
}
?>`

	l := New(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("parser has %d errors", len(p.Errors()))
		for _, err := range p.Errors() {
			t.Errorf("parser error: %q", err)
		}
		return
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*TryStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *TryStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Body == nil {
		t.Fatalf("stmt.Body is nil")
	}

	if len(stmt.Catches) != 1 {
		t.Fatalf("expected 1 catch clause. got=%d", len(stmt.Catches))
	}

	if stmt.Finally == nil {
		t.Fatalf("stmt.Finally is nil")
	}

	catch := stmt.Catches[0]
	if catch.ExceptionType.Value != "Exception" {
		t.Errorf("exception type not 'Exception'. got=%s", catch.ExceptionType.Value)
	}

	if catch.Variable.Name != "e" {
		t.Errorf("exception variable not 'e'. got=%s", catch.Variable.Name)
	}
}

func TestParseThrowStatement(t *testing.T) {
	input := `<?php
throw new Exception("Error message");
?>`

	l := New(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("parser has %d errors", len(p.Errors()))
		for _, err := range p.Errors() {
			t.Errorf("parser error: %q", err)
		}
		return
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ThrowStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ThrowStatement. got=%T",
			program.Statements[0])
	}

	newExpr, ok := stmt.Expression.(*NewExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *NewExpression. got=%T",
			stmt.Expression)
	}

	if newExpr.ClassName.Value != "Exception" {
		t.Errorf("exception class not 'Exception'. got=%s", newExpr.ClassName.Value)
	}
}

func TestParseAnonymousFunction(t *testing.T) {
	input := `<?php
$callback = function($x, $y) use ($multiplier) {
    return $x * $y * $multiplier;
};
?>`

	l := New(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("parser has %d errors", len(p.Errors()))
		for _, err := range p.Errors() {
			t.Errorf("parser error: %q", err)
		}
		return
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ExpressionStatement. got=%T",
			program.Statements[0])
	}

	assignExpr, ok := stmt.Expression.(*AssignmentExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *AssignmentExpression. got=%T",
			stmt.Expression)
	}

	anonFunc, ok := assignExpr.Value.(*AnonymousFunction)
	if !ok {
		t.Fatalf("assignExpr.Value is not *AnonymousFunction. got=%T",
			assignExpr.Value)
	}

	if len(anonFunc.Parameters) != 2 {
		t.Errorf("anonymous function parameters length not 2. got=%d", len(anonFunc.Parameters))
	}

	if len(anonFunc.UseClause) != 1 {
		t.Errorf("anonymous function use clause length not 1. got=%d", len(anonFunc.UseClause))
	}

	if anonFunc.Body == nil {
		t.Fatalf("anonymous function body is nil")
	}
}

func TestParseYieldExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`<?php yield; ?>`, "yield"},
		{`<?php yield $value; ?>`, "yield"},
		{`<?php yield $key => $value; ?>`, "yield"},
	}

	for _, tt := range tests {
		l := New(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			t.Errorf("parser has %d errors", len(p.Errors()))
			for _, err := range p.Errors() {
				t.Errorf("parser error: %q", err)
			}
			continue
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ExpressionStatement. got=%T",
				program.Statements[0])
		}

		_, ok = stmt.Expression.(*YieldExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not *YieldExpression. got=%T",
				stmt.Expression)
		}
	}
}

func TestParseInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		input := "<?php " + tt.input + "; ?>"
		l := New(input)
		p := NewParser(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			t.Errorf("parser has %d errors", len(p.Errors()))
			for _, err := range p.Errors() {
				t.Errorf("parser error: %q", err)
			}
			continue
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*InfixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not InfixExpression. got=%T", stmt.Expression)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il Expression, value int64) bool {
	integ, ok := il.(*IntegerLiteral)
	if !ok {
		t.Errorf("il not *IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	return true
}
