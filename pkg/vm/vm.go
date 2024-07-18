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

package vm

import (
	"fmt"

	"go.uber.org/zap"
)

type Opcode int

const (
	// Arithmetic operations
	OpAdd Opcode = iota
	OpSub
	OpMul
	OpDiv

	// Stack operations
	OpPush
	OpPop

	// I/O operations
	OpPrint

	// Control flow
	OpHalt
	OpJump
	OpJumpIfFalse

	// Variable operations
	OpSetLocal
	OpGetLocal

	// Function operations
	OpCall
	OpReturn

	// Agent operations
	OpCreateAgent
	OpSetAgentGoal
	OpAddAgentCapability
	OpCreateEventHandler
	OpSetEventHandlerEvent
	OpAddAgentEventHandler
	OpCreateFunction
	OpAddFunctionArgument
	OpAddAgentFunction

	// Comparison operations
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpLessThan
	OpGreaterThanOrEqual
	OpLessThanOrEqual

	// Logical operations
	OpAnd
	OpOr
	OpNot

	// Type-specific operations
	OpConcatString
	OpPushString

	// Built-in function calls
	OpSyscall
	OpExec
	OpLog

	// Data structure operations
	OpCreateList
	OpAppendList
	OpGetListItem
	OpSetListItem
)

type Instruction struct {
	Opcode  Opcode
	Operand int
}

type VM struct {
	logger       *zap.Logger
	stack        []interface{}
	locals       []interface{}
	pc           int
	instructions []Instruction
	running      bool
	callStack    []int
}

func New(instructions []Instruction) *VM {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to initialize Zap logger: " + err.Error())
	}
	return &VM{
		logger:       logger,
		stack:        make([]interface{}, 0),
		locals:       make([]interface{}, 256),
		instructions: instructions,
		running:      true,
		callStack:    make([]int, 0),
	}
}

// Run starts the VM and executes the bytecode instructions
func (vm *VM) Run() {
	for vm.running {
		vm.step()
	}
}

func (vm *VM) step() {
	if vm.pc >= len(vm.instructions) {
		vm.running = false
		return
	}

	instr := vm.instructions[vm.pc]
	switch instr.Opcode {
	case OpAdd, OpSub, OpMul, OpDiv:
		vm.executeBinaryOp(instr.Opcode)
	case OpPush:
		vm.stack = append(vm.stack, instr.Operand)
	case OpPop:
		vm.popStack()
	case OpPrint:
		vm.logger.Debug("", zap.Any("value", vm.popStack()))
	case OpSetLocal:
		value := vm.popStack()
		vm.locals[instr.Operand] = value
	case OpGetLocal:
		value := vm.locals[instr.Operand]
		vm.stack = append(vm.stack, value)
	case OpCall:
		vm.callStack = append(vm.callStack, vm.pc+1)
		vm.pc = instr.Operand
		return
	case OpReturn:
		if len(vm.callStack) == 0 {
			vm.running = false
			return
		}
		vm.pc = vm.callStack[len(vm.callStack)-1]
		vm.callStack = vm.callStack[:len(vm.callStack)-1]
		return
	case OpHalt:
		vm.running = false
	case OpCreateAgent:
		// TODO: Implement agent creation logic
		vm.logger.Debug("Creating agent")
	case OpSetAgentGoal:
		// TODO: Implement setting agent goal logic
		vm.logger.Debug("Setting agent goal")
	case OpAddAgentCapability:
		// TODO: Implement adding agent capability logic
		vm.logger.Debug("Adding agent capability")
	case OpCreateEventHandler:
		// TODO: Implement event handler creation logic
		vm.logger.Debug("Creating event handler")
	case OpSetEventHandlerEvent:
		// TODO: Implement setting event handler event logic
		vm.logger.Debug("Setting event handler event")
	case OpAddAgentEventHandler:
		// TODO: Implement adding event handler to agent logic
		vm.logger.Debug("Adding event handler to agent")
	case OpCreateFunction:
		// TODO: Implement function creation logic
		vm.logger.Debug("Creating function")
	case OpAddFunctionArgument:
		// TODO: Implement adding function argument logic
		vm.logger.Debug("Adding function argument")
	case OpAddAgentFunction:
		// TODO: Implement adding function to agent logic
		vm.logger.Debug("Adding function to agent")
	case OpSyscall:
		// TODO: Implement syscall logic
		vm.logger.Debug("Executing syscall")
	case OpExec:
		// TODO: Implement exec logic
		vm.logger.Debug("Executing external command")
	case OpLog:
		// TODO: Implement log logic
		vm.logger.Debug("Logging:", zap.Any("value", vm.popStack()))
	case OpPushString:
		// Implement string pushing logic
		stringValue := vm.getStringConstant(instr.Operand)
		vm.stack = append(vm.stack, stringValue)
		vm.logger.Debug("Pushing string:", zap.String("value", stringValue))
	default:
		fmt.Printf("Error: Unknown opcode %d\n", instr.Opcode)
		vm.running = false
	}

	vm.pc++
}

