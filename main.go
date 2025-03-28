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
	err = lex.ReadLines()

	if err != nil {
		customerrors.Log(customerrors.EmptyFile, &err, customerrors.ErrorLevel)
		return
	}

	err = lex.MovelookAhead()

	if err != nil {
		return
	}

	for lex.Token != models.TInputEnd && lex.Token != models.TLexError {
		err = lex.NextToken()
		if err != nil {
			break
		}
		if !comment(&lex) {
			lex.DisplayToken()
		}
	}

	if lex.Token == models.TLexError {
		err = errors.New(lex.ErrorMessage)
		customerrors.Log(fmt.Sprintf("Lexical error: %s", lex.ErrorMessage), &err, customerrors.ErrorLevel)
	} else {
		fmt.Println("Lexical analysis completed with no errors")
		err = lex.WriteOutput()
	}

	//lex.ShowTokens()
	return
}

func comment(lex *models.Lexical) bool {
	return lex.Token == models.TSingleLineComment ||
		lex.Token == models.TOpenMultilineComment ||
		lex.Token == models.TCloseMultilineComment ||
		lex.CommentBlock == true
}
