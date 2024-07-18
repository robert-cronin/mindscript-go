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
	"errors"
	"fmt"
	"strconv"

	"github.com/robert-cronin/mindscript-go/pkg/lexer"
	"github.com/robert-cronin/mindscript-go/pkg/logger"
	"go.uber.org/zap"
)

type Parser struct {
	l *lexer.Lexer

	curToken  lexer.Token
	peekToken lexer.Token

	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,

		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekError(expectedType lexer.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s instead",
		expectedType, p.peekToken.Type)
	p.addError(msg)
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
	switch p.curToken.Type {
	case lexer.AGENT:
		// TODO: make err handling like this everywhere else
		agent, err := p.parseAgentStatement()
		if err != nil {
			logger.Log.Error("Error parsing agent statement", zap.Error(err))
			return nil
		}
		return agent
	case lexer.VAR:
		return p.parseVarStatement()
	case lexer.IDENT:
		return p.parseExpressionStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.FUNCTION:
		return p.parseFunction()
	default:
		msg := fmt.Sprintf("Unexpected token %s encountered", p.curToken.Type)
		p.addError(msg)
		return nil
	}
}

func (p *Parser) parseAgentStatement() (*AgentStatement, error) {
	stmt := &AgentStatement{}
	stmt.Token = p.curToken

	if !p.expectPeek(lexer.IDENT) {
		err := errors.New("Agent statement must have a name")
		return nil, err
	}

	stmt.Name = &Identifier{}
	stmt.Name.Token = p.curToken
	stmt.Name.Value = p.curToken.Literal

	if !p.expectPeek(lexer.LBRACE) {
		err := errors.New("Agent statement must have a body")
		return nil, err
	}

Loop:
	for !p.curTokenIs(lexer.EOF) {
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
		case lexer.RBRACE:
			break Loop
		}
	}

	return stmt, nil
}

func (p *Parser) parseGoal() *Goal {
	goal := &Goal{}
	goal.Token = p.curToken

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	if !p.expectPeek(lexer.STRING) {
		return nil
	}

	goal.Value = p.curToken.Literal

	return goal
}

func (p *Parser) parseCapabilities() *Capabilities {
	capabilities := &Capabilities{}
	capabilities.Token = p.curToken

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACKET) {
		return nil
	}

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		p.nextToken()
		if p.curToken.Type == lexer.STRING {
			capabilities.Values = append(capabilities.Values, p.curToken.Literal)
		} else if p.curToken.Type == lexer.COMMA {
			continue
		} else if p.curToken.Type == lexer.RBRACKET {
			break
		} else {
			logger.Log.Error("Error parsing capabilities")
			return nil
		}
	}

	return capabilities
}

func (p *Parser) parseBehavior() *Behavior {
	behavior := &Behavior{}
	behavior.Token = p.curToken
	behavior.EventHandlers = []*EventHandler{}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

Loop:
	for !p.curTokenIs(lexer.EOF) {
		p.nextToken()

		switch p.curToken.Type {
		case lexer.ON:
			behavior.EventHandlers = append(behavior.EventHandlers, p.parseEventHandler())
		case lexer.RBRACE:
			break Loop
		default:
			logger.Log.Error("Error parsing behavior")
			return nil
		}
	}

	return behavior
}

func (p *Parser) parseEventHandler() *EventHandler {
	eventHandler := &EventHandler{}
	eventHandler.Token = p.curToken

	if !p.expectPeek(lexer.STRING) {
		return nil
	}

	eventHandler.Event = &Event{}
	eventHandler.Event.Name = &Identifier{}
	eventHandler.Event.Name.Token = p.curToken
	eventHandler.Event.Name.Value = p.curToken.Literal

	if !p.expectPeek(lexer.LBRACE) {
		logger.Log.Error("Error parsing event handler")
		return nil
	}

	eventHandler.BlockStatement = p.parseBlockStatement()

	return eventHandler
}

func (p *Parser) parseFunction() *Function {
	function := &Function{}
	function.Token = p.curToken

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	function.Name = &Identifier{}
	function.Name.Token = p.curToken
	function.Name.Value = p.curToken.Literal

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	function.Arguments = p.parseFunctionArguments()

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	function.ReturnType = p.parseReturnDataType()

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	function.Body = p.parseBlockStatement()

	return function
}

