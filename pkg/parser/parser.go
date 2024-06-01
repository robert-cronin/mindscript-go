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

package parser

import (
	"fmt"

	"github.com/robert-cronin/mindscript-go/pkg/lexer"
)

type Parser struct {
	l *lexer.Lexer

	curToken  lexer.Token
	peekToken lexer.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// Read two tokens, so curToken and peekToken are both set
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

	for p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() Statement {
	fmt.Println("parseStatement: ", p.curToken.Type)
	switch p.curToken.Type {
	case lexer.AGENT:
		return p.parseAgentStatement()
	default:
		return nil
	}
}

func (p *Parser) parseAgentStatement() *Agent {
	stmt := &Agent{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	// Parse goal, capabilities, behaviors, and functions here

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		p.nextToken()

		switch p.curToken.Type {
		case lexer.GOAL:
			stmt.Goal = p.parseGoal()
		case lexer.CAPABILITIES:
			stmt.Capabilities = p.parseCapabilities()
		case lexer.BEHAVIOR:
			stmt.Behaviors = append(stmt.Behaviors, p.parseBehavior())
		case lexer.FUNCTION:
			stmt.Functions = append(stmt.Functions, p.parseFunction())
		}
	}

	return stmt
}

func (p *Parser) parseGoal() *Goal {
	goal := &Goal{Token: p.curToken}

	if !p.expectPeek(lexer.STRING) {
		return nil
	}

	goal.Value = p.curToken.Literal

	return goal
}

func (p *Parser) parseCapabilities() *Capabilities {
	capabilities := &Capabilities{Token: p.curToken}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		p.nextToken()
		if p.curToken.Type == lexer.STRING {
			capabilities.Values = append(capabilities.Values, p.curToken.Literal)
		}
	}

	return capabilities
}

func (p *Parser) parseBehavior() *Behavior {
	behavior := &Behavior{Token: p.curToken}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	// Parse events and actions inside behavior block

	return behavior
}

func (p *Parser) parseFunction() *Function {
	function := &Function{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	function.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	function.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	function.ReturnType = p.curToken.Literal

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	function.Body = p.parseBlockStatement()

	return function
}

func (p *Parser) parseFunctionParameters() []*Identifier {
	identifiers := []*Identifier{}

	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Token: p.curToken}
	block.Statements = []Statement{}

	p.nextToken()

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		return false
	}
}
