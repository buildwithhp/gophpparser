package main

import "fmt"

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	// Literals
	IDENT
	INT
	FLOAT
	STRING
	
	// PHP Tags
	PHP_OPEN  // <?php
	PHP_CLOSE // ?>
	
	// Variables
	VARIABLE // $var
	
	// Operators
	ASSIGN   // =
	PLUS     // +
	MINUS    // -
	MULTIPLY // *
	DIVIDE   // /
	MODULO   // %
	CONCAT   // .
	INCREMENT // ++
	DECREMENT // --
	
	// Comparison
	EQ    // ==
	NOT_EQ // !=
	LT     // <
	GT     // >
	LTE    // <=
	GTE    // >=
	
	// Logical
	AND    // &&
	OR     // ||
	NOT    // !
	
	// Delimiters
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	RBRACKET  // ]
	
	// Keywords
	FUNCTION
	CLASS
	IF
	ELSE
	ELSEIF
	WHILE
	FOR
	FOREACH
	RETURN
	ECHO
	PRINT
	VAR
	PUBLIC
	PRIVATE
	PROTECTED
	STATIC
	CONST
	NEW
	EXTENDS
	IMPLEMENTS
	INTERFACE
	NAMESPACE
	USE
	REQUIRE
	INCLUDE
	TRUE
	FALSE
	NULL
	ARRAY
	BREAK
	CONTINUE
	DO
	AS
	ARROW // =>
	DOUBLE_ARROW // =>
	OBJECT_ACCESS // ->
	STATIC_ACCESS // ::
	NAMESPACE_SEPARATOR // \
	TRY
	CATCH
	FINALLY
	THROW
	EXCEPTION
	CLOSURE // closure/anonymous function
	YIELD
	YIELD_FROM
	// Modern PHP operators
	QUESTION // ?
	QUESTION_QUESTION // ??
	QUESTION_QUESTION_ASSIGN // ??=
	SPACESHIP // <=>
	QUESTION_ARROW // ?->
	// Advanced constructs
	TRAIT
	ABSTRACT
	FINAL
	GLOBAL
	LIST
	UNSET
	ISSET
	EMPTY
	CLONE
	INSTANCEOF
	MATCH
	// Type system
	UNION_TYPE // |
	INTERSECTION_TYPE // &
	// References and variables
	REFERENCE // &
	VARIABLE_VAR // $$
	// Include system
	INCLUDE_ONCE
	REQUIRE_ONCE
	// Magic constants
	MAGIC_CONSTANT
	// Arrow function
	ARROW_FUNCTION // fn
)

type Token struct {
	Type     TokenType
	Literal  string
	Line     int
	Column   int
	Position int
}

