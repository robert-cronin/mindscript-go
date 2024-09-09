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
	"os/exec"
	"strings"

	"github.com/robert-cronin/mindscript-go/pkg/logger"
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
	stack           []interface{}
	locals          []interface{}
	pc              int
	instructions    []Instruction
	running         bool
	callStack       []int
	stringConstants []string
}

func New(instructions []Instruction) *VM {
	return &VM{
		stack:           make([]interface{}, 0),
		locals:          make([]interface{}, 256),
		instructions:    instructions,
		running:         true,
		callStack:       make([]int, 0),
		stringConstants: make([]string, 0),
	}
}

// Run starts the VM and executes the bytecode instructions
func (vm *VM) Run() {
	logger.Log.Info("Starting VM execution")
	for vm.running {
		vm.step()
	}
	logger.Log.Info("VM execution completed")
}

func (vm *VM) step() {
	if vm.pc >= len(vm.instructions) {
		vm.running = false
		logger.Log.Info("Reached end of instructions", zap.Int("pc", vm.pc))
		return
	}

	instr := vm.instructions[vm.pc]
	logger.Log.Debug("Executing instruction", zap.Int("pc", vm.pc), zap.Any("instruction", instr))

	switch instr.Opcode {
	case OpAdd, OpSub, OpMul, OpDiv:
		vm.executeBinaryOp(instr.Opcode)
	case OpPush:
		vm.stack = append(vm.stack, instr.Operand)
		logger.Log.Debug("Pushed value to stack", zap.Any("value", instr.Operand))
	case OpPop:
		value := vm.popStack()
		logger.Log.Debug("Popped value from stack", zap.Any("value", value))
	case OpPrint:
		value := vm.popStack()
		fmt.Println(value)
		logger.Log.Debug("Printed value", zap.Any("value", value))
	case OpSetLocal:
		value := vm.popStack()
		vm.locals[instr.Operand] = value
		logger.Log.Debug("Set local variable", zap.Int("index", instr.Operand), zap.Any("value", value))
	case OpGetLocal:
		value := vm.locals[instr.Operand]
		vm.stack = append(vm.stack, value)
		logger.Log.Debug("Got local variable", zap.Int("index", instr.Operand), zap.Any("value", value))
	case OpCall:
		vm.callStack = append(vm.callStack, vm.pc+1)
		vm.pc = instr.Operand
		logger.Log.Debug("Function call", zap.Int("returnAddress", vm.pc+1), zap.Int("functionAddress", instr.Operand))
		return
	case OpReturn:
		if len(vm.callStack) == 0 {
			vm.running = false
			logger.Log.Info("Return from main function, halting VM")
			return
		}
		vm.pc = vm.callStack[len(vm.callStack)-1]
		vm.callStack = vm.callStack[:len(vm.callStack)-1]
		logger.Log.Debug("Function return", zap.Int("returnAddress", vm.pc))
		return
	case OpHalt:
		vm.running = false
		logger.Log.Info("Halt instruction encountered, stopping VM")
	case OpCreateAgent:
		logger.Log.Debug("Creating agent", zap.Int("agentIndex", instr.Operand))
		// TODO: Implement actual agent creation logic
	case OpSetAgentGoal:
		goal := vm.popStack()
		logger.Log.Debug("Setting agent goal", zap.Int("agentIndex", instr.Operand), zap.Any("goal", goal))
		// TODO: Implement actual agent goal setting logic
	case OpAddAgentCapability:
		capability := vm.popStack()
		logger.Log.Debug("Adding agent capability", zap.Int("agentIndex", instr.Operand), zap.Any("capability", capability))
		// TODO: Implement actual agent capability adding logic
	case OpCreateEventHandler:
		logger.Log.Debug("Creating event handler", zap.Int("handlerIndex", instr.Operand))
		// TODO: Implement actual event handler creation logic
	case OpSetEventHandlerEvent:
		event := vm.popStack()
		logger.Log.Debug("Setting event handler event", zap.Int("handlerIndex", instr.Operand), zap.Any("event", event))
		// TODO: Implement actual event handler event setting logic
	case OpAddAgentEventHandler:
		handlerIndex := vm.popStack()
		logger.Log.Debug("Adding event handler to agent", zap.Int("agentIndex", instr.Operand), zap.Any("handlerIndex", handlerIndex))
		// TODO: Implement actual logic to add event handler to agent
	case OpCreateFunction:
		logger.Log.Debug("Creating function", zap.Int("functionIndex", instr.Operand))
		// TODO: Implement actual function creation logic
	case OpAddFunctionArgument:
		argName := vm.popStack()
		logger.Log.Debug("Adding function argument", zap.Int("functionIndex", instr.Operand), zap.Any("argumentName", argName))
		// TODO: Implement actual function argument adding logic
	case OpAddAgentFunction:
		functionIndex := vm.popStack()
		logger.Log.Debug("Adding function to agent", zap.Int("agentIndex", instr.Operand), zap.Any("functionIndex", functionIndex))
		// TODO: Implement actual logic to add function to agent
	case OpSyscall:
		command := vm.popStack().(string)
		args := vm.popStack().(string)
		logger.Log.Debug("Executing syscall", zap.String("command", command), zap.String("args", args))
		cmd := exec.Command(command, strings.Split(args, " ")...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			logger.Log.Error("Syscall failed", zap.Error(err))
		} else {
			logger.Log.Debug("Syscall output", zap.String("output", string(output)))
		}
	case OpExec:
		command := vm.popStack().(string)
		args := vm.popStack().(string)
		logger.Log.Debug("Executing external command", zap.String("command", command), zap.String("args", args))
		cmd := exec.Command(command, strings.Split(args, " ")...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			logger.Log.Error("External command failed", zap.Error(err))
		} else {
			vm.stack = append(vm.stack, string(output))
			logger.Log.Debug("External command output", zap.String("output", string(output)))
		}
	case OpLog:
		message := vm.popStack()
		logger.Log.Info("Log message", zap.Any("message", message))
	case OpPushString:
		stringValue := vm.getStringConstant(instr.Operand)
		vm.stack = append(vm.stack, stringValue)
		logger.Log.Debug("Pushed string to stack", zap.String("value", stringValue))
	default:
		logger.Log.Error("Unknown opcode", zap.Int("opcode", int(instr.Opcode)))
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
		logger.Log.Error("Attempted to pop from empty stack")
		vm.running = false
		return nil
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

func (vm *VM) AddStringConstant(s string) int {
	vm.stringConstants = append(vm.stringConstants, s)
	return len(vm.stringConstants) - 1
}

func (vm *VM) GetLastResult() interface{} {
	if len(vm.stack) > 0 {
		return vm.stack[len(vm.stack)-1]
	}
	return nil
}
