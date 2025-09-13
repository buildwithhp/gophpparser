package gophpparser

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
}

func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) peekCharAt(offset int) byte {
	pos := l.readPosition + offset
	if pos >= len(l.input) {
		return 0
	}
	return l.input[pos]
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: DOUBLE_ARROW, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(ASSIGN, l.ch, l.line, l.column)
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: INCREMENT, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(PLUS, l.ch, l.line, l.column)
		}
	case '-':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: OBJECT_ACCESS, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(MINUS, l.ch, l.line, l.column)
		}
	case '*':
		tok = newToken(MULTIPLY, l.ch, l.line, l.column)
	case '/':
		if l.peekChar() == '/' {
			l.skipComment()
			return l.NextToken()
		} else {
			tok = newToken(DIVIDE, l.ch, l.line, l.column)
		}
	case '%':
		tok = newToken(MODULO, l.ch, l.line, l.column)
	case '.':
		tok = newToken(CONCAT, l.ch, l.line, l.column)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: NOT_EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(NOT, l.ch, l.line, l.column)
		}
	case '<':
		if l.peekChar() == '=' && l.peekCharAt(1) == '>' {
			ch := l.ch
			l.readChar()
			l.readChar()
			tok = Token{Type: SPACESHIP, Literal: string(ch) + "=>", Line: l.line, Column: l.column}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: LTE, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else if l.peekChar() == '?' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == 'p' {
				l.readChar()
				if l.peekChar() == 'h' {
					l.readChar()
					if l.peekChar() == 'p' {
						l.readChar()
						tok = Token{Type: PHP_OPEN, Literal: "<?php", Line: l.line, Column: l.column}
					} else {
						tok = newToken(LT, ch, l.line, l.column)
					}
				} else {
					tok = newToken(LT, ch, l.line, l.column)
				}
			} else {
				tok = newToken(LT, ch, l.line, l.column)
			}
		} else {
			tok = newToken(LT, l.ch, l.line, l.column)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: GTE, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(GT, l.ch, l.line, l.column)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: AND, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(ILLEGAL, l.ch, l.line, l.column)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: OR, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(ILLEGAL, l.ch, l.line, l.column)
		}
	case '?':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: PHP_CLOSE, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else if l.peekChar() == '?' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = Token{Type: QUESTION_QUESTION_ASSIGN, Literal: string(ch) + "?=", Line: l.line, Column: l.column}
			} else {
				tok = Token{Type: QUESTION_QUESTION, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
			}
		} else if l.peekChar() == '-' && l.peekCharAt(1) == '>' {
			ch := l.ch
			l.readChar()
			l.readChar()
			tok = Token{Type: QUESTION_ARROW, Literal: string(ch) + "->", Line: l.line, Column: l.column}
		} else {
			tok = newToken(QUESTION, l.ch, l.line, l.column)
		}
	case ',':
		tok = newToken(COMMA, l.ch, l.line, l.column)
	case ';':
		tok = newToken(SEMICOLON, l.ch, l.line, l.column)
	case '(':
		tok = newToken(LPAREN, l.ch, l.line, l.column)
	case ')':
		tok = newToken(RPAREN, l.ch, l.line, l.column)
	case '{':
		tok = newToken(LBRACE, l.ch, l.line, l.column)
	case '}':
		tok = newToken(RBRACE, l.ch, l.line, l.column)
	case '[':
		tok = newToken(LBRACKET, l.ch, l.line, l.column)
	case ']':
		tok = newToken(RBRACKET, l.ch, l.line, l.column)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString('"')
		tok.Line = l.line
		tok.Column = l.column
	case '\'':
		tok.Type = STRING
		tok.Literal = l.readString('\'')
		tok.Line = l.line
		tok.Column = l.column
	case '$':
		tok.Type = VARIABLE
		l.readChar()
		tok.Literal = "$" + l.readIdentifier()
		tok.Line = l.line
		tok.Column = l.column
		return tok
	case ':':
		if l.peekChar() == ':' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: STATIC_ACCESS, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = newToken(COLON, l.ch, l.line, l.column)
		}
	case '\\':
		tok = newToken(NAMESPACE_SEPARATOR, l.ch, l.line, l.column)
	case 0:
		tok.Literal = ""
		tok.Type = EOF
		tok.Line = l.line
		tok.Column = l.column
	default:
		if isLetter(l.ch) {
			tok.Line = l.line
			tok.Column = l.column
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type, tok.Literal = l.readNumber()
			tok.Line = l.line
			tok.Column = l.column
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch, l.line, l.column)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() (TokenType, string) {
	position := l.position
	tokenType := INT

	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' && isDigit(l.peekChar()) {
		tokenType = FLOAT
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return tokenType, l.input[position:l.position]
}

func (l *Lexer) readString(delimiter byte) string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == delimiter || l.ch == 0 {
			break
		}
		if l.ch == '\\' {
			l.readChar()
		}
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch > 127
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType TokenType, ch byte, line, column int) Token {
	return Token{Type: tokenType, Literal: string(ch), Line: line, Column: column}
}

func (l *Lexer) skipComment() {
	if l.ch == '/' && l.peekChar() == '/' {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
	} else if l.ch == '/' && l.peekChar() == '*' {
		l.readChar()
		l.readChar()
		for {
			if l.ch == 0 {
				break
			}
			if l.ch == '*' && l.peekChar() == '/' {
				l.readChar()
				l.readChar()
				break
			}
			l.readChar()
		}
	} else if l.ch == '#' {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
	}
}
