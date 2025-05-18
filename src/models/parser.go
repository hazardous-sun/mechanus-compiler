package models

import (
	"fmt"
	log "mechanus-compiler/src/error"
	"os"
)

// Parser :
// This is the structure responsible for analyzing the syntax of the source file. It interacts directly with the lexical
// analyzer and checks recursively for invalid tokens.
type Parser struct {
	source  *os.File
	output  *os.File
	lexical Lexical
}

// NewParser :
// Initializes a new instance of Parser with the provided input and output files.
//
// Fails if it is not possible to read the source file, or if the Lexical initialization fails.
func NewParser(source, output *os.File) (Parser, error) {
	// Initialize Lexical
	lex, err := NewLexical(source, output)

	// Check for errors during Lexical initialization
	if err != nil {
		err = log.EnrichError(err, "NewParser()")
		log.Log(err.Error(), log.ErrorLevel)
		return Parser{}, err
	}

	return Parser{
		source:  source,
		output:  output,
		lexical: lex,
	}, nil
}

func (p *Parser) Run() error {
	for p.lexical.WIP() {
		_, err := p.lexical.NextToken()

		if err != nil {
			break
		}
	}

	// Check if Lexical failed to reach EOF
	if p.lexical.Fail() {
		err := fmt.Errorf("%s -> %s", log.LexicalError, p.lexical.errorMessage)
		log.Log(err.Error(), log.ErrorLevel)
	}

	log.Log(log.LexicalSuccess, log.SuccessLevel)
	err := p.lexical.WriteOutput()

	if err != nil {
		err = log.EnrichError(err, "Parser.Run()")
		log.Log(err.Error(), log.ErrorLevel)
		return err
	}

	return nil
}
