package main

import (
	"errors"
	"flag"
	log "mechanus-compiler/src/error"
	"mechanus-compiler/src/models"
	"os"
)

func main() {
	// Collect source and output files
	sourceFile, outputFile, err := getFiles()
	if err != nil {
		err = log.EnrichError(err, log.FileOpenError)
		log.Log(err.Error(), log.ErrorLevel)
		os.Exit(1)
	}
	defer func() {
		if err := sourceFile.Close(); err != nil {
			err = log.EnrichError(err, log.FileCloseError)
			log.Log(err.Error(), log.ErrorLevel)
		}
		if err := outputFile.Close(); err != nil {
			err = log.EnrichError(err, log.FileCloseError)
			log.Log(err.Error(), log.ErrorLevel)
		}
	}()

	// Initialize the parser
	parser, err := models.NewParser(sourceFile, outputFile)

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
		err = log.EnrichError(err, log.FileOpenError)
		log.Log(err.Error(), log.ErrorLevel)
		return nil, nil, err
	}

	// Create output file
	outputFile, err := os.Create(filePaths[1])
	if err != nil {
		err = log.EnrichError(err, log.FileCreateError)
		err2 := sourceFile.Close()

		if err2 != nil {
			err = log.EnrichError(err, err2.Error())
		}

		log.Log(err.Error(), log.ErrorLevel)
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
		err := errors.New(log.NoSourceFile)
		log.Log(err.Error(), log.ErrorLevel)
		return nil, err
	}

	// If the output file name was not provided, default to "output"
	if *outputFile == "" {
		*outputFile = "output"
	}

	return []string{*inputFile, *outputFile}, nil
}
