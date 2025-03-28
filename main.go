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
	err = lex.MovelookAhead()

	if err != nil {
		custom_errors.Log(custom_errors.EmptyFile, &err, custom_errors.ErrorLevel)
		return
	}

	for lex.Token != models.TInputEnd && lex.Token != models.TLexError {
		lex.NextToken()
		lex.DisplayToken()
	}

	if lex.Token == models.TLexError {
		err = errors.New(lex.ErrorMessage)
		custom_errors.Log(fmt.Sprintf("Lexical error: %s", lex.ErrorMessage), &err, custom_errors.ErrorLevel)
	} else {
		fmt.Println("Lexical analys completed with no errors")
		err = lex.WriteOutput()
	}

	lex.ShowTokens()
	return
}
