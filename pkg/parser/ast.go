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

// the parser package is responsible for parsing the tokens from the lexer
package parser

import "github.com/robert-cronin/mindscript-go/pkg/lexer"

// Node interface for all nodes
type Node interface {
	TokenLiteral() string
}

// BaseNode provides common fields and methods for nodes
type BaseNode struct {
	Node
	Token lexer.Token `json:"token"`
}

func (b *BaseNode) TokenLiteral() string {
	return b.Token.Literal
}

// Program represents the entire program
type Program struct {
	Statements []Statement `json:"statements"`
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// Statement interface for all statements
type Statement interface {
	Node
	statementNode()
}

// Expression interface for all expressions
type Expression interface {
	Node
	expressionNode()
}

// Identifier represents a variable or function name
type Identifier struct {
	BaseNode
	Value string `json:"value"`
}

func (i *Identifier) expressionNode() {}

// AgentStatement represents an agent declaration
type AgentStatement struct {
	BaseNode
	Name         *Identifier   `json:"name"`
	Goal         *Goal         `json:"goal"`
	Capabilities *Capabilities `json:"capabilities"`
	Behaviors    []*Behavior   `json:"behaviors"`
	Functions    []*Function   `json:"functions"`
}

func (a *AgentStatement) statementNode() {}

// Goal represents the agent's goal
type Goal struct {
	BaseNode
	Value string `json:"value"`
}

func (g *Goal) expressionNode() {}

// Capabilities represents the agent's capabilities
type Capabilities struct {
	BaseNode
	Values []string `json:"values"`
}

// Event represents an event in a behavior block
type Event struct {
	BaseNode
	Name *Identifier `json:"name"`
}

func (e *Event) expressionNode() {}

// Behavior represents an action in a behavior block
type Behavior struct {
	BaseNode
	EventHandlers []*EventHandler `json:"event_handlers"`
}

func (b *Behavior) expressionNode() {}

// EventHandler represents an event handler in a behavior block
type EventHandler struct {
	BaseNode
	Event          *Event          `json:"event"`
	BlockStatement *BlockStatement `json:"block_statement"`
}

// FunctionArgument represents a function argument
type FunctionArgument struct {
	BaseNode
	Name *Identifier `json:"name"`
	Type *DataType   `json:"type"`
}

// Function represents a function declaration
type Function struct {
	BaseNode
	Name       *Identifier         `json:"name"`
	Arguments  []*FunctionArgument `json:"arguments"`
	Body       *BlockStatement     `json:"body"`
	ReturnType *DataType           `json:"return_type"`
}

func (f *Function) statementNode() {}

// ReturnStatement represents a return statement
type ReturnStatement struct {
	BaseNode
	Value *Expression `json:"value"`
}

func (rs *ReturnStatement) statementNode() {}

// BlockStatement represents a block of statements
type BlockStatement struct {
	BaseNode
	Statements map[int]*Statement `json:"statements"`
}

func (bs *BlockStatement) statementNode() {}

// VarStatement represents a variable declaration
type VarStatement struct {
	Statement
	Token lexer.Token `json:"token"`
	Name  *Identifier `json:"name"`
	Type  *DataType   `json:"type"`
	Value *Expression `json:"value"`
}

func (vs *VarStatement) statementNode() {}

// DataType represents a data type
type DataType struct {
	BaseNode
	Token lexer.Token `json:"token"`
}

// IdentifierLiteral represents an identifier literal
type IdentifierLiteral struct {
	BaseNode
	Value string `json:"value"`
}

func (il *IdentifierLiteral) expressionNode() {}

// IntegerLiteral represents an integer literal
type IntegerLiteral struct {
	BaseNode
	Value int64 `json:"value"`
}

func (il *IntegerLiteral) expressionNode() {}

// FloatLiteral represents a float literal
type FloatLiteral struct {
	BaseNode
	Value float64 `json:"value"`
}

func (fl *FloatLiteral) expressionNode() {}

// StringLiteral represents a string literal
type StringLiteral struct {
	BaseNode
	Value string `json:"value"`
}

func (sl *StringLiteral) expressionNode() {}

// BooleanLiteral represents a boolean literal
type BooleanLiteral struct {
	BaseNode
	Value bool `json:"value"`
}

func (b *BooleanLiteral) expressionNode() {}

// InfixExpression represents binary operations like 42 * 7
type InfixExpression struct {
	BaseNode
	Left     *Expression  `json:"left"`
	Operator *lexer.Token `json:"operator"`
	Right    *Expression  `json:"right"`
}

func (ie *InfixExpression) expressionNode() {}

// CallExpression represents a function call
type CallExpression struct {
	BaseNode
	Function  *Expression   `json:"function"`
	Arguments []*Expression `json:"arguments"`
}

func (ce *CallExpression) expressionNode() {}

// ExpressionStatement represents an expression statement
type ExpressionStatement struct {
	Statement
	BaseNode
	Expression *Expression `json:"expression"`
}

// TokenLiteral implements Statement.
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) statementNode() {}
