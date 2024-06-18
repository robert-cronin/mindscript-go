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
)

type Scope struct {
	variables map[string]string
	functions map[string]string
	parent    *Scope
}

type SymbolTable struct {
	currentScope *Scope
}

func NewSymbolTable() *SymbolTable {
	globalScope := &Scope{variables: make(map[string]string)}
	return &SymbolTable{currentScope: globalScope}
}

func (st *SymbolTable) pushScope() {
	newScope := &Scope{variables: make(map[string]string), parent: st.currentScope}
	st.currentScope = newScope
}

func (st *SymbolTable) popScope() {
	if st.currentScope.parent == nil {
		panic("cannot pop the global scope")
	}
	st.currentScope = st.currentScope.parent
}

// DeclareVariable adds a new variable to the current scope
func (st *SymbolTable) DeclareVariable(name string, varType string) error {
	if _, exists := st.currentScope.variables[name]; exists {
		return errors.New("variable already declared in this scope")
	}
	st.currentScope.variables[name] = varType
	return nil
}

// GetVariableType returns the type of a variable
func (st *SymbolTable) GetVariableType(name string) (string, error) {
	for scope := st.currentScope; scope != nil; scope = scope.parent {
		if varType, exists := scope.variables[name]; exists {
			return varType, nil
		}
	}
	return "", errors.New("variable not declared")
}

// CheckVariableType checks if the type of a variable matches the expected type
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

// DeclareFunction adds a new function to the current scope
func (st *SymbolTable) DeclareFunction(name string, returnType string) error {
	return st.DeclareVariable(name, returnType)
}