package main

import (
	"mechanus-compiler/src/error"
	"mechanus-compiler/src/models"
	"os"
)

func main() {
	// Open source file
	sourceFile, err := os.Open("examples/example1.mecha")

	if err != nil {
		custom_errors.Log(custom_errors.FileOpenError, &err, custom_errors.ErrorLevel)
		return
	}
	defer func(inputFile *os.File) {
		err := inputFile.Close()
		if err != nil {
			custom_errors.Log(custom_errors.FileCloseError, &err, custom_errors.ErrorLevel)
		}
	}(sourceFile)

	// Create output file
	outputFile, err := os.Create("output.txt")

	if err != nil {
		custom_errors.Log(custom_errors.FileCreateError, &err, custom_errors.ErrorLevel)
		return
	}
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			custom_errors.Log(custom_errors.FileCloseError, &err, custom_errors.ErrorLevel)
		}
	}(outputFile)

	// Initialize the parser
	parser := models.Parser{
		Source: sourceFile,
		Output: outputFile,
	}

	// Start the syntax analysis
	parser.Run()
}
