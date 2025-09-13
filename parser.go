package gophpparser

import (
	"fmt"
	"strconv"
	"strings"
	"os"
)

const (
	_ int = iota
	LOWEST
	TERNARY     // ? :
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[TokenType]int{
	QUESTION:                 TERNARY,
	QUESTION_QUESTION:        EQUALS,
	QUESTION_QUESTION_ASSIGN: EQUALS,
	QUESTION_ARROW:           CALL,
	EQ:                       EQUALS,
	NOT_EQ:                   EQUALS,
	LT:                       LESSGREATER,
	GT:                       LESSGREATER,
	LTE:                      LESSGREATER,
	GTE:                      LESSGREATER,
	SPACESHIP:                LESSGREATER,
	PLUS:                     SUM,
	MINUS:                    SUM,
	CONCAT:                   SUM,
	DIVIDE:                   PRODUCT,
	MULTIPLY:                 PRODUCT,
	MODULO:                   PRODUCT,
	LPAREN:                   CALL,
	OBJECT_ACCESS:            CALL,
	STATIC_ACCESS:            CALL,
}

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

type Parser struct {
	l *Lexer

	curToken  Token
	peekToken Token

	errors []string

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.registerPrefix(IDENT, p.parseIdentifier)
	p.registerPrefix(VARIABLE, p.parseVariable)
	p.registerPrefix(INT, p.parseIntegerLiteral)
	p.registerPrefix(FLOAT, p.parseFloatLiteral)
	p.registerPrefix(STRING, p.parseStringLiteral)
	p.registerPrefix(TRUE, p.parseBooleanLiteral)
	p.registerPrefix(FALSE, p.parseBooleanLiteral)
	p.registerPrefix(NOT, p.parsePrefixExpression)
	p.registerPrefix(MINUS, p.parsePrefixExpression)
	p.registerPrefix(INCREMENT, p.parsePrefixExpression)
	p.registerPrefix(DECREMENT, p.parsePrefixExpression)
	p.registerPrefix(NEW, p.parseNewExpression)
	p.registerPrefix(FUNCTION, p.parseAnonymousFunction)
	p.registerPrefix(YIELD, p.parseYieldExpression)
	p.registerPrefix(LPAREN, p.parseGroupedExpression)
	p.registerPrefix(LBRACKET, p.parseArrayLiteral)

	p.infixParseFns = make(map[TokenType]infixParseFn)
	p.registerInfix(PLUS, p.parseInfixExpression)
	p.registerInfix(MINUS, p.parseInfixExpression)
	p.registerInfix(MULTIPLY, p.parseInfixExpression)
	p.registerInfix(DIVIDE, p.parseInfixExpression)
	p.registerInfix(MODULO, p.parseInfixExpression)
	p.registerInfix(CONCAT, p.parseInfixExpression)
	p.registerInfix(EQ, p.parseInfixExpression)
	p.registerInfix(NOT_EQ, p.parseInfixExpression)
	p.registerInfix(LT, p.parseInfixExpression)
	p.registerInfix(GT, p.parseInfixExpression)
	p.registerInfix(LTE, p.parseInfixExpression)
	p.registerInfix(GTE, p.parseInfixExpression)
	p.registerInfix(SPACESHIP, p.parseInfixExpression)
	p.registerInfix(AND, p.parseInfixExpression)
	p.registerInfix(OR, p.parseInfixExpression)
	p.registerInfix(QUESTION, p.parseTernaryExpression)
	p.registerInfix(QUESTION_QUESTION, p.parseInfixExpression)
	p.registerInfix(QUESTION_QUESTION_ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(QUESTION_ARROW, p.parseObjectAccessExpression)
	p.registerInfix(ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(LPAREN, p.parseCallExpression)
	p.registerInfix(LBRACKET, p.parseIndexExpression)
	p.registerInfix(INCREMENT, p.parsePostfixExpression)
	p.registerInfix(DECREMENT, p.parsePostfixExpression)
	p.registerInfix(OBJECT_ACCESS, p.parseObjectAccessExpression)
	p.registerInfix(STATIC_ACCESS, p.parseStaticAccessExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for !p.curTokenIs(EOF) {
		if p.curTokenIs(PHP_OPEN) {
			p.nextToken()
			continue
		}
		if p.curTokenIs(PHP_CLOSE) {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case FUNCTION:
		return p.parseFunctionDeclaration()
	case CLASS:
		return p.parseClassDeclaration()
	case INTERFACE:
		return p.parseInterfaceDeclaration()
	case TRAIT:
		return p.parseTraitDeclaration()
	case CONST:
		return p.parseConstantDeclaration()
	case NAMESPACE:
		return p.parseNamespaceDeclaration()
	case USE:
		return p.parseUseStatement()
	case TRY:
		return p.parseTryStatement()
	case THROW:
		return p.parseThrowStatement()
	case RETURN:
		return p.parseReturnStatement()
	case IF:
		return p.parseIfStatement()
	case ECHO:
		return p.parseEchoStatement()
	case FOR:
		return p.parseForStatement()
	case WHILE:
		return p.parseWhileStatement()
	case FOREACH:
		return p.parseForeachStatement()
	case BREAK:
		return p.parseBreakStatement()
	case CONTINUE:
		return p.parseContinueStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseFunctionDeclaration() *FunctionDeclaration {
	stmt := &FunctionDeclaration{Token: p.curToken}

	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseFunctionParameters() []*Variable {
	identifiers := []*Variable{}

	if p.peekTokenIs(RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	if p.curToken.Type == VARIABLE {
		ident := &Variable{Token: p.curToken, Name: p.curToken.Literal[1:]}
		identifiers = append(identifiers, ident)
	}

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		if p.curToken.Type == VARIABLE {
			ident := &Variable{Token: p.curToken, Name: p.curToken.Literal[1:]}
			identifiers = append(identifiers, ident)
		}
	}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Token: p.curToken}
	block.Statements = []Statement{}

	p.nextToken()

	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Token: p.curToken}

	p.nextToken()

	if !p.curTokenIs(SEMICOLON) {
		stmt.ReturnValue = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIfStatement() *IfStatement {
	stmt := &IfStatement{Token: p.curToken}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(RPAREN) {
		return nil
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	stmt.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(ELSE) {
		p.nextToken()

		if !p.expectPeek(LBRACE) {
			return nil
		}

		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseEchoStatement() *EchoStatement {
	stmt := &EchoStatement{Token: p.curToken}
	stmt.Values = []Expression{}

	p.nextToken()
	stmt.Values = append(stmt.Values, p.parseExpression(LOWEST))

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		stmt.Values = append(stmt.Values, p.parseExpression(LOWEST))
	}

	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}

	// Check for assignment patterns: $var = value
	if p.curToken.Type == VARIABLE && p.peekToken.Type == ASSIGN {
		stmt.Expression = p.parseAssignmentExpressionFromVariable()
	} else {
		stmt.Expression = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseAssignmentExpressionFromVariable() Expression {
	variable := &Variable{Token: p.curToken, Name: p.curToken.Literal[1:]}

	if !p.expectPeek(ASSIGN) {
		return nil
	}

	assignment := &AssignmentExpression{
		Token: p.curToken,
		Name:  variable,
	}

	p.nextToken()
	assignment.Value = p.parseExpression(LOWEST)

	return assignment
}

func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseVariable() Expression {
	return &Variable{Token: p.curToken, Name: p.curToken.Literal[1:]}
}

func (p *Parser) parseIntegerLiteral() Expression {
	lit := &IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() Expression {
	lit := &FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() Expression {
	literal := p.curToken.Literal

	// Check if string contains variables (simple detection for $var)
	if strings.Contains(literal, "$") {
		return p.parseInterpolatedString()
	}

	return &StringLiteral{Token: p.curToken, Value: literal}
}

func (p *Parser) parseInterpolatedString() Expression {
	literal := p.curToken.Literal
	interpolated := &InterpolatedString{Token: p.curToken}

	// Simple parsing: split on $ and create string parts and variable parts
	parts := strings.Split(literal, "$")

	// First part is always a string (may be empty)
	if parts[0] != "" {
		stringToken := Token{Type: STRING, Literal: parts[0], Line: p.curToken.Line, Column: p.curToken.Column}
		interpolated.Parts = append(interpolated.Parts, &StringLiteral{Token: stringToken, Value: parts[0]})
	}

	// Process variable parts
	for i := 1; i < len(parts); i++ {
		part := parts[i]

		// Extract variable name (up to first non-identifier character)
		varName := ""
		j := 0
		for j < len(part) && (isLetter(part[j]) || (j > 0 && isDigit(part[j]))) {
			j++
		}

		if j > 0 {
			varName = part[:j]
			varToken := Token{Type: VARIABLE, Literal: "$" + varName, Line: p.curToken.Line, Column: p.curToken.Column}
			interpolated.Parts = append(interpolated.Parts, &Variable{Token: varToken, Name: varName})

			// Add remaining string part if any
			if j < len(part) {
				remaining := part[j:]
				stringToken := Token{Type: STRING, Literal: remaining, Line: p.curToken.Line, Column: p.curToken.Column}
				interpolated.Parts = append(interpolated.Parts, &StringLiteral{Token: stringToken, Value: remaining})
			}
		} else {
			// Not a valid variable, treat as string
			stringToken := Token{Type: STRING, Literal: "$" + part, Line: p.curToken.Line, Column: p.curToken.Column}
			interpolated.Parts = append(interpolated.Parts, &StringLiteral{Token: stringToken, Value: "$" + part})
		}
	}

	return interpolated
}

func (p *Parser) parseBooleanLiteral() Expression {
	return &BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(TRUE)}
}

func (p *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseAssignmentExpression(left Expression) Expression {
	variable, ok := left.(*Variable)
	if !ok {
		p.errors = append(p.errors, "left side of assignment must be a variable")
		return nil
	}

	expression := &AssignmentExpression{
		Token: p.curToken,
		Name:  variable,
	}

	p.nextToken()
	expression.Value = p.parseExpression(LOWEST)

	return expression
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseCallExpression(fn Expression) Expression {
	exp := &CallExpression{Token: p.curToken, Function: fn}
	exp.Arguments = p.parseExpressionList(RPAREN)
	return exp
}

func (p *Parser) parseExpressionList(end TokenType) []Expression {
	args := []Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return args
}

func (p *Parser) parseArrayLiteral() Expression {
	tok := p.curToken

	if p.peekTokenIs(RBRACKET) {
		p.nextToken() // consume RBRACKET
		return &ArrayLiteral{Token: tok, Elements: []Expression{}}
	}

	p.nextToken() // move to first element

	// Check if this is an associative array by looking for =>
	firstElement := p.parseExpression(LOWEST)

	if p.peekTokenIs(DOUBLE_ARROW) {
		// This is an associative array
		assocArray := &AssociativeArrayLiteral{Token: tok}

		// Parse first key-value pair
		p.nextToken() // consume =>
		p.nextToken() // move to value
		value := p.parseExpression(LOWEST)
		assocArray.Pairs = append(assocArray.Pairs, ArrayPair{Key: firstElement, Value: value})

		// Parse remaining pairs
		for p.peekTokenIs(COMMA) {
			p.nextToken() // consume comma
			p.nextToken() // move to next key

			key := p.parseExpression(LOWEST)

			if !p.expectPeek(DOUBLE_ARROW) {
				return nil
			}

			p.nextToken() // move to value
			value := p.parseExpression(LOWEST)

			assocArray.Pairs = append(assocArray.Pairs, ArrayPair{Key: key, Value: value})
		}

		if !p.expectPeek(RBRACKET) {
			return nil
		}

		return assocArray
	} else {
		// This is a regular indexed array
		array := &ArrayLiteral{Token: tok}
		array.Elements = []Expression{firstElement}

		// Parse remaining elements
		for p.peekTokenIs(COMMA) {
			p.nextToken() // consume comma
			p.nextToken() // move to next element
			array.Elements = append(array.Elements, p.parseExpression(LOWEST))
		}

		if !p.expectPeek(RBRACKET) {
			return nil
		}

		return array
	}
}

func (p *Parser) parseForStatement() *ForStatement {
	stmt := &ForStatement{Token: p.curToken}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	p.nextToken()
	// Handle assignment in init part of for loop
	if p.curToken.Type == VARIABLE && p.peekToken.Type == ASSIGN {
		stmt.Init = p.parseAssignmentExpressionFromVariable()
	} else {
		stmt.Init = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(SEMICOLON) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(SEMICOLON) {
		return nil
	}

	p.nextToken()
	// Handle assignment or increment in update part of for loop
	if p.curToken.Type == VARIABLE && p.peekToken.Type == ASSIGN {
		stmt.Update = p.parseAssignmentExpressionFromVariable()
	} else if p.curToken.Type == VARIABLE && p.peekToken.Type == INCREMENT {
		// Parse variable first, then parse as postfix
		variable := &Variable{Token: p.curToken, Name: p.curToken.Literal[1:]}
		p.nextToken() // move to INCREMENT
		stmt.Update = &PostfixExpression{
			Token:    p.curToken,
			Left:     variable,
			Operator: p.curToken.Literal,
		}
	} else {
		stmt.Update = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseIndexExpression(left Expression) Expression {
	exp := &IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parsePostfixExpression(left Expression) Expression {
	return &PostfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parseWhileStatement() *WhileStatement {
	stmt := &WhileStatement{Token: p.curToken}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(RPAREN) {
		return nil
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseForeachStatement() *ForeachStatement {
	stmt := &ForeachStatement{Token: p.curToken}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Array = p.parseExpression(LOWEST)

	if !p.expectPeek(AS) {
		return nil
	}

	p.nextToken()

	// Check if we have key => value syntax
	if p.peekTokenIs(DOUBLE_ARROW) {
		// Parse key
		if p.curToken.Type != VARIABLE {
			p.errors = append(p.errors, "foreach key must be a variable")
			return nil
		}
		stmt.Key = &Variable{Token: p.curToken, Name: p.curToken.Literal[1:]}

		p.nextToken() // consume =>
		p.nextToken() // move to value
	}

	// Parse value
	if p.curToken.Type != VARIABLE {
		p.errors = append(p.errors, "foreach value must be a variable")
		return nil
	}
	stmt.Value = &Variable{Token: p.curToken, Name: p.curToken.Literal[1:]}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseBreakStatement() *BreakStatement {
	stmt := &BreakStatement{Token: p.curToken}

	// Check if there's a level specified
	if p.peekTokenIs(INT) {
		p.nextToken()
		stmt.Level = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseContinueStatement() *ContinueStatement {
	stmt := &ContinueStatement{Token: p.curToken}

	// Check if there's a level specified
	if p.peekTokenIs(INT) {
		p.nextToken()
		stmt.Level = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseClassDeclaration() *ClassDeclaration {
	stmt := &ClassDeclaration{Token: p.curToken}

	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Check for inheritance
	if p.peekTokenIs(EXTENDS) {
		p.nextToken() // consume 'extends'
		if !p.expectPeek(IDENT) {
			return nil
		}
		stmt.SuperClass = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	// Check for interface implementations
	if p.peekTokenIs(IMPLEMENTS) {
		p.nextToken() // consume 'implements'
		p.nextToken()
		for !p.curTokenIs(LBRACE) && !p.curTokenIs(EOF) {
			if p.curTokenIs(IDENT) {
				stmt.Interfaces = append(stmt.Interfaces, &Identifier{
					Token: p.curToken,
					Value: p.curToken.Literal,
				})
			}

			if p.peekTokenIs(COMMA) {
				p.nextToken()
			}

			if p.peekTokenIs(LBRACE) {
				break
			}
			p.nextToken()
		}
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	// Parse class body
	p.nextToken()
	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		// Handle trait uses
		if p.curTokenIs(USE) {
			if traitUse := p.parseTraitUse(); traitUse != nil {
				stmt.TraitUses = append(stmt.TraitUses, traitUse)
			}
		} else {
			// Check for visibility modifiers or static
			visibility := "public" // default visibility
			static := false

			if p.curTokenIs(PUBLIC) || p.curTokenIs(PRIVATE) || p.curTokenIs(PROTECTED) {
				visibility = p.curToken.Literal
				p.nextToken()
			}

			if p.curTokenIs(STATIC) {
				static = true
				p.nextToken()
			}

			if p.curTokenIs(CONST) {
				// Class constant
				constant := p.parseConstantDeclaration()
				if constant != nil {
					constant.Visibility = visibility
					stmt.Constants = append(stmt.Constants, constant)
				}
			} else if p.curTokenIs(FUNCTION) {
				// Parse method
				method := p.parseMethodDeclaration(visibility, static)
				if method != nil {
					stmt.Methods = append(stmt.Methods, method)
				}
			} else if p.curTokenIs(VARIABLE) {
				// Parse property
				property := p.parsePropertyDeclaration(visibility, static)
				if property != nil {
					stmt.Properties = append(stmt.Properties, property)
				}
			}
		}

		p.nextToken()
	}

	return stmt
}

func (p *Parser) parsePropertyDeclaration(visibility string, static bool) *PropertyDeclaration {
	if !p.curTokenIs(VARIABLE) {
		return nil
	}

	prop := &PropertyDeclaration{
		Token:      p.curToken,
		Visibility: visibility,
		Static:     static,
		Name:       &Variable{Token: p.curToken, Name: p.curToken.Literal[1:]},
	}

	// Check for default value
	if p.peekTokenIs(ASSIGN) {
		p.nextToken() // consume =
		p.nextToken() // move to value
		prop.Value = p.parseExpression(LOWEST)
	}

	// Expect semicolon
	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return prop
}

func (p *Parser) parseMethodDeclaration(visibility string, static bool) *MethodDeclaration {
	if !p.curTokenIs(FUNCTION) {
		return nil
	}

	method := &MethodDeclaration{
		Token:      p.curToken,
		Visibility: visibility,
		Static:     static,
	}

	if !p.expectPeek(IDENT) {
		return nil
	}

	method.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	method.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(LBRACE) {
		return nil
	}

	method.Body = p.parseBlockStatement()

	return method
}

func (p *Parser) parseNewExpression() Expression {
	expr := &NewExpression{Token: p.curToken}

	if !p.expectPeek(IDENT) {
		return nil
	}

	expr.ClassName = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(LPAREN) {
		p.nextToken() // consume (
		expr.Arguments = p.parseExpressionList(RPAREN)
	}

	return expr
}

func (p *Parser) parseObjectAccessExpression(left Expression) Expression {
	expr := &ObjectAccessExpression{
		Token:  p.curToken,
		Object: left,
	}

	p.nextToken()
	expr.Property = p.parseExpression(CALL)

	return expr
}

func (p *Parser) parseStaticAccessExpression(left Expression) Expression {
	expr := &StaticAccessExpression{
		Token: p.curToken,
		Class: left,
	}

	p.nextToken()
	expr.Property = p.parseExpression(CALL)

	return expr
}

func (p *Parser) parseNamespaceDeclaration() *NamespaceDeclaration {
	stmt := &NamespaceDeclaration{Token: p.curToken}

	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseUseStatement() *UseStatement {
	stmt := &UseStatement{Token: p.curToken}

	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.Namespace = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Check for alias
	if p.peekTokenIs(AS) {
		p.nextToken() // consume 'as'
		if !p.expectPeek(IDENT) {
			return nil
		}
		stmt.Alias = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseTryStatement() *TryStatement {
	stmt := &TryStatement{Token: p.curToken}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	// Parse catch clauses
	for p.peekTokenIs(CATCH) {
		p.nextToken() // consume 'catch'
		catch := p.parseCatchClause()
		if catch != nil {
			stmt.Catches = append(stmt.Catches, catch)
		}
	}

	// Parse optional finally clause
	if p.peekTokenIs(FINALLY) {
		p.nextToken() // consume 'finally'
		if p.expectPeek(LBRACE) {
			stmt.Finally = p.parseBlockStatement()
		}
	}

	return stmt
}

func (p *Parser) parseCatchClause() *CatchClause {
	clause := &CatchClause{Token: p.curToken}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	p.nextToken()

	// Check if there's an exception type
	if p.curToken.Type == IDENT {
		clause.ExceptionType = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken()
	}

	// Parse variable
	if p.curToken.Type != VARIABLE {
		p.errors = append(p.errors, "expected variable in catch clause")
		return nil
	}

	clause.Variable = &Variable{Token: p.curToken, Name: p.curToken.Literal[1:]}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	clause.Body = p.parseBlockStatement()

	return clause
}

func (p *Parser) parseThrowStatement() *ThrowStatement {
	stmt := &ThrowStatement{Token: p.curToken}

	p.nextToken()
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseAnonymousFunction() Expression {
	fn := &AnonymousFunction{Token: p.curToken}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	fn.Parameters = p.parseFunctionParameters()

	// Check for use clause
	if p.peekTokenIs(USE) {
		p.nextToken() // consume 'use'
		if !p.expectPeek(LPAREN) {
			return nil
		}

		p.nextToken()
		for !p.curTokenIs(RPAREN) && !p.curTokenIs(EOF) {
			if p.curToken.Type == VARIABLE {
				fn.UseClause = append(fn.UseClause, &Variable{
					Token: p.curToken,
					Name:  p.curToken.Literal[1:],
				})
			}

			if p.peekTokenIs(COMMA) {
				p.nextToken()
			}
			p.nextToken()
		}
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	fn.Body = p.parseBlockStatement()

	return fn
}

func (p *Parser) parseYieldExpression() Expression {
	expr := &YieldExpression{Token: p.curToken}

	if !p.peekTokenIs(SEMICOLON) && !p.peekTokenIs(RBRACE) && !p.peekTokenIs(EOF) {
		p.nextToken()

		// Parse value or key => value
		value := p.parseExpression(LOWEST)

		if p.peekTokenIs(DOUBLE_ARROW) {
			expr.Key = value
			p.nextToken() // consume =>
			p.nextToken() // move to value
			expr.Value = p.parseExpression(LOWEST)
		} else {
			expr.Value = value
		}
	}

	return expr
}

func (p *Parser) parseInterfaceDeclaration() *InterfaceDeclaration {
	stmt := &InterfaceDeclaration{Token: p.curToken}

	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	p.nextToken()
	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		if method := p.parseInterfaceMethod(); method != nil {
			stmt.Methods = append(stmt.Methods, method)
		}
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseInterfaceMethod() *InterfaceMethod {
	method := &InterfaceMethod{Token: p.curToken}

	// Parse visibility
	if p.curTokenIs(PUBLIC) || p.curTokenIs(PRIVATE) || p.curTokenIs(PROTECTED) {
		method.Visibility = p.curToken.Literal
		p.nextToken()
	} else {
		method.Visibility = "public" // default
	}

	if !p.curTokenIs(FUNCTION) {
		return nil
	}
	p.nextToken()

	if !p.curTokenIs(IDENT) {
		return nil
	}

	method.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	method.Parameters = p.parseFunctionParameters()

	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return method
}

func (p *Parser) parseTraitDeclaration() *TraitDeclaration {
	stmt := &TraitDeclaration{Token: p.curToken}

	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	p.nextToken()
	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		visibility := "public"
		static := false

		// Handle visibility and static modifiers
		if p.curTokenIs(PUBLIC) || p.curTokenIs(PRIVATE) || p.curTokenIs(PROTECTED) {
			visibility = p.curToken.Literal
			p.nextToken()

			if p.curTokenIs(STATIC) {
				static = true
				p.nextToken()
			}
		} else if p.curTokenIs(STATIC) {
			static = true
			p.nextToken()
		}

		if p.curTokenIs(VARIABLE) {
			if property := p.parsePropertyDeclaration(visibility, static); property != nil {
				stmt.Properties = append(stmt.Properties, property)
			}
		} else if p.curTokenIs(FUNCTION) {
			if method := p.parseMethodDeclaration(visibility, static); method != nil {
				stmt.Methods = append(stmt.Methods, method)
			}
		}
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseConstantDeclaration() *ConstantDeclaration {
	stmt := &ConstantDeclaration{Token: p.curToken}

	// Handle visibility for class constants
	if p.curTokenIs(PUBLIC) || p.curTokenIs(PRIVATE) || p.curTokenIs(PROTECTED) {
		stmt.Visibility = p.curToken.Literal
		if !p.expectPeek(CONST) {
			return nil
		}
	} else {
		stmt.Visibility = "public" // default
	}

	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseTraitUse() *TraitUse {
	stmt := &TraitUse{Token: p.curToken}

	p.nextToken()
	for !p.curTokenIs(SEMICOLON) && !p.curTokenIs(EOF) {
		if p.curTokenIs(IDENT) {
			stmt.Traits = append(stmt.Traits, &Identifier{
				Token: p.curToken,
				Value: p.curToken.Literal,
			})
		}

		if p.peekTokenIs(COMMA) {
			p.nextToken()
		}
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseTernaryExpression(condition Expression) Expression {
	expr := &TernaryExpression{
		Token:     p.curToken,
		Condition: condition,
	}

	p.nextToken() // consume '?'
	expr.TrueValue = p.parseExpression(LOWEST)

	if !p.expectPeek(COLON) {
		return nil
	}

	p.nextToken() // consume ':'
	expr.FalseValue = p.parseExpression(LOWEST)

	return expr
}


// Parsefile parses the given PHP file and returns the parsed program
// and any errors encountered during parsing. If the file does not
// exist, it returns an error with a message indicating the file
// does not exist.
func Parsefile(filepath string) (*Program, error) {
	// Check if the file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("File '%s' does not exist", filepath)
	}

	// Read the contents of the file
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("Error reading file '%s': %v", filepath, err)
	}

	// Convert the content to a string
	input := string(content)

	// Create a lexer with the input string
	lexer := New(input)

	// Create a parser with the lexer
	parser := NewParser(lexer)

	// Parse the program
	program := parser.ParseProgram()

	// Check if there are any parser errors
	if len(parser.Errors()) != 0 {
		// If there are errors, construct an error string
		errStr := fmt.Sprintf("Parser errors for file '%s':\n", filepath)
		return nil, fmt.Errorf("%s%s", errStr, strings.Join(parser.Errors(), "\n"))
	}

	// Return the parsed program and nil for the error
	return program, nil
}