func (vm *VM) getStringConstant(index int) string {
	// TODO: Implement string constant retrieval logic
	return fmt.Sprintf("String constant %d", index)
}

// executeBinaryOp executes a binary operation
func (vm *VM) executeBinaryOp(opcode Opcode) {
	right := vm.popStack()
	left := vm.popStack()

	var result interface{}

	switch opcode {
	case OpAdd:
		result = vm.add(left, right)
	case OpSub:
		result = vm.sub(left, right)
	case OpMul:
		result = vm.mul(left, right)
	case OpDiv:
		result = vm.div(left, right)
	}

	vm.stack = append(vm.stack, result)
}

// popStack pops the top value from the stack
func (vm *VM) popStack() interface{} {
	if len(vm.stack) == 0 {
		panic("Stack underflow")
	}
	value := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return value
}

func (vm *VM) add(a, b interface{}) interface{} {
	switch x := a.(type) {
	case int:
		switch y := b.(type) {
		case int:
			return x + y
		case float64:
			return float64(x) + y
		}
	case float64:
		switch y := b.(type) {
		case int:
			return x + float64(y)
		case float64:
			return x + y
		}
	}
	panic(fmt.Sprintf("Unsupported types for addition: %T and %T", a, b))
}

func (vm *VM) sub(a, b interface{}) interface{} {
	switch x := a.(type) {
	case int:
		switch y := b.(type) {
		case int:
			return x - y
		case float64:
			return float64(x) - y
		}
	case float64:
		switch y := b.(type) {
		case int:
			return x - float64(y)
		case float64:
			return x - y
		}
	}
	panic(fmt.Sprintf("Unsupported types for subtraction: %T and %T", a, b))
}

func (vm *VM) mul(a, b interface{}) interface{} {
	switch x := a.(type) {
	case int:
		switch y := b.(type) {
		case int:
			return x * y
		case float64:
			return float64(x) * y
		}
	case float64:
		switch y := b.(type) {
		case int:
			return x * float64(y)
		case float64:
			return x * y
		}
	}
	panic(fmt.Sprintf("Unsupported types for multiplication: %T and %T", a, b))
}

func (vm *VM) div(a, b interface{}) interface{} {
	switch x := a.(type) {
	case int:
		switch y := b.(type) {
		case int:
			if y == 0 {
				panic("Division by zero")
			}
			return x / y
		case float64:
			if y == 0 {
				panic("Division by zero")
			}
			return float64(x) / y
		}
	case float64:
		switch y := b.(type) {
		case int:
			if y == 0 {
				panic("Division by zero")
			}
			return x / float64(y)
		case float64:
			if y == 0 {
				panic("Division by zero")
			}
			return x / y
		}
	}
	panic(fmt.Sprintf("Unsupported types for division: %T and %T", a, b))
}