var keywords = map[string]TokenType{
	"function":   FUNCTION,
	"class":      CLASS,
	"if":         IF,
	"else":       ELSE,
	"elseif":     ELSEIF,
	"while":      WHILE,
	"for":        FOR,
	"foreach":    FOREACH,
	"return":     RETURN,
	"echo":       ECHO,
	"print":      PRINT,
	"var":        VAR,
	"public":     PUBLIC,
	"private":    PRIVATE,
	"protected":  PROTECTED,
	"static":     STATIC,
	"const":      CONST,
	"new":        NEW,
	"extends":    EXTENDS,
	"implements": IMPLEMENTS,
	"interface":  INTERFACE,
	"namespace":  NAMESPACE,
	"use":        USE,
	"require":    REQUIRE,
	"include":    INCLUDE,
	"true":       TRUE,
	"false":      FALSE,
	"null":       NULL,
	"array":      ARRAY,
	"break":      BREAK,
	"continue":   CONTINUE,
	"do":         DO,
	"as":         AS,
	"try":        TRY,
	"catch":      CATCH,
	"finally":    FINALLY,
	"throw":      THROW,
	"exception":  EXCEPTION,
	"closure":      CLOSURE,
	"yield":        YIELD,
	"trait":        TRAIT,
	"abstract":     ABSTRACT,
	"final":        FINAL,
	"global":       GLOBAL,
	"list":         LIST,
	"unset":        UNSET,
	"isset":        ISSET,
	"empty":        EMPTY,
	"clone":        CLONE,
	"instanceof":   INSTANCEOF,
	"match":        MATCH,
	"include_once": INCLUDE_ONCE,
	"require_once": REQUIRE_ONCE,
	"fn":           ARROW_FUNCTION,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func (t TokenType) String() string {
	switch t {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case IDENT:
		return "IDENT"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case STRING:
		return "STRING"
	case PHP_OPEN:
		return "PHP_OPEN"
	case PHP_CLOSE:
		return "PHP_CLOSE"
	case VARIABLE:
		return "VARIABLE"
	case ASSIGN:
		return "ASSIGN"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case MULTIPLY:
		return "MULTIPLY"
	case DIVIDE:
		return "DIVIDE"
	case MODULO:
		return "MODULO"
	case CONCAT:
		return "CONCAT"
	case INCREMENT:
		return "INCREMENT"
	case DECREMENT:
		return "DECREMENT"
	case EQ:
		return "EQ"
	case NOT_EQ:
		return "NOT_EQ"
	case LT:
		return "LT"
	case GT:
		return "GT"
	case LTE:
		return "LTE"
	case GTE:
		return "GTE"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case NOT:
		return "NOT"
	case COMMA:
		return "COMMA"
	case SEMICOLON:
		return "SEMICOLON"
	case COLON:
		return "COLON"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case LBRACKET:
		return "LBRACKET"
	case RBRACKET:
		return "RBRACKET"
	case BREAK:
		return "BREAK"
	case CONTINUE:
		return "CONTINUE"
	case DO:
		return "DO"
	case AS:
		return "AS"
	case ARROW:
		return "ARROW"
	case DOUBLE_ARROW:
		return "DOUBLE_ARROW"
	case OBJECT_ACCESS:
		return "OBJECT_ACCESS"
	case STATIC_ACCESS:
		return "STATIC_ACCESS"
	case NAMESPACE_SEPARATOR:
		return "NAMESPACE_SEPARATOR"
	case TRY:
		return "TRY"
	case CATCH:
		return "CATCH"
	case FINALLY:
		return "FINALLY"
	case THROW:
		return "THROW"
	case EXCEPTION:
		return "EXCEPTION"
	case CLOSURE:
		return "CLOSURE"
	case YIELD:
		return "YIELD"
	case YIELD_FROM:
		return "YIELD_FROM"
	case QUESTION:
		return "QUESTION"
	case QUESTION_QUESTION:
		return "QUESTION_QUESTION"
	case QUESTION_QUESTION_ASSIGN:
		return "QUESTION_QUESTION_ASSIGN"
	case SPACESHIP:
		return "SPACESHIP"
	case QUESTION_ARROW:
		return "QUESTION_ARROW"
	case TRAIT:
		return "TRAIT"
	case ABSTRACT:
		return "ABSTRACT"
	case FINAL:
		return "FINAL"
	case GLOBAL:
		return "GLOBAL"
	case LIST:
		return "LIST"
	case UNSET:
		return "UNSET"
	case ISSET:
		return "ISSET"
	case EMPTY:
		return "EMPTY"
	case CLONE:
		return "CLONE"
	case INSTANCEOF:
		return "INSTANCEOF"
	case MATCH:
		return "MATCH"
	case UNION_TYPE:
		return "UNION_TYPE"
	case INTERSECTION_TYPE:
		return "INTERSECTION_TYPE"
	case REFERENCE:
		return "REFERENCE"
	case VARIABLE_VAR:
		return "VARIABLE_VAR"
	case INCLUDE_ONCE:
		return "INCLUDE_ONCE"
	case REQUIRE_ONCE:
		return "REQUIRE_ONCE"
	case MAGIC_CONSTANT:
		return "MAGIC_CONSTANT"
	case ARROW_FUNCTION:
		return "ARROW_FUNCTION"
	default:
		return fmt.Sprintf("KEYWORD(%d)", t)
	}
}