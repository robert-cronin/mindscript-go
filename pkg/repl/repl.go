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

package repl

import (
	"bufio"
	"fmt"
	"os"

	"github.com/robert-cronin/mindscript-go/pkg/codegen"
	"github.com/robert-cronin/mindscript-go/pkg/lexer"
	"github.com/robert-cronin/mindscript-go/pkg/logger"
	"github.com/robert-cronin/mindscript-go/pkg/parser"
	"github.com/robert-cronin/mindscript-go/pkg/semantic"
	"github.com/robert-cronin/mindscript-go/pkg/vm"
	"go.uber.org/zap"
)

func Start() {
	fmt.Println("Welcome to the MindScript REPL!")
	fmt.Println("Type 'exit' to quit.")

	scanner := bufio.NewScanner(os.Stdin)
	symbolTable := semantic.NewSymbolTable(lexer.New(""))

	for {
		fmt.Print(">> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "exit" {
			break
		}

		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			for _, msg := range p.Errors() {
				logger.Log.Error("Parser error", zap.String("error", msg))
			}
			continue
		}

		err := symbolTable.Analyse(program)
		if err != nil {
			logger.Log.Error("Semantic error", zap.Error(err))
			continue
		}

		instructions := codegen.GenerateBytecode(program, symbolTable)
		virtualMachine := vm.New(instructions)
		virtualMachine.Run()

		result := virtualMachine.GetLastResult()
		fmt.Printf("%v\n", result)
	}

	fmt.Println("Goodbye!")
}
