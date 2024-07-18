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

package codegen

import (
	"fmt"

	"github.com/robert-cronin/mindscript-go/pkg/lexer"
	"github.com/robert-cronin/mindscript-go/pkg/parser"
	"github.com/robert-cronin/mindscript-go/pkg/semantic"
	"github.com/robert-cronin/mindscript-go/pkg/vm"
	"go.uber.org/zap"
)

type CodeGenerator struct {
	logger           *zap.Logger
	instructions     []vm.Instruction
	symbolTable      *semantic.SymbolTable
	functions        map[string]int
	symbols          map[string]int
	nextFuncIndex    int
	nextSymbolIndex  int
	builtinFunctions map[string]vm.Opcode
}

func NewCodeGenerator(symbolTable *semantic.SymbolTable) *CodeGenerator {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to initialize Zap logger: " + err.Error())
	}
	cg := &CodeGenerator{
		logger:          logger,
		instructions:    []vm.Instruction{},
		symbolTable:     symbolTable,
		functions:       make(map[string]int),
		symbols:         make(map[string]int),
		nextFuncIndex:   0,
		nextSymbolIndex: 0,
		builtinFunctions: map[string]vm.Opcode{
			"log":     vm.OpLog,
			"syscall": vm.OpSyscall,
			"exec":    vm.OpExec,
		},
	}
	return cg
}

func (cg *CodeGenerator) declareSymbol(name string) int {
	if index, exists := cg.symbols[name]; exists {
		return index
	}
	index := cg.nextSymbolIndex
	cg.symbols[name] = index
	cg.nextSymbolIndex++
	return index
}

func (cg *CodeGenerator) declareFunction(name string) int {
	if index, exists := cg.functions[name]; exists {
		return index
	}
	index := cg.nextFuncIndex
	cg.functions[name] = index
	cg.nextFuncIndex++
	return index
}

func (cg *CodeGenerator) generateAgentStatement(agent *parser.AgentStatement) {
	agentIndex := cg.declareSymbol(agent.Name.Value)
	cg.emit(vm.OpCreateAgent, agentIndex)

	if agent.Goal != nil {
		cg.generateStringLiteral(agent.Goal.Value)
		cg.emit(vm.OpSetAgentGoal, agentIndex)
	}

	if agent.Capabilities != nil {
		for _, capability := range agent.Capabilities.Values {
			cg.generateStringLiteral(capability)
			cg.emit(vm.OpAddAgentCapability, agentIndex)
		}
	}

	for _, behavior := range agent.Behaviors {
		cg.generateBehavior(behavior, agentIndex)
	}

	for _, function := range agent.Functions {
		cg.generateFunction(function, agentIndex)
	}
}

func (cg *CodeGenerator) generateBehavior(behavior *parser.Behavior, agentIndex int) {
	for _, eventHandler := range behavior.EventHandlers {
		eventHandlerIndex := cg.nextSymbolIndex
		cg.nextSymbolIndex++

		cg.emit(vm.OpCreateEventHandler, eventHandlerIndex)

		cg.generateStringLiteral(eventHandler.Event.Name.Value)
		cg.emit(vm.OpSetEventHandlerEvent, eventHandlerIndex)

		cg.generateBlockStatement(eventHandler.BlockStatement)

		cg.emit(vm.OpAddAgentEventHandler, agentIndex)
		cg.emit(vm.OpPush, eventHandlerIndex)
	}
}

func (cg *CodeGenerator) generateFunction(function *parser.Function, agentIndex int) {
	functionIndex := cg.declareFunction(function.Name.Value)

	cg.emit(vm.OpCreateFunction, functionIndex)

	for _, arg := range function.Arguments {
		cg.generateStringLiteral(arg.Name.Value)
		cg.emit(vm.OpAddFunctionArgument, functionIndex)
	}

	cg.generateBlockStatement(function.Body)

	cg.emit(vm.OpAddAgentFunction, agentIndex)
	cg.emit(vm.OpPush, functionIndex)
}

