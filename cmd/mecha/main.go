package main

import (
	"flag"
	"fmt"
	"mechanus-compiler/internal/compiler_error"
	"mechanus-compiler/internal/parser"
	"os"
)

var debug bool = false

func main() {
	// Collect source and output files
	sourceFile, outputFile, err := getFiles()
	errSalt := "main"

	if err != nil {
		os.Exit(1)
	}
	defer func() {
		if err := sourceFile.Close(); err != nil {
			err = compiler_error.FileErrorf(errSalt, err)
			compiler_error.LogError(err)
		}
		if err := outputFile.Close(); err != nil {
			err = compiler_error.FileErrorf(errSalt, err)
			compiler_error.LogError(err)
		}
	}()

	// Initialize the parser
	parser, err := parser.NewParser(sourceFile, outputFile, debug)

	if err != nil {
		os.Exit(1)
	}

	// Start the syntax analysis
	if err = parser.Run(); err != nil {
		os.Exit(1)
	}
}

func getFiles() (*os.File, *os.File, error) {
	// Get source file path
	filePaths, err := getFilePaths()
	if err != nil {
		return nil, nil, err
	}

	// Open source file
	sourceFile, err := os.Open(filePaths[0])

	if err != nil {
		err = compiler_error.FileErrorf("getFiles", err)
		compiler_error.LogError(err)
		return nil, nil, err
	}

	// Create output file
	outputFile, err := os.Create(filePaths[1])

	if err != nil {
		err = compiler_error.FileErrorf("getFiles", err)
		compiler_error.LogError(err)
		return nil, nil, err
	}

	return sourceFile, outputFile, nil
}

func getFilePaths() ([]string, error) {
	// Define flags
	inputFile := flag.String("i", "", "Source file path")
	outputFile := flag.String("o", "", "Output file path")
	flag.BoolVar(&debug, "d", false, "Debug mode")

	// Parse command line arguments
	flag.Parse()

	// Check if required flags are provided
	if *inputFile == "" {
		err := compiler_error.FileErrorf("getFilePaths", fmt.Errorf(compiler_error.NoSourceFile))
		compiler_error.LogError(err)
		return nil, err
	}

	// If the output file name was not provided, default to "output"
	if *outputFile == "" {
		*outputFile = "output"
	}

	return []string{*inputFile, *outputFile}, nil
}
