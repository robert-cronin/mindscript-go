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
	Token lexer.Token
}

func (b *BaseNode) TokenLiteral() string {
	return b.Token.Literal
}

// Program represents the entire program
type Program struct {
	Statements []Statement
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
	Value string
}

func (i *Identifier) expressionNode() {}

// Agent represents an agent declaration
type Agent struct {
	BaseNode
	Name         *Identifier
	Goal         *Goal
	Capabilities *Capabilities
	Behaviors    []*Behavior
	Functions    []*Function
}

func (a *Agent) statementNode() {}

// Goal represents the agent's goal
type Goal struct {
	BaseNode
	Value string
}

func (g *Goal) expressionNode() {}

// Capabilities represents the agent's capabilities
type Capabilities struct {
	BaseNode
	Values []string
}

// Event represents an event in a behavior block
type Event struct {
	BaseNode
	Name *Identifier
}

func (e *Event) expressionNode() {}

// Behavior represents an action in a behavior block
type Behavior struct {
	BaseNode
	EventHandlers []*EventHandler
}

func (b *Behavior) expressionNode() {}

// EventHandler represents an event handler in a behavior block
type EventHandler struct {
	BaseNode
	Event          *Event
	BlockStatement *BlockStatement
}

// Function represents a function declaration
type Function struct {
	BaseNode
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
	ReturnType string
}

func (f *Function) statementNode() {}

// BlockStatement represents a block of statements
type BlockStatement struct {
	BaseNode
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// VarStatement represents a variable declaration
type VarStatement struct {
	Statement
	Token lexer.Token
	Name  *Identifier
	Type  *DataType
	Value *Expression
}

func (vs *VarStatement) statementNode() {}

// DataType
type DataType struct {
	BaseNode
	Token lexer.Token
}

// IdentifierLiteral
type IdentifierLiteral struct {
	BaseNode
	Value string
}

func (il *IdentifierLiteral) expressionNode() {}

// IntegerLiteral
type IntegerLiteral struct {
	BaseNode
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// FloatLiteral
type FloatLiteral struct {
	BaseNode
	Value float64
}

func (fl *FloatLiteral) expressionNode() {}

// StringLiteral
type StringLiteral struct {
	BaseNode
	Value string
}

func (sl *StringLiteral) expressionNode() {}

// Boolean
type BooleanLiteral struct {
	BaseNode
	Value bool
}

func (b *BooleanLiteral) expressionNode() {}

// InfixExpression represents binary operations like 42 * 7
type InfixExpression struct {
	BaseNode
	Left     Expression
	Operator lexer.Token
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
