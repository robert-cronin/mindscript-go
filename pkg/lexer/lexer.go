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

// the lexer takes the input and converts it into tokens
package lexer

import (
	"strings"
	"unicode"

	"go.uber.org/zap"
)

type TokenType string

const (
	IDENT     TokenType = "IDENT"
	LBRACE    TokenType = "LBRACE"
	RBRACE    TokenType = "RBRACE"
	LPAREN    TokenType = "LPAREN"
	RPAREN    TokenType = "RPAREN"
	LBRACKET  TokenType = "LBRACKET"
	RBRACKET  TokenType = "RBRACKET"
	COLON     TokenType = "COLON"
	SEMICOLON TokenType = "SEMICOLON"
	COMMA     TokenType = "COMMA"
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
	AGENT     TokenType = "AGENT"
	ON        TokenType = "ON"
	VAR       TokenType = "VAR"
	RETURN    TokenType = "RETURN"

	GOAL         TokenType = "GOAL"
	CAPABILITIES TokenType = "CAPABILITIES"
	BEHAVIOR     TokenType = "BEHAVIOR"
	FUNCTION     TokenType = "FUNCTION"
	EOF          TokenType = "EOF"
)

// Data types
const (
	STRING TokenType = "STRING"
	INT    TokenType = "INT"
	FLOAT  TokenType = "FLOAT"
	BOOL   TokenType = "BOOL"
)

// Store a list of keywords
var keywords = map[string]TokenType{
	"agent":        AGENT,
	"goal":         GOAL,
	"capabilities": CAPABILITIES,
	"behavior":     BEHAVIOR,
	"function":     FUNCTION,
	"on":           ON,
	"var":          VAR,
	"int":          INT,
	"float":        FLOAT,
	"string":       STRING,
	"bool":         BOOL,
	"return":       RETURN,
}

type Token struct {
	Type    TokenType
	Literal string
	Loc     int
}

type Lexer struct {
	logger       *zap.Logger
	input        string
	position     int
	readPosition int
	ch           byte
}

// Line gets the line number of the provided token
func (l *Lexer) Line(tok Token) int {
	return 1 + strings.Count(l.Prefix(tok.Loc), "\n")
}

// Column gets the column number of the provided token
func (l *Lexer) Column(tok Token) int {
	return 1 + strings.LastIndex(l.Prefix(tok.Loc), "\n")
}

func New(input string) *Lexer {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to initialize Zap logger: " + err.Error())
	}
	l := &Lexer{logger: logger, input: input}
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
	tok.Loc = l.position
	switch l.ch {
	case '{':
		tok = Token{Type: LBRACE, Literal: string(l.ch), Loc: l.position}
	case '}':
		tok = Token{Type: RBRACE, Literal: string(l.ch), Loc: l.position}
	case '(':
		tok = Token{Type: LPAREN, Literal: string(l.ch), Loc: l.position}
	case ')':
		tok = Token{Type: RPAREN, Literal: string(l.ch), Loc: l.position}
	case '[':
		tok = Token{Type: LBRACKET, Literal: string(l.ch), Loc: l.position}
	case ']':
		tok = Token{Type: RBRACKET, Literal: string(l.ch), Loc: l.position}
	case ':':
		tok = Token{Type: COLON, Literal: string(l.ch), Loc: l.position}
	case ';':
		tok = Token{Type: SEMICOLON, Literal: string(l.ch), Loc: l.position}
	case ',':
		tok = Token{Type: COMMA, Literal: string(l.ch), Loc: l.position}
	case '+':
		tok = Token{Type: PLUS, Literal: string(l.ch), Loc: l.position}
	case '-':
		tok = Token{Type: MINUS, Literal: string(l.ch), Loc: l.position}
	case '*':
		tok = Token{Type: ASTERISK, Literal: string(l.ch), Loc: l.position}
	case '/':
		tok = Token{Type: SLASH, Literal: string(l.ch), Loc: l.position}
	case '=':
		tok = Token{Type: ASSIGN, Literal: string(l.ch), Loc: l.position}
	case '>':
		tok = Token{Type: GT, Literal: string(l.ch), Loc: l.position}
	case '<':
		tok = Token{Type: LT, Literal: string(l.ch), Loc: l.position}
	case '&':
		tok = Token{Type: AND, Literal: string(l.ch), Loc: l.position}
	case '|':
		tok = Token{Type: OR, Literal: string(l.ch), Loc: l.position}
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
		tok.Loc = l.position
	case 0:
		tok.Type = EOF
		tok.Literal = "EOF"
		tok.Loc = l.position
	default:
		if isDigit(l.ch) {
			if l.peekChar() == '.' {
				tok.Literal = l.readFloat()
				tok.Type = FLOAT
				tok.Loc = l.position
			} else {
				tok.Literal = l.readInt()
				tok.Type = INT
				tok.Loc = l.position
			}
			return tok
		} else if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = IDENT
			tok.Loc = l.position
			if keywordType, ok := keywords[tok.Literal]; ok {
				tok.Type = keywordType
			}
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
	return l.peekCharOffset(0)
}

func (l *Lexer) peekCharOffset(offset int) byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition+offset]
	}
}

// Helper to get prefix up to loc
func (l *Lexer) Prefix(loc int) string {
	return l.input[:loc]
}
