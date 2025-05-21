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
	debug   bool
	source  *os.File
	output  *os.File
	lexical Lexical
}

// NewParser :
// Initializes a new instance of Parser with the provided input and output files.
//
// Fails if it is not possible to read the source file, or if the Lexical initialization fails.
func NewParser(source, output *os.File, debug bool) (Parser, error) {
	// Initialize Lexical
	lex, err := NewLexical(source, output, debug)

	// Check for errors during Lexical initialization
	if err != nil {
		err = log.SyntaxErrorf("NewParser", err)
		log.LogError(err)
		return Parser{}, err
	}

	return Parser{
		debug:   debug,
		source:  source,
		output:  output,
		lexical: lex,
	}, nil
}

// Run :
// Runs the syntax analysis.
func (p *Parser) Run() error {
	for p.lexical.WIP() {
		_, err := p.lexical.NextToken()

		if err != nil {
			break
		}
	}

	// Check if Lexical failed to reach EOF
	if err := p.lexical.Fail(); err != nil {
		err = log.SyntaxErrorf("Parser.Run", err)
		log.LogError(err)
		return err
	}

	log.LogSuccess(log.LexicalSuccess)

	if err := p.lexical.WriteOutput(); err != nil {
		err = log.SyntaxErrorf("Parser.Run", err)
		log.LogError(err)
		return err
	}

	return nil
}

func (p *Parser) parse() error {
	if err := p.g(); err != nil {
		err = log.SyntaxErrorf("Parser.parse", err)
		log.LogError(err)
		return err
	}

	return nil
}

// <G> ::= '{' <BODY> '}' <TEXT_WITHOUT_NUMBERS> 'Construct'
func (p *Parser) g() error {
	errSalt := "(Parser.g)"

	// Check for TConstruct
	if p.lexical.GetToken() != TConstruct {
		err := unexpectedLexeme(p, Construct)
		err = log.SyntaxErrorf(errSalt, err)
		log.LogError(err)
		return err
	}

	// Check if EOF was reached
	if !p.lexical.WIP() {
		err := fmt.Errorf(log.MissingConstructBody)
		err = log.SyntaxErrorf(errSalt, err)
		log.LogError(err)
		return err
	}

	_, err := p.lexical.NextToken()

	if err != nil {
		err = log.SyntaxErrorf(errSalt, err)
		log.LogError(err)
		return err
	}

	// Check for Construct name
	err = p.textWithoutNumbers()

	if err != nil {
		err = log.SyntaxErrorf(errSalt, err)
		log.LogError(err)
		return err
	}

	token, err := p.lexical.NextToken()

	if err != nil {
		err = log.SyntaxErrorf(errSalt, err)
		log.LogError(err)
		return err
	}

	// Check for "}"
	if token != TCloseBraces {
		err := unexpectedLexeme(p, "}")
		err = log.SyntaxErrorf(errSalt, err)
		log.LogError(err)
		return err
	}

	if _, err = p.lexical.NextToken(); err != nil {
		err = log.SyntaxErrorf(errSalt, err)
		log.LogError(err)
		return err
	}

	// Check for Construct body
	err = p.body()

	if err != nil {
		err = log.SyntaxErrorf(errSalt, err)
		log.LogError(err)
		return err
	}

	token, err = p.lexical.NextToken()

	if err != nil {
		err = log.SyntaxErrorf(errSalt, err)
		log.LogError(err)
		return err
	}

	// Check for "{"
	if token != TOpenBraces {
		err := unexpectedLexeme(p, "{")
		err = log.SyntaxErrorf(errSalt, err)
		log.LogError(err)
		return err
	}

	_, err = p.lexical.NextToken()

	// Check for EOF
	if err == nil {
		err = unexpectedLexeme(p, "EOF")
		err = log.SyntaxErrorf(errSalt, err)
		log.LogError(err)
		return err
	}

	return nil
}

// (([A-Z]|[a-z])+(_)*)+
func (p *Parser) textWithoutNumbers() error {
	if p.lexical.GetToken() != TId {
		err := unexpectedLexeme(p, "ID")
		err = log.SyntaxErrorf("Parser.textWithoutNumbers", err)
		log.LogError(err)
		return err
	}

	return nil
}

// <BODY> ::= '{' <CMDS> '}' '(' <PARAMETERS> ')' <TEXT_WITH_NUMBERS> 'Architect'
// <BODY> ::= '{' <CMDS> '}' <TYPE> '(' <PARAMETERS> ')' <TEXT_WITH_NUMBERS> 'Architect'
// <BODY> ::= <BODY_REST>
func (p *Parser) body() error {
	// Check for "TArchitect" OR bodyRest()
	return nil
}

// <BODY_REST> ::= '{' <CMDS> '}' '(' <PARAMETERS> ')' <TEXT_WITH_NUMBERS> 'Architect' <BODY_REST>
// <BODY_REST> ::= '{' <CMDS> '}' <TYPE> '(' <PARAMETERS> ')' <TEXT_WITH_NUMBERS> 'Architect' <BODY_REST>
// <BODY_REST> ::= ε
func (p *Parser) bodyRest() error {
	// TODO implement this logic
	return nil
}

// Helper methods
func unexpectedLexeme(p *Parser, expectedLexeme string) error {
	return log.SyntaxErrorf("Parser.unexpectedLexeme", fmt.Errorf("%s: (%s) expected lexeme '%s', found '%s'", log.SyntaxError, p.lexical.DisplayPos(), expectedLexeme, p.lexical.GetLexeme()))
}
