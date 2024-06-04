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

type SymbolTable struct {
	variables map[string]string
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{variables: make(map[string]string)}
}

func (st *SymbolTable) DeclareVariable(name string, varType string) error {
	if _, exists := st.variables[name]; exists {
		return errors.New("variable already declared")
	}
	st.variables[name] = varType
	return nil
}

func (st *SymbolTable) GetVariableType(name string) (string, error) {
	varType, exists := st.variables[name]
	if !exists {
		return "", errors.New("variable not declared")
	}
	return varType, nil
}

func (st *SymbolTable) CheckVariableType(name string, expectedType string) error {
	varType, err := st.GetVariableType(name)
	if err != nil {
		return err
	}
	if varType != expectedType {
		return errors.New(fmt.Sprintf("type mismatch: expected %s but got %s", expectedType, varType))
	}
	return nil
}

func Analyze(program *parser.Program) error {
	st := NewSymbolTable()

	for _, stmt := range program.Statements {
		switch stmt := stmt.(type) {
		case *parser.VarStatement:
			// TODO: support other types
			err := st.DeclareVariable(stmt.Name.Value, "int")
			if err != nil {
				return err
			}
		default:
			return errors.New("unsupported statement type")
		}
	}

	return nil
}
