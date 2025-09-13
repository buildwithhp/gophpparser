package gophpparser

import (
	"fmt"
	"strings"
)

// TokenInfo provides detailed information about a token for debugging
type TokenInfo struct {
	Type     TokenType `json:"type"`
	TypeName string    `json:"type_name"`
	Literal  string    `json:"literal"`
	Line     int       `json:"line"`
	Column   int       `json:"column"`
}

// DebugParseErrors helps identify why parsing is failing
type DebugParseErrors struct {
	Input              string      `json:"input"`
	Tokens             []TokenInfo `json:"tokens"`
	ParsingErrors      []string    `json:"parsing_errors"`
	UnknownTokens      []TokenInfo `json:"unknown_tokens"`
	MissingPrefixFuncs []string    `json:"missing_prefix_functions"`
}

// DebugParsePHP provides detailed debugging information for failed parsing
func DebugParsePHP(input string) *DebugParseErrors {
	debug := &DebugParseErrors{
		Input:              input,
		Tokens:             []TokenInfo{},
		ParsingErrors:      []string{},
		UnknownTokens:      []TokenInfo{},
		MissingPrefixFuncs: []string{},
	}

	// Tokenize the input
	lexer := New(input)
	
	// Collect all tokens
	for {
		token := lexer.NextToken()
		
		tokenInfo := TokenInfo{
			Type:     token.Type,
			TypeName: getTokenTypeName(token.Type),
			Literal:  token.Literal,
			Line:     token.Line,
			Column:   token.Column,
		}
		
		debug.Tokens = append(debug.Tokens, tokenInfo)
		
		// Check for unknown or problematic tokens
		if token.Type == ILLEGAL {
			debug.UnknownTokens = append(debug.UnknownTokens, tokenInfo)
		}
		
		if token.Type == EOF {
			break
		}
	}

	// Try parsing and collect errors
	lexer = New(input)
	parser := NewParser(lexer)
	_ = parser.ParseProgram()
	
	debug.ParsingErrors = parser.Errors()
	
	// Analyze missing prefix functions
	debug.analyzeMissingPrefixFunctions()
	
	return debug
}

// getTokenTypeName returns a human-readable name for a token type
func getTokenTypeName(tokenType TokenType) string {
	switch tokenType {
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
	case FUNCTION:
		return "FUNCTION"
	case CLASS:
		return "CLASS"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case ELSEIF:
		return "ELSEIF"
	case WHILE:
		return "WHILE"
	case FOR:
		return "FOR"
	case FOREACH:
		return "FOREACH"
	case RETURN:
		return "RETURN"
	case ECHO:
		return "ECHO"
	case PRINT:
		return "PRINT"
	case VAR:
		return "VAR"
	case PUBLIC:
		return "PUBLIC"
	case PRIVATE:
		return "PRIVATE"
	case PROTECTED:
		return "PROTECTED"
	case STATIC:
		return "STATIC"
	case CONST:
		return "CONST"
	case NEW:
		return "NEW"
	case EXTENDS:
		return "EXTENDS"
	case IMPLEMENTS:
		return "IMPLEMENTS"
	case INTERFACE:
		return "INTERFACE"
	case NAMESPACE:
		return "NAMESPACE"
	case USE:
		return "USE"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case NULL:
		return "NULL"
	case ARRAY:
		return "ARRAY"
	case BREAK:
		return "BREAK"
	case CONTINUE:
		return "CONTINUE"
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
	case YIELD:
		return "YIELD"
	case QUESTION:
		return "QUESTION"
	case QUESTION_QUESTION:
		return "QUESTION_QUESTION"
	case QUESTION_ARROW:
		return "QUESTION_ARROW"
	case SPACESHIP:
		return "SPACESHIP"
	case TRAIT:
		return "TRAIT"
	case ABSTRACT:
		return "ABSTRACT"
	case FINAL:
		return "FINAL"
	case GLOBAL:
		return "GLOBAL"
	case CLONE:
		return "CLONE"
	case INSTANCEOF:
		return "INSTANCEOF"
	case MAGIC_CONSTANT:
		return "MAGIC_CONSTANT"
	case COMMENT:
		return "COMMENT"
	case DOCBLOCK:
		return "DOCBLOCK"
	default:
		return fmt.Sprintf("UNKNOWN_TOKEN(%d)", int(tokenType))
	}
}

