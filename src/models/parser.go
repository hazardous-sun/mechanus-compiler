package models

import (
	"errors"
	custom_errors "mechanus-compiler/src/error"
	"os"
)

// Parser :
// This is the structure responsible for analyzing the syntax of the Source file. It interacts directly with the lexical
// analyzer and checks recursively for invalid tokens.
type Parser struct {
	Source       *os.File
	Output       *os.File
	SyntaxError  bool
	LexicalError bool
}

func (p *Parser) Run() {
	// Initialize input and Output files insde Lexical
	lex := NewLexical(p.Source, p.Output)

	// Start looking for tokens
	err := lex.ReadLines()

	if err != nil {
		custom_errors.Log(custom_errors.EmptyFile, &err, custom_errors.ErrorLevel)
		return
	}

	err = lex.MoveLookAhead()

	if err != nil {
		return
	}

	for lex.token != TInputEnd && lex.token != TLexError {
		err = lex.NextToken()

		if !comment(&lex) {
			lex.DisplayToken()
		}

		if err != nil {
			break
		}
	}

	if lex.token == TLexError {
		err = errors.New(lex.errorMessage)
		custom_errors.Log(custom_errors.LexicalError, &err, custom_errors.ErrorLevel)
	} else {
		custom_errors.Log(custom_errors.LexicalSuccess, nil, custom_errors.SuccessLevel)
		err = lex.WriteOutput()
	}

	//lex.ShowTokens()
	return
}

func comment(lex *Lexical) bool {
	return lex.token == TSingleLineComment ||
		lex.token == TOpenMultilineComment ||
		lex.token == TCloseMultilineComment ||
		lex.commentBlock == true
}
