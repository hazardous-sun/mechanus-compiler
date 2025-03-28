package main

import (
	"errors"
	"fmt"
	customerrors "mechanus-compiler/error"
	"mechanus-compiler/models"
	"os"
)

func main() {
	// Initialize input and output files insde Lexical
	lex := models.NewLexical(nil, nil)
	inputFile, err := os.Open("example.mecha")

	if err != nil {
		customerrors.Log(customerrors.FileOpenError, &err, customerrors.ErrorLevel)
		return
	}
	lex.InputFile = inputFile
	defer lex.Close("input")

	outputFile, err := os.Create("output")

	if err != nil {
		customerrors.Log(customerrors.FileCreateError, &err, customerrors.ErrorLevel)
		return
	}
	lex.OutputFile = outputFile
	defer lex.Close("output")

	// Start looking for tokens
	lex.MovelookAhead()

	for lex.Token != models.TInputEnd && lex.Token != models.TLexError {
		lex.NextToken()
		lex.DisplayToken()
	}

	if lex.Token == models.TLexError {
		err = errors.New(lex.ErrorMessage)
		customerrors.Log(fmt.Sprintf("Lexical error: %s", lex.ErrorMessage), &err, customerrors.ErrorLevel)
	} else {
		fmt.Println("Lexical analysis completed with no errors")
		err = lex.WriteOutput()
	}

	lex.ShowTokens()
	return
}
