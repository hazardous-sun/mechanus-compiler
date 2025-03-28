package main

import (
	"errors"
	"mechanus-compiler/src/error"
	"mechanus-compiler/src/models"
	"os"
)

// TODO check for Active and Inative values for State variables
// TODO check why Send function is not being correctly shown
// TODO check why the Construct name is not being collected

func main() {
	// Initialize input and output files insde Lexical
	lex := models.NewLexical(nil, nil)
	inputFile, err := os.Open("examples/example1.mecha")

	if err != nil {
		custom_errors.Log(custom_errors.FileOpenError, &err, custom_errors.ErrorLevel)
		return
	}
	lex.InputFile = inputFile
	defer lex.Close("input")

	outputFile, err := os.Create("output.txt")

	if err != nil {
		custom_errors.Log(custom_errors.FileCreateError, &err, custom_errors.ErrorLevel)
		return
	}
	lex.OutputFile = outputFile
	defer lex.Close("output")

	// Start looking for tokens
	err = lex.ReadLines()

	if err != nil {
		custom_errors.Log(custom_errors.EmptyFile, &err, custom_errors.ErrorLevel)
		return
	}

	err = lex.MoveLookAhead()

	if err != nil {
		return
	}

	for lex.Token != models.TInputEnd && lex.Token != models.TLexError {
		err = lex.NextToken()

		if !comment(&lex) {
			lex.DisplayToken()
		}

		if err != nil {
			break
		}
	}

	if lex.Token == models.TLexError {
		err = errors.New(lex.ErrorMessage)
		custom_errors.Log(custom_errors.LexicalError, &err, custom_errors.ErrorLevel)
	} else {
		custom_errors.Log(custom_errors.LexicalSuccess, nil, custom_errors.SuccessLevel)
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
