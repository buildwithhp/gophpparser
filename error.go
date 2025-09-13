package gophpparser

import (
	"fmt"
)

type ParseError struct {
	Message string
	Line    int
	Column  int
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("Parse error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}

type ErrorHandler struct {
	errors []ParseError
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		errors: []ParseError{},
	}
}

func (e *ErrorHandler) AddError(message string, line, column int) {
	e.errors = append(e.errors, ParseError{
		Message: message,
		Line:    line,
		Column:  column,
	})
}

func (e *ErrorHandler) HasErrors() bool {
	return len(e.errors) > 0
}

func (e *ErrorHandler) GetErrors() []ParseError {
	return e.errors
}

func (e *ErrorHandler) Clear() {
	e.errors = []ParseError{}
}

func (e *ErrorHandler) PrintErrors() {
	if len(e.errors) == 0 {
		return
	}

	fmt.Printf("Found %d error(s):\n", len(e.errors))
	for _, err := range e.errors {
		fmt.Printf("  - %s\n", err.Error())
	}
}