func (cg *CodeGenerator) generateBlockStatement(block *parser.BlockStatement) {
	for _, stmt := range block.Statements {
		cg.generateStatement(*stmt)
	}
}

func (cg *CodeGenerator) generateStatement(stmt parser.Statement) {
	switch s := stmt.(type) {
	case *parser.AgentStatement:
		cg.generateAgentStatement(s)
	case *parser.ExpressionStatement:
		cg.generateExpression(*s.Expression)
	case *parser.VarStatement:
		cg.generateVarStatement(s)
	case *parser.ReturnStatement:
		cg.generateExpression(*s.Value)
		cg.emit(vm.OpReturn, 0)
	default:
		// Handle unknown statement types
		cg.logger.Panic("Unsupported statement type", zap.String("type", fmt.Sprintf("%T", s)))
	}
}

func (cg *CodeGenerator) generateExpression(expr parser.Expression) {
	switch e := expr.(type) {
	case *parser.IntegerLiteral:
		cg.emit(vm.OpPush, int(e.Value))
	case *parser.FloatLiteral:
		// TODO: handle float literals and not just cast to int
		cg.emit(vm.OpPush, int(e.Value))
	case *parser.StringLiteral:
		cg.generateStringLiteral(e.Value)
	case *parser.BooleanLiteral:
		if e.Value {
			cg.emit(vm.OpPush, 1)
		} else {
			cg.emit(vm.OpPush, 0)
		}
	case *parser.IdentifierLiteral:
		varIndex, exists := cg.symbols[e.Value]
		if !exists {
			cg.logger.Panic("Undefined variable", zap.String("variable", e.Value))
		}
		cg.emit(vm.OpGetLocal, varIndex)
	case *parser.InfixExpression:
		cg.generateExpression(*e.Left)
		cg.generateExpression(*e.Right)
		switch e.Operator.Type {
		case lexer.PLUS:
			cg.emit(vm.OpAdd, 0)
		case lexer.MINUS:
			cg.emit(vm.OpSub, 0)
		case lexer.ASTERISK:
			cg.emit(vm.OpMul, 0)
		case lexer.SLASH:
			cg.emit(vm.OpDiv, 0)
		default:
			cg.logger.Panic("Unknown operator", zap.String("operator", e.Operator.Literal))
		}
	case *parser.CallExpression:
		for _, arg := range e.Arguments {
			cg.generateExpression(*arg)
		}
		funcName := (*e.Function).(*parser.IdentifierLiteral).Value
		if opcode, isBuiltin := cg.builtinFunctions[funcName]; isBuiltin {
			cg.emit(opcode, len(e.Arguments))
		} else {
			funcIndex, exists := cg.functions[funcName]
			if !exists {
				cg.logger.Panic("Undefined function", zap.String("function", funcName))
			}
			cg.emit(vm.OpCall, funcIndex)
		}
	default:
		cg.logger.Panic("Unsupported expression type", zap.String("type", fmt.Sprintf("%T", e)))
	}
}

func (cg *CodeGenerator) generateStringLiteral(value string) {
	// TODO: maybe store string literals in a separate table
	stringIndex := cg.declareSymbol(value)
	cg.emit(vm.OpPushString, stringIndex)
}

func (cg *CodeGenerator) generateVarStatement(stmt *parser.VarStatement) {
	cg.generateExpression(*stmt.Value)
	varIndex := cg.declareSymbol(stmt.Name.Value)
	cg.emit(vm.OpSetLocal, varIndex)
}

func (cg *CodeGenerator) emit(opcode vm.Opcode, operand int) {
	cg.instructions = append(cg.instructions, vm.Instruction{Opcode: opcode, Operand: operand})
}

// GenerateBytecode is the main function to generate bytecode from the AST
func GenerateBytecode(program *parser.Program, symbolTable *semantic.SymbolTable) []vm.Instruction {
	cg := NewCodeGenerator(symbolTable)
	for _, stmt := range program.Statements {
		cg.generateStatement(stmt)
	}
	cg.emit(vm.OpHalt, 0)
	return cg.instructions
}
