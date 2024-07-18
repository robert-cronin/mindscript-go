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
	"fmt"
	"os"

	"github.com/robert-cronin/mindscript-go/pkg/codegen"
	"github.com/robert-cronin/mindscript-go/pkg/lexer"
	"github.com/robert-cronin/mindscript-go/pkg/logger"
	"github.com/robert-cronin/mindscript-go/pkg/parser"
	"github.com/robert-cronin/mindscript-go/pkg/repl"
	"github.com/robert-cronin/mindscript-go/pkg/semantic"
	"github.com/robert-cronin/mindscript-go/pkg/vm"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	inputFile  string
	outputFile string
	logLevel   string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "msc",
		Short: "MindScript Compiler",
		Long:  `MindScript Compiler is a tool for compiling and running MindScript code.`,
	}

	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "info", "Log level (debug, info, warn, error)")

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Build MindScript code",
		Run:   runBuild,
	}

	buildCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file")
	buildCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file")
	buildCmd.MarkFlagRequired("input")

	replCmd := &cobra.Command{
		Use:   "repl",
		Short: "Start MindScript REPL",
		Run:   runRepl,
	}

	rootCmd.AddCommand(buildCmd, replCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initLogger() {
	var zapLevel zapcore.Level
	switch logLevel {
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	default:
		zapLevel = zap.InfoLevel
	}
	logger.Init(zapLevel)
}

func runBuild(cmd *cobra.Command, args []string) {
	initLogger()
	logger.Log.Info("msc: Starting build")

	if outputFile == "" {
		outputFile = inputFile[:len(inputFile)-3] + ".mind"
	}
	logger.Log.Info("Processing files", zap.String("input", inputFile), zap.String("output", outputFile))

	input, err := os.ReadFile(inputFile)
	if err != nil {
		logger.Log.Error("Error reading input file", zap.Error(err))
		os.Exit(1)
	}

	inputStr := string(input)
	l := lexer.New(inputStr)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		logger.Log.Error("Parser errors", zap.Strings("errors", p.Errors()))
		os.Exit(1)
	}

	st := semantic.NewSymbolTable(l)
	err = st.Analyse(program)
	if err != nil {
		logger.Log.Error("Error analyzing program", zap.Error(err))
		os.Exit(1)
	}

	instructions := codegen.GenerateBytecode(program, st)

	virtualMachine := vm.New(instructions)
	virtualMachine.Run()

	jsonOutput, err := dumpProgramToJson(program)
	if err != nil {
		logger.Log.Error("Error dumping program to JSON", zap.Error(err))
		os.Exit(1)
	}

	jsonDumpFile := outputFile + ".json"
	err = os.WriteFile(jsonDumpFile, []byte(jsonOutput), 0644)
	if err != nil {
		logger.Log.Error("Error writing JSON dump file", zap.Error(err))
		os.Exit(1)
	}

	logger.Log.Info("msc: Build finished")
}

func runRepl(cmd *cobra.Command, args []string) {
	initLogger()
	logger.Log.Info("msc: Starting REPL")
	repl.Start()
	logger.Log.Info("msc: REPL finished")
}

func dumpProgramToJson(program *parser.Program) (string, error) {
	jsonData, err := json.MarshalIndent(program, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
