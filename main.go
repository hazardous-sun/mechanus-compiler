package main

import (
	"errors"
	"fmt"
	custom_errors "mechanus-compiler/error"
	"mechanus-compiler/models"
	"os"
)

func main() {
	// Initialize input and output files insde Lexical
	lex := models.NewLexical(nil, nil)
	inputFile, err := os.Open("example.mecha")

	if err != nil {
		custom_errors.Log(custom_errors.FileOpenError, &err, custom_errors.ErrorLevel)
		return
	}
	lex.InputFile = inputFile
	defer lex.Close("input")

	outputFile, err := os.Create("output")

	if err != nil {
		custom_errors.Log(custom_errors.FileCreateError, &err, custom_errors.ErrorLevel)
		return
	}
	lex.OutputFile = outputFile
	defer lex.Close("output")

	// Start looking for tokens

	lex.movelookAhead()

	for lex.token != models.TInputEnd && lex.token != TLexError {
		lex.nextToken()
		lex.displayToken()
	}

	if lex.token == TLexError {
		err = errors.New(lex.errorMessage)
		custom_errors.Log(fmt.Sprintf("Lexical error: %s", lex.errorMessage), &err, custom_errors.ErrorLevel)
	} else {
		fmt.Println("Lexical analys completed with no errors")
	}

	lex.showTokens()
	lex.writeOutput()

	return lex.lexeme, err
}