// analyzeMissingPrefixFunctions identifies which prefix functions are missing
func (d *DebugParseErrors) analyzeMissingPrefixFunctions() {
	missingPrefixes := make(map[string]bool)
	
	for _, errMsg := range d.ParsingErrors {
		if strings.Contains(errMsg, "no prefix parse function for") {
			// Extract the token type from error message
			start := strings.Index(errMsg, "no prefix parse function for ") + len("no prefix parse function for ")
			end := strings.Index(errMsg[start:], " found")
			if end != -1 {
				tokenType := errMsg[start : start+end]
				missingPrefixes[tokenType] = true
			}
		}
	}
	
	for prefix := range missingPrefixes {
		d.MissingPrefixFuncs = append(d.MissingPrefixFuncs, prefix)
	}
}

// PrintDebugInfo prints a human-readable debug report
func (d *DebugParseErrors) PrintDebugInfo() {
	fmt.Println("=== PHP Parser Debug Report ===")
	fmt.Printf("Input length: %d characters\n", len(d.Input))
	fmt.Printf("Total tokens: %d\n", len(d.Tokens))
	fmt.Printf("Parsing errors: %d\n", len(d.ParsingErrors))
	fmt.Printf("Unknown tokens: %d\n", len(d.UnknownTokens))
	
	if len(d.UnknownTokens) > 0 {
		fmt.Println("\n--- Unknown/Illegal Tokens ---")
		for _, token := range d.UnknownTokens {
			fmt.Printf("  Line %d:%d - %s (%s)\n", token.Line, token.Column, token.Literal, token.TypeName)
		}
	}
	
	if len(d.MissingPrefixFuncs) > 0 {
		fmt.Println("\n--- Missing Prefix Parse Functions ---")
		for _, missing := range d.MissingPrefixFuncs {
			fmt.Printf("  - %s\n", missing)
		}
	}
	
	if len(d.ParsingErrors) > 0 && len(d.ParsingErrors) <= 10 {
		fmt.Println("\n--- Parsing Errors ---")
		for i, err := range d.ParsingErrors {
			fmt.Printf("  %d. %s\n", i+1, err)
		}
	} else if len(d.ParsingErrors) > 10 {
		fmt.Println("\n--- First 10 Parsing Errors ---")
		for i := 0; i < 10; i++ {
			fmt.Printf("  %d. %s\n", i+1, d.ParsingErrors[i])
		}
		fmt.Printf("  ... and %d more errors\n", len(d.ParsingErrors)-10)
	}
	
	fmt.Println("\n--- Token Stream (first 20 tokens) ---")
	limit := len(d.Tokens)
	if limit > 20 {
		limit = 20
	}
	for i := 0; i < limit; i++ {
		token := d.Tokens[i]
		fmt.Printf("  %2d. %-20s %q\n", i+1, token.TypeName, token.Literal)
	}
	if len(d.Tokens) > 20 {
		fmt.Printf("  ... and %d more tokens\n", len(d.Tokens)-20)
	}
}

// GetMostCommonErrors returns the most frequently occurring error types
func (d *DebugParseErrors) GetMostCommonErrors() map[string]int {
	errorCounts := make(map[string]int)
	
	for _, err := range d.ParsingErrors {
		if strings.Contains(err, "no prefix parse function") {
			errorCounts["missing_prefix_function"]++
		} else if strings.Contains(err, "expected next token") {
			errorCounts["unexpected_token"]++
		} else {
			errorCounts["other"]++
		}
	}
	
	return errorCounts
}

// SuggestFixes provides suggestions for fixing the most common issues
func (d *DebugParseErrors) SuggestFixes() []string {
	var suggestions []string
	
	errorCounts := d.GetMostCommonErrors()
	
	if errorCounts["missing_prefix_function"] > 0 {
		suggestions = append(suggestions, 
			"Add missing prefix parse functions for tokens: "+strings.Join(d.MissingPrefixFuncs, ", "))
	}
	
	if errorCounts["unexpected_token"] > 0 {
		suggestions = append(suggestions, 
			"The PHP code contains syntax that the parser doesn't support yet")
	}
	
	if len(d.UnknownTokens) > 0 {
		suggestions = append(suggestions, 
			"The lexer encountered tokens it doesn't recognize - check for unsupported PHP syntax")
	}
	
	return suggestions
}