package models

import (
	"fmt"
	log "mechanus-compiler/src/error"
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
	lex, err := NewLexical(p.Source, p.Output)

	if err != nil {
		err = log.EnrichError(err, "Parser.Run()")
		log.Log(err.Error(), log.ErrorLevel)
		return
	}

	for lex.WIP() {
		err = lex.NextToken()

		if !comment(&lex) {
			lex.DisplayToken()
		}

		if err != nil {
			break
		}
	}

	if lex.token == TLexError {
		err = fmt.Errorf("%s -> %s", log.LexicalError, lex.errorMessage)
		log.Log(err.Error(), log.ErrorLevel)
	} else {
		log.Log(log.LexicalSuccess, log.SuccessLevel)
		err = lex.WriteOutput()
	}

	return
}

func comment(lex *Lexical) bool {
	return lex.token == TSingleLineComment ||
		lex.token == TOpenMultilineComment ||
		lex.token == TCloseMultilineComment ||
		lex.commentBlock == true
}
