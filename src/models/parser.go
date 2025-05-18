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
	// Initialize Lexical
	lex, err := NewLexical(p.Source, p.Output)

	// Check for errors during Lexical initialization
	if err != nil {
		err = log.EnrichError(err, "Parser.Run()")
		log.Log(err.Error(), log.ErrorLevel)
		return
	}

	for lex.WIP() {
		lexeme, err := lex.NextToken()

		log.Log(lexeme, log.WarningLevel)

		if err != nil {
			break
		}
	}

	// Check if Lexical failed to reach EOF
	if lex.Fail() {
		err = fmt.Errorf("%s -> %s", log.LexicalError, lex.errorMessage)
		log.Log(err.Error(), log.ErrorLevel)
	}

	log.Log(log.LexicalSuccess, log.SuccessLevel)
	err = lex.WriteOutput()

	return
}
