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

package main

import (
	"encoding/json"
	"os"

	"github.com/robert-cronin/mindscript-go/pkg/codegen"
	"github.com/robert-cronin/mindscript-go/pkg/lexer"
	"github.com/robert-cronin/mindscript-go/pkg/logger"
	"github.com/robert-cronin/mindscript-go/pkg/parser"
	"github.com/robert-cronin/mindscript-go/pkg/semantic"
	"github.com/robert-cronin/mindscript-go/pkg/vm"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func dumpProgramToJson(program *parser.Program) (string, error) {
	jsonData, err := json.MarshalIndent(program, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func main() {
	var inputFile string
	var outputFile string
	var logLevel string

	var rootCmd = &cobra.Command{
		Use: "mindcript",
		Run: func(cmd *cobra.Command, args []string) {
			switch logLevel {
			case "debug":
				logger.Init(zap.DebugLevel)
			case "warn":
				logger.Init(zap.WarnLevel)
			case "error":
				logger.Init(zap.ErrorLevel)
			case "info":
				fallthrough
			default:
				logger.Init(zap.InfoLevel)
			}

			logger.Log.Info("msc: Starting")

			// Check if inputFile and outputFile are provided
			if inputFile == "" {
				logger.Log.Error("Input file not provided")
				os.Exit(1)
			}
			if outputFile == "" {
				// default output file, strips the extension
				outputFile = inputFile[:len(inputFile)-3] + ".mind"
			}
			logger.Log.Info("Processing files", zap.String("input", inputFile), zap.String("output", outputFile))

			// Read input file
			input, err := os.ReadFile(inputFile)
			if err != nil {
				logger.Log.Error("Error reading input file", zap.Error(err))
				os.Exit(1)
			}

			inputStr := string(input)

			l := lexer.New(inputStr)
			p := parser.New(l)
			program := p.ParseProgram()

			// Analyse the program
			st := semantic.NewSymbolTable(l)
			err = st.Analyse(program)
			if err != nil {
				logger.Log.Error("Error analyzing program", zap.Error(err))
				os.Exit(1)
			}

			// Generate bytecode
			instructions := codegen.GenerateBytecode(program, st)

			// Run VM
			virtualMachine := vm.New(instructions)
			virtualMachine.Run()

			jsonOutput, err := dumpProgramToJson(program)
			if err != nil {
				logger.Log.Error("Error dumping program to JSON", zap.Error(err))
				os.Exit(1)
			}

			// Write output to file
			jsonDumpFile := outputFile + ".json"
			err = os.WriteFile(jsonDumpFile, []byte(jsonOutput), 0644)
			if err != nil {
				logger.Log.Error("Error writing JSON dump file", zap.Error(err))
				os.Exit(1)
			}

			// Finished
			logger.Log.Info("msc: Finished")
		},
	}

	rootCmd.PersistentFlags().StringVarP(&inputFile, "input", "i", "", "Input file")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "Output file")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "error", "Log level (debug, info, warn, error)")

	err := rootCmd.Execute()
	if err != nil {
		logger.Log.Error("Error executing command", zap.Error(err))
		os.Exit(1)
	}
}
