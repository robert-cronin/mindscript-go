/**
 * Copyright 2024 Robert Cronin
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package lexer

import (
	"unicode"
)

type TokenType string

const (
	IDENT     TokenType = "IDENT"
	KEYWORD   TokenType = "KEYWORD"
	STRING    TokenType = "STRING"
	INT       TokenType = "INT"
	FLOAT     TokenType = "FLOAT"
	BOOL      TokenType = "BOOL"
	LBRACE    TokenType = "LBRACE"
	RBRACE    TokenType = "RBRACE"
	LPAREN    TokenType = "LPAREN"
	RPAREN    TokenType = "RPAREN"
	SEMICOLON TokenType = "SEMICOLON"
	PLUS      TokenType = "PLUS"
	MINUS     TokenType = "MINUS"
	ASTERISK  TokenType = "ASTERISK"
	SLASH     TokenType = "SLASH"
	ASSIGN    TokenType = "ASSIGN"
	GT        TokenType = "GT"
	LT        TokenType = "LT"
	EQ        TokenType = "EQ"
	AND       TokenType = "AND"
	OR        TokenType = "OR"
)

type Token struct {
	Type    TokenType
	Literal string
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
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
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()
	switch l.ch {
	case '{':
		tok = Token{Type: LBRACE, Literal: string(l.ch)}
	case '}':
		tok = Token{Type: RBRACE, Literal: string(l.ch)}
	case '(':
		tok = Token{Type: LPAREN, Literal: string(l.ch)}
	case ')':
		tok = Token{Type: RPAREN, Literal: string(l.ch)}
	case ';':
		tok = Token{Type: SEMICOLON, Literal: string(l.ch)}
	case '+':
		tok = Token{Type: PLUS, Literal: string(l.ch)}
	case '-':
		tok = Token{Type: MINUS, Literal: string(l.ch)}
	case '*':
		tok = Token{Type: ASTERISK, Literal: string(l.ch)}
	case '/':
		tok = Token{Type: SLASH, Literal: string(l.ch)}
	case '=':
		tok = Token{Type: ASSIGN, Literal: string(l.ch)}
	case '>':
		tok = Token{Type: GT, Literal: string(l.ch)}
	case '<':
		tok = Token{Type: LT, Literal: string(l.ch)}
	case '&':
		tok = Token{Type: AND, Literal: string(l.ch)}
	case '|':
		tok = Token{Type: OR, Literal: string(l.ch)}
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
	default:
		if isDigit(l.ch) {
			if l.peekChar() == '.' {
				tok.Literal = l.readFloat()
				tok.Type = FLOAT
			} else {
				tok.Literal = l.readInt()
				tok.Type = INT
			}
			return tok
		} else if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = KEYWORD
			return tok
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readInt() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readFloat() string {
	position := l.position
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch))
}

func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}
