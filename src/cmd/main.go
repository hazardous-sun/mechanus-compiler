package main

import (
	"flag"
	log "mechanus-compiler/src/error"
	"mechanus-compiler/src/models"
	"os"
)

func main() {
	// Collect source and output files
	sourceFile, outputFile, err := getFiles()
	errSalt := "(main) -> %w"

	if err != nil {
		err = log.FileErrorf(errSalt, err)
		log.LogError(err)
		os.Exit(1)
	}
	defer func() {
		if err := sourceFile.Close(); err != nil {
			err = log.FileErrorf(errSalt, err)
			log.LogError(err)
		}
		if err := outputFile.Close(); err != nil {
			err = log.FileErrorf(errSalt, err)
			log.LogError(err)
		}
	}()

	// Initialize the parser
	parser, err := models.NewParser(sourceFile, outputFile)

	if err != nil {
		err = log.EnrichError(err, "(main)")
		log.LogError(err)
		os.Exit(1)
	}

	// Start the syntax analysis
	if err = parser.Run(); err != nil {
		err = log.EnrichError(err, "(main)")
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
		err = log.FileErrorf("(getFiles) -> %w", err)
		log.LogError(err)
		return nil, nil, err
	}

	// Create output file
	outputFile, err := os.Create(filePaths[1])

	if err != nil {
		err = log.FileErrorf("(getFiles) -> %w", err)
		log.LogError(err)
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
		err := log.FileErrorf("(getFilePaths) -> %w", log.NoSourceFile)
		log.LogError(err)
		return nil, err
	}

	// If the output file name was not provided, default to "output"
	if *outputFile == "" {
		*outputFile = "output"
	}

	return []string{*inputFile, *outputFile}, nil
}
