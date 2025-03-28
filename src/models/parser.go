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

	for lex.Token != TInputEnd && lex.Token != TLexError {
		err = lex.NextToken()

		if !comment(&lex) {
			lex.DisplayToken()
		}

		if err != nil {
			break
		}
	}

	if lex.Token == TLexError {
		err = errors.New(lex.ErrorMessage)
		custom_errors.Log(custom_errors.LexicalError, &err, custom_errors.ErrorLevel)
	} else {
		custom_errors.Log(custom_errors.LexicalSuccess, nil, custom_errors.SuccessLevel)
		err = lex.WriteOutput()
	}

	//lex.ShowTokens()
	return
}

func comment(lex *Lexical) bool {
	return lex.Token == TSingleLineComment ||
		lex.Token == TOpenMultilineComment ||
		lex.Token == TCloseMultilineComment ||
		lex.CommentBlock == true
}
