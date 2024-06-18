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

package semantic

import (
	"errors"
	"fmt"
	"github.com/robert-cronin/mindscript-go/pkg/parser"
)

func (st *SymbolTable) Analyse(program *parser.Program) error {
	st.initSystemFunctions()
	for _, stmt := range program.Statements {
		if err := st.analyseStatement(stmt); err != nil {
			return err
		}
	}
	return nil
}

// Initialise the system functions like log and syscall
func (st *SymbolTable) initSystemFunctions() {
	st.DeclareFunction("log", FunctionSignature{
		Arguments:  []string{"string"},
		ReturnType: "void",
	})
	st.DeclareFunction("syscall", FunctionSignature{
		Arguments:  []string{"string", "string"},
		ReturnType: "void",
	})
}

func (st *SymbolTable) analyseStatement(stmt parser.Statement) error {
	switch s := stmt.(type) {
	case *parser.AgentStatement:
		if err := st.DeclareVariable(s.Name.Value, "agent"); err != nil {
			return err
		}
		if err := st.analyseAgentStatement(s); err != nil {
			return err
		}
	case *parser.VarStatement:
		if err := st.DeclareVariable(s.Name.Value, s.Type.TokenLiteral()); err != nil {
			return err
		}
		return st.analyseExpression(*s.Value)
	case *parser.Function:
		signature := FunctionSignature{
			Arguments:  st.getArgumentsTypes(s.Arguments),
			ReturnType: s.ReturnType.TokenLiteral(),
		}
		if err := st.DeclareFunction(s.Name.Value, signature); err != nil {
			return err
		}
		st.pushScope()
		for _, arg := range s.Arguments {
			if err := st.DeclareVariable(arg.Name.Value, arg.Type.TokenLiteral()); err != nil {
				return err
			}
		}
		for _, stmt := range s.Body.Statements {
			if err := st.analyseStatement(*stmt); err != nil {
				return err
			}
		}
		st.popScope()
	case *parser.ExpressionStatement:
		return st.analyseExpression(*s.Expression)
	case *parser.ReturnStatement:
		return st.analyseExpression(*s.Value)
	}
	return nil
}

func (st *SymbolTable) analyseAgentStatement(agent *parser.AgentStatement) error {
	for _, behavior := range agent.Behaviors {
		for _, eventHandler := range behavior.EventHandlers {
			st.pushScope()
			if err := st.analyseBlockStatement(eventHandler.BlockStatement); err != nil {
				return err
			}
			st.popScope()
		}
	}
	for _, function := range agent.Functions {
		if err := st.analyseStatement(function); err != nil {
			return err
		}
	}
	return nil
}

func (st *SymbolTable) analyseBlockStatement(block *parser.BlockStatement) error {
	for _, stmt := range block.Statements {
		if err := st.analyseStatement(*stmt); err != nil {
			return err
		}
	}
	return nil
}

func (st *SymbolTable) analyseExpression(expr parser.Expression) error {
	switch e := expr.(type) {
	case *parser.IdentifierLiteral:
		if _, err := st.GetVariableType(e.Value); err != nil {
			return err
		}
	case *parser.InfixExpression:
		if err := st.analyseExpression(*e.Left); err != nil {
			return err
		}
		if err := st.analyseExpression(*e.Right); err != nil {
			return err
		}
	case *parser.CallExpression:
		funcName := (*e.Function).(*parser.IdentifierLiteral).Value
		funcSig, err := st.GetFunctionSignature(funcName)
		if err != nil {
			return fmt.Errorf("line %d: %s", st.l.Line(e.Token), err)
		}
		if len(funcSig.Arguments) != len(e.Arguments) {
			return fmt.Errorf("line %d: expected %d arguments but got %d", st.l.Line(e.Token), len(funcSig.Arguments), len(e.Arguments))
		}
		for i, arg := range e.Arguments {
			if err := st.analyseExpression(*arg); err != nil {
				return fmt.Errorf("line %d: %s", st.l.Line(e.Token), err)
			}
			argType, err := st.getExpressionType(*arg)
			if err != nil {
				return fmt.Errorf("line %d: %s", st.l.Line(e.Token), err)
			}
			if funcSig.Arguments[i] != argType {
				return fmt.Errorf("line %d: type mismatch for argument %d: expected %s but got %s", st.l.Line(e.Token), i+1, funcSig.Arguments[i], argType)
			}
		}
	}
	return nil
}

func (st *SymbolTable) getArgumentsTypes(args []*parser.FunctionArgument) []string {
	types := []string{}
	for _, arg := range args {
		types = append(types, arg.Type.TokenLiteral())
	}
	return types
}

func (st *SymbolTable) getExpressionType(expr parser.Expression) (string, error) {
	switch e := expr.(type) {
	case *parser.IdentifierLiteral:
		return st.GetVariableType(e.Value)
	case *parser.IntegerLiteral:
		return "int", nil
	case *parser.FloatLiteral:
		return "float", nil
	case *parser.StringLiteral:
		return "string", nil
	case *parser.BooleanLiteral:
		return "bool", nil
	case *parser.InfixExpression:
		leftType, err := st.getExpressionType(*e.Left)
		if err != nil {
			return "", err
		}
		rightType, err := st.getExpressionType(*e.Right)
		if err != nil {
			return "", err
		}
		if leftType != rightType {
			return "", errors.New("type mismatch in infix expression")
		}
		return leftType, nil
	case *parser.CallExpression:
		funcName := (*e.Function).(*parser.IdentifierLiteral).Value
		funcSig, err := st.GetFunctionSignature(funcName)
		if err != nil {
			return "", err
		}
		return funcSig.ReturnType, nil
	}
	return "", errors.New("unknown expression type")
}