func (p *Parser) parseVarStatement() *VarStatement {
	stmt := &VarStatement{}
	stmt.Token = p.curToken

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = &Identifier{}
	stmt.Name.Token = p.curToken
	stmt.Name.Value = p.curToken.Literal

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	stmt.Type = p.parseDataType()
	if stmt.Type == nil {
		return nil
	}

	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseDataType() *DataType {
	dataType := &DataType{}

	switch p.peekToken.Type {
	case lexer.INT, lexer.FLOAT, lexer.STRING, lexer.BOOL:
		p.nextToken()
		dataType.Token = p.curToken
	default:
		logger.Log.Error("Error parsing data type")
		return nil
	}

	return dataType
}

func (p *Parser) parseFunctionArguments() []*FunctionArgument {
	args := []*FunctionArgument{}

	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()

	arg := &FunctionArgument{}
	arg.Name = &Identifier{}
	arg.Name.Token = p.curToken
	arg.Name.Value = p.curToken.Literal

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	arg.Type = p.parseDataType()

	args = append(args, arg)

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()

		arg := &FunctionArgument{}
		arg.Name = &Identifier{}
		arg.Name.Token = p.curToken
		arg.Name.Value = p.curToken.Literal

		if !p.expectPeek(lexer.COLON) {
			return nil
		}

		arg.Type = p.parseDataType()

		args = append(args, arg)
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{}
	block.Token = p.curToken
	block.Statements = make(map[int]*Statement, 0)

	count := 0

	p.nextToken()

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements[count] = &stmt
			count++
		}
		p.nextToken()
	}

	return block
}

const (
	_ int = iota
	LOWEST
	SUM     // + or -
	PRODUCT // * or /
	PREFIX  // -X or !X
	CALL    // myFunction(X)
)

var precedences = map[lexer.TokenType]int{
	lexer.PLUS:     SUM,
	lexer.MINUS:    SUM,
	lexer.ASTERISK: PRODUCT,
	lexer.SLASH:    PRODUCT,
	lexer.LPAREN:   CALL,
}

func (p *Parser) parseExpression(precedence int) *Expression {
	var leftExp Expression

	switch p.curToken.Type {
	case lexer.IDENT:
		leftExp = p.parseIdentifier()
	case lexer.INT:
		leftExp = p.parseIntegerLiteral()
	case lexer.FLOAT:
		leftExp = p.parseFloatLiteral()
	case lexer.STRING:
		leftExp = p.parseStringLiteral()
	case lexer.BOOL:
		leftExp = p.parseBooleanLiteral()
	default:
		// Check first if its a function call
		if p.peekToken.Type != lexer.LPAREN {
			return nil
		}
	}

	for !p.peekTokenIs(lexer.SEMICOLON) && precedence < p.peekPrecedence() {
		switch p.peekToken.Type {
		case lexer.PLUS, lexer.MINUS, lexer.ASTERISK, lexer.SLASH:
			p.nextToken()
			leftExp = p.parseInfixExpression(leftExp)
		case lexer.LPAREN:
			p.nextToken()
			leftExp = p.parseCallExpression(leftExp)
		default:
			return &leftExp
		}
	}

	return &leftExp
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		BaseNode: BaseNode{Token: p.curToken},
		Left:     &left,
		Operator: &p.curToken,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseCallExpression(function Expression) Expression {
	exp := &CallExpression{BaseNode: BaseNode{Token: p.curToken}, Function: &function}
	exp.Arguments = p.parseExpressionList(lexer.RPAREN)
	return exp
}

func (p *Parser) parseExpressionList(end lexer.TokenType) []*Expression {
	list := []*Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseIdentifier() *IdentifierLiteral {
	ident := &IdentifierLiteral{}
	ident.Token = p.curToken
	ident.Value = p.curToken.Literal
	return ident
}

func (p *Parser) parseIntegerLiteral() *IntegerLiteral {
	integer := &IntegerLiteral{}
	integer.Token = p.curToken

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		logger.Log.Error("Error parsing integer literal", zap.Error(err))
		return nil
	}

	integer.Value = value
	return integer
}

func (p *Parser) parseFloatLiteral() *FloatLiteral {
	float := &FloatLiteral{}
	float.Token = p.curToken

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		logger.Log.Error("Error parsing float literal", zap.Error(err))
		return nil
	}

	float.Value = value

	return float
}

func (p *Parser) parseStringLiteral() *StringLiteral {
	str := &StringLiteral{}
	str.Token = p.curToken

	str.Value = p.curToken.Literal

	return str
}

func (p *Parser) parseBooleanLiteral() *BooleanLiteral {
	boolean := &BooleanLiteral{}
	boolean.Token = p.curToken

	value, err := strconv.ParseBool(p.curToken.Literal)
	if err != nil {
		logger.Log.Error("Error parsing boolean literal", zap.Error(err))
		return nil
	}

	boolean.Value = value

	return boolean
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{}
	stmt.Token = p.curToken

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnDataType() *DataType {
	dataType := &DataType{}

	switch p.peekToken.Type {
	case lexer.INT, lexer.FLOAT, lexer.STRING, lexer.BOOL:
		p.nextToken()
		dataType.Token = p.curToken
	default:
		logger.Log.Error("Error parsing return data type")
		return nil
	}

	return dataType
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{}
	stmt.Token = p.curToken

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
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
	}
	p.peekError(t)
	return false
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
