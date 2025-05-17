package main

import (
	"errors"
	"flag"
	"mechanus-compiler/src/error"
	"mechanus-compiler/src/models"
	"os"
)

func main() {
	// Collect source and output files
	sourceFile, outputFile, err := getFiles()
	if err != nil {
		custom_errors.Log(custom_errors.FileOpenError, &err, custom_errors.ErrorLevel)
		os.Exit(1)
	}
	defer func() {
		if err := sourceFile.Close(); err != nil {
			custom_errors.Log(custom_errors.FileCloseError, &err, custom_errors.ErrorLevel)
		}
		if err := outputFile.Close(); err != nil {
			custom_errors.Log(custom_errors.FileCloseError, &err, custom_errors.ErrorLevel)
		}
	}()

	// Initialize the parser
	parser := models.Parser{
		Source: sourceFile,
		Output: outputFile,
	}

	// Start the syntax analysis
	parser.Run()
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
		custom_errors.Log(custom_errors.FileOpenError, &err, custom_errors.ErrorLevel)
		return nil, nil, err
	}

	// Create output file
	outputFile, err := os.Create(filePaths[1])
	if err != nil {
		sourceFile.Close()
		custom_errors.Log(custom_errors.FileCreateError, &err, custom_errors.ErrorLevel)
		return nil, nil, err
	}

	return sourceFile, outputFile, nil
}

func getFilePaths() ([]string, error) {
	// Define flags
	inputFile := flag.String("i", "", "Source file path")
	outputFile := flag.String("o", "", "Output file path")

	// Parse command line arguments
	flag.Parse()

	// Check if required flags are provided
	if *inputFile == "" {
		err := errors.New(custom_errors.NoSourceFile)
		custom_errors.Log(custom_errors.NoSourceFile, &err, custom_errors.ErrorLevel)
		return nil, err
	}

	// If the output file name was not provided, default to "output"
	if *outputFile == "" {
		*outputFile = "output"
	}

	return []string{*inputFile, *outputFile}, nil
}
