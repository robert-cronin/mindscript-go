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

type Node interface {
	TokenLiteral() string
}

// Program represents the entire program
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
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
	Token lexer.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// Agent represents an agent declaration
type Agent struct {
	Token        lexer.Token
	Name         *Identifier
	Goal         *Goal
	Capabilities *Capabilities
	Behaviors    []*Behavior
	Functions    []*Function
}

func (a *Agent) statementNode()       {}
func (a *Agent) TokenLiteral() string { return a.Token.Literal }

// Goal represents the agent's goal
type Goal struct {
	Token lexer.Token
	Value string
}

func (g *Goal) expressionNode()      {}
func (g *Goal) TokenLiteral() string { return g.Token.Literal }

// Capabilities represents the agent's capabilities
type Capabilities struct {
	Token  lexer.Token
	Values []string
}

// Event represents an event in a behavior block
type Event struct {
	Token lexer.Token
	Name  *Identifier
}

func (e *Event) expressionNode()      {}
func (e *Event) TokenLiteral() string { return e.Token.Literal }

// Behavior represents an action in a behavior block
type Behavior struct {
	Token lexer.Token
	Name  *Identifier
}

func (a *Behavior) expressionNode()      {}
func (a *Behavior) TokenLiteral() string { return a.Token.Literal }

// Function represents a function declaration
type Function struct {
	Token      lexer.Token
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
	ReturnType string
}

func (f *Function) statementNode()       {}
func (f *Function) TokenLiteral() string { return f.Token.Literal }

// BlockStatement represents a block of statements
type BlockStatement struct {
	Token      lexer.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
