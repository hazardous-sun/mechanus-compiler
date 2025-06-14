package models

import (
	"fmt"
	log "mechanus-compiler/src/error"
	"os"
	"strings"
)

// Parser :
// This is the structure responsible for making the syntactical analysis of the source file. It checks for unrecognized
// syntaxes and, if it finds one, it returns an error code.
type Parser struct {
	debug           bool
	lexer           Lexer
	outputFile      *os.File
	token           int
	lexeme          string
	errorMessage    error
	recognizedRules strings.Builder
}

const (
	errExpectedCloseBraces      = "expected '}', got '%s'"
	errExpectedOpenBraces       = "expected '{', got '%s'"
	errExpectedOpenParenthesis  = "expected '(', got '%s'"
	errExpectedCloseParenthesis = "expected ')', got '%s'"
	errExpectedIdentifier       = "expected an identifier, got '%s'"
	errExpectedColon            = "expected ':', got '%s'"
)

// NewParser :
// Initializes a new Parser instance with the provided input and output files.
//
// Fails if it is not possible to initialize the lexer.
func NewParser(inputFile, outputFile *os.File, debug bool) (Parser, error) {
	// Initialize the Lexer
	lexer, err := NewLexer(inputFile, outputFile, debug)
	errSalt := "NewParser"

	if err != nil {
		err = log.SyntaxErrorf(errSalt, err)
		log.LogError(err)
		return Parser{}, err
	}

	// Initialize the structure
	parser := Parser{
		debug:        debug,
		lexer:        lexer,
		outputFile:   outputFile,
		token:        TNilValue,
		errorMessage: nil,
	}

	return parser, nil
}

// Run :
// Starts the syntactical analysis.
//
// Fails if the lexer fails or if a syntactical error is found.
func (parser *Parser) Run() error {
	errSalt := "Parser.Run"

	if err := parser.advanceToken(); err != nil {
		err = log.SyntaxErrorf(errSalt, err)
		log.LogError(err)
		return err
	}

	if err := parser.g(); err != nil {
		err = log.SyntaxErrorf(errSalt, err)
		log.LogError(err)
		return err
	}

	log.LogSuccess(log.SyntaxSuccess)
	return nil
}

// g :
// <G> ::= '{' <BODY> '}' <ID> 'Construct'
func (parser *Parser) g() error {
	errSalt := "Parser.g"
	parser.accumulateRule("<G> ::= '{' <BODY> '}' <ID> 'Construct'")

	// Expect 'Construct'
	if parser.token != TConstruct {
		return parser.handleSyntaxError(fmt.Errorf("expected 'Construct', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect <ID>
	if parser.token != TId {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedIdentifier, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect '}'
	if parser.token != TCloseBraces {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedCloseBraces, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect <BODY>
	if err := parser.body(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect '{'
	if parser.token != TOpenBraces {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedOpenBraces, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		if strings.Contains(err.Error(), log.EndOfFileReached) {
			return nil
		}
		return log.SyntaxErrorf(errSalt, err)
	}

	return nil
}

// <BODY> :
//
// <BODY> ::= <BODY_REST> '{' <CMDS> '}' '(' <PARAMETERS_DECL> ')' <ID> 'Architect'
// <BODY> ::= <BODY_REST> '{' <CMDS> '}' '(' ')' <ID> 'Architect'
// <BODY> ::= <BODY_REST> '{' <CMDS> '}' <TYPE> '(' ')' <ID> 'Architect'
// <BODY> ::= <BODY_REST> '{' <CMDS> '}' <TYPE> '(' <PARAMETERS_DECL> ')' <ID> 'Architect'
func (parser *Parser) body() error {
	errSalt := "Parser.body"
	parser.accumulateRule("<BODY> ::= <BODY_REST> '{' <CMDS> '}' '(' <PARAMETERS> ')' <ID> 'Architect' | ...")

	// 1. Expect 'Architect'
	if parser.token != TArchitect {
		return parser.handleSyntaxError(fmt.Errorf("expected 'Architect', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// 2. Expect <ID>
	if parser.token != TId {
		return parser.handleSyntaxError(fmt.Errorf("expected ID after Architect, got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// 3. Expect ')'
	if parser.token != TCloseParentheses {
		return parser.handleSyntaxError(fmt.Errorf("expected ')', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// 4. Optionally parse <PARAMETERS> (may be ε)
	if parser.token != TOpenParentheses {
		_ = parser.parametersDecl() // fail silently if no parametersDecl
	}

	// 5. Expect '('
	if parser.token != TOpenParentheses {
		return parser.handleSyntaxError(fmt.Errorf("expected '(', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// 6. Expect '}'
	if parser.token != TCloseBraces {
		return parser.handleSyntaxError(fmt.Errorf("expected '}', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// 7. Parse <CMDS>
	if err := parser.cmds(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// 8. Expect '{'
	if parser.token != TOpenBraces {
		return parser.handleSyntaxError(fmt.Errorf("expected '{', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		if strings.Contains(err.Error(), log.EndOfFileReached) {
			return nil
		}
		return log.SyntaxErrorf(errSalt, err)
	}

	// 9. Recursively parse any additional Architect bodies
	if err := parser.bodyRest(); err != nil {
		if strings.Contains(err.Error(), log.EndOfFileReached) {
			return nil
		}
		return log.SyntaxErrorf(errSalt, err)
	}

	return nil
}

// <BODY_REST> :
//
// <BODY_REST> ::= <BODY_REST> '{' <CMDS> '}' '(' <PARAMETERS_DECL> ')' <ID> 'Architect'
// <BODY_REST> ::= <BODY_REST> '{' <CMDS> '}' <TYPE> '(' <PARAMETERS_DECL> ')' <ID> 'Architect'
// <BODY_REST> ::= ε
func (parser *Parser) bodyRest() error {
	errSalt := "Parser.bodyRest"
	parser.accumulateRule("<BODY_REST> ::= <BODY_REST> '{' <CMDS> '}' '(' <PARAMETERS> ')' <ID> 'Architect' | ... | ε")

	// 1. Base case: ε
	if parser.token == TCloseBraces || parser.token == TInputEnd {
		parser.accumulateRule("<BODY_REST> ::= ε")
		return nil
	}

	// 2. Expect 'Architect'
	if parser.token != TArchitect {
		return parser.handleSyntaxError(fmt.Errorf("expected 'Architect', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// 3. Expect <ID>
	if parser.token != TId {
		return parser.handleSyntaxError(fmt.Errorf("expected ID after Architect, got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// 4. Expect ')'
	if parser.token != TCloseParentheses {
		return parser.handleSyntaxError(fmt.Errorf("expected ')', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// 5. Optionally parse <PARAMETERS>
	if parser.token != TOpenParentheses {
		_ = parser.parametersDecl() // fails silently if epsilon
	}

	// 6. Expect '('
	if parser.token != TOpenParentheses {
		return parser.handleSyntaxError(fmt.Errorf("expected '(', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// 7. Expect '}'
	if parser.token != TCloseBraces {
		return parser.handleSyntaxError(fmt.Errorf("expected '}', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// 8. Parse <CMDS>
	if err := parser.cmds(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// 9. Expect '{'
	if parser.token != TOpenBraces {
		return parser.handleSyntaxError(fmt.Errorf("expected '{', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// 10. Recurse to parse next body
	return parser.bodyRest()
}

// <TYPE> :
//
// <TYPE> ::= 'Nil'
// <TYPE> ::= 'Gear'
// <TYPE> ::= 'Tensor'
// <TYPE> ::= 'State'
// <TYPE> ::= 'Monodrone'
// <TYPE> ::= 'Omnidrone'
func (parser *Parser) typeToken() error {
	parser.accumulateRule("<TYPE> ::= 'Nil' | 'Gear' | 'Tensor' | 'State' | 'Monodrone' | 'Omnidrone'")
	if parser.token != TNil && parser.token != TGear && parser.token != TTensor &&
		parser.token != TState && parser.token != TMonodrone && parser.token != TOmnidrone {
		return parser.handleSyntaxError(fmt.Errorf("expected a Type keyword, got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err // Propagation of error
	}
	return nil
}

// <CMDS> :
// <CMDS> ::= <CMDS_REST> <CMD>
func (parser *Parser) cmds() error {
	parser.accumulateRule("<CMDS> ::= <CMDS_REST> <CMD>")

	for {
		// Skip any newlines
		for parser.token == TNewLine {
			parser.displayToken()
			if err := parser.advanceToken(); err != nil {
				return err
			}
		}

		// At this level, hitting '{' means the parser is done with <CMDS>
		if parser.token == TOpenBraces {
			break
		}

		// Attempt to parse one command
		if err := parser.cmd(); err != nil {
			return err
		}
	}

	return nil
}

// <CMDS_REST> :
//
// <CMDS_REST> ::= '\n' <CMDS>
// <CMDS_REST> ::= ε
func (parser *Parser) cmdsRest() error {
	parser.accumulateRule("<CMDS_REST> ::= '\\n' <CMDS> | ε")

	if parser.token == TNewLine {
		if err := parser.advanceToken(); err != nil {
			return err
		}
		return parser.cmds()
	}
	// epsilon
	parser.accumulateRule("<CMDS_REST> ::= ε")
	return nil
}

// <CMD> :
//
// <CMD> ::= <CMD_IF>
// <CMD> ::= <CMD_FOR>
// <CMD> ::= <CMD_DECLARATION>
// <CMD> ::= <CMD_ASSIGNMENT>
// <CMD> ::= <CMD_RECEIVE>
// <CMD> ::= <CMD_SEND>
// <CMD> ::= <CMD_INTEGRATE>
func (parser *Parser) cmd() error {
	errSalt := "Parser.cmd"
	parser.accumulateRule("<CMD> ::= <CMD_IF> | <CMD_FOR> | <CMD_DECLARATION> | <CMD_ASSIGNMENT> | <CMD_RECEIVE> | <CMD_SEND> | <CMD_INTEGRATE> | <CMD_CALL>")

	if err := parser.cmdIf(); err == nil {
		return nil
	}
	if err := parser.cmdFor(); err == nil {
		return nil
	}
	if err := parser.cmdDeclaration(); err == nil {
		return nil
	}
	if err := parser.cmdAssignment(); err == nil {
		return nil
	}
	if err := parser.cmdReceive(); err == nil {
		return nil
	}
	if err := parser.cmdSend(); err == nil {
		return nil
	}
	if err := parser.cmdIntegrate(); err == nil {
		return nil
	}

	if parser.token == TCloseParentheses {
		parser.accumulateRule("<CMD_CALL> ::= '(' <PARAMETERS_CALL> ')' <ID>")

		// Parse parameters
		if err := parser.parametersCall(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}

		return nil
	}

	// If no command matches, it's a syntax error
	return parser.handleSyntaxError(fmt.Errorf("unrecognized command starting with token %s", parser.lexeme))
}

// <CMD_IF> :
//
// <CMD_IF> ::= '{' <CMDS> '}' <CONDITION> 'if'
// <CMD_IF> ::= '{' <CMDS> '}' 'else' '{' <CMDS> '}' <CONDITION> 'if'
// <CMD_IF> ::= <CMD_ELIF> '{' <CMDS> '}' <CONDITION> 'if'
func (parser *Parser) cmdIf() error {
	errSalt := "Parser.cmdIf"
	parser.accumulateRule("<CMD_IF> ::= '{' <CMDS> '}' 'if' <CONDITION> | '{' <CMDS> '}' 'else' '{' <CMDS> '}' 'if' <CONDITION> | <CMD_ELIF> '{' <CMDS> '}' 'if' <CONDITION>")

	// Expect 'if'
	if parser.token != TIf {
		return parser.handleSyntaxError(fmt.Errorf("expected 'if', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect <CONDITION>
	if err := parser.condition(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect '}'
	if parser.token != TCloseBraces {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedCloseBraces, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect <CMDS>
	if err := parser.cmds(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect '{'
	if parser.token != TOpenBraces {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedOpenBraces, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Check for 'elif'
	for parser.token == TElif {
		if err := parser.cmdElif(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
	}

	// Check for 'else'
	if parser.token == TElse {
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
		// Expect '}' after 'else' block
		if parser.token != TCloseBraces {
			return parser.handleSyntaxError(fmt.Errorf(errExpectedCloseBraces, parser.lexeme))
		}
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
		// Expect <CMDS> for 'else' block
		if err := parser.cmds(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
		// Expect '{' for 'else' block
		if parser.token != TOpenBraces {
			return parser.handleSyntaxError(fmt.Errorf(errExpectedOpenBraces, parser.lexeme))
		}
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
	}

	return nil
}

// <CMD_ELIF> :
//
// <CMD_ELIF> ::= '{' <CMDS> '}' <CONDITION> 'elif'
// <CMD_ELIF> ::= <CMD_ELIF_REST>
func (parser *Parser) cmdElif() error {
	errSalt := "Parser.cmdElif"
	parser.accumulateRule("<CMD_ELIF> ::= '{' <CMDS> '}' 'elif' <CONDITION> | <CMD_ELIF_REST>")

	// If the current token is TElif, it's a direct elif. Otherwise, it must be CMD_ELIF_REST.
	if parser.token != TElif {
		return parser.cmdElifRest()
	}

	// Expect 'elif'
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect <CONDITION>
	if err := parser.condition(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect '}'
	if parser.token != TCloseBraces {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedCloseBraces, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect <CMDS>
	if err := parser.cmds(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect '{'
	if parser.token != TOpenBraces {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedOpenBraces, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	return nil
}

// <CMD_ELIF_REST> :
//
// <CMD_ELIF_REST> ::= <CMD_ELIF_REST> '{' <CMDS> '}' <CONDITION> 'elif'
// <CMD_ELIF_REST> ::= '{' <CMDS> '}' 'else' <CMD_ELIF_REST> '{' <CMDS> '}' <CONDITION> 'elif'
// <CMD_ELIF_REST> ::= ε
func (parser *Parser) cmdElifRest() error {
	errSalt := "Parser.cmdElifRest"
	parser.accumulateRule("<CMD_ELIF_REST> ::= '{' <CMDS> '}' 'elif' <CONDITION> <CMD_ELIF_REST> | '{' <CMDS> '}' 'else' '{' <CMDS> '}' 'elif' <CONDITION> <CMD_ELIF_REST> | ε")

	// If the next token is not 'elif', it's epsilon.
	if parser.token != TElif {
		parser.accumulateRule("<CMD_ELIF_REST> ::= ε")
		return nil
	}

	// Expect 'elif'
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect <CONDITION>
	if err := parser.condition(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect '}'
	if parser.token != TCloseBraces {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedCloseBraces, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect <CMDS>
	if err := parser.cmds(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect '{'
	if parser.token != TOpenBraces {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedOpenBraces, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Check for 'else'
	if parser.token == TElse {
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
		// Expect '}' after 'else' block
		if parser.token != TCloseBraces {
			return parser.handleSyntaxError(fmt.Errorf(errExpectedCloseBraces, parser.lexeme))
		}
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
		// Expect <CMDS> for 'else' block
		if err := parser.cmds(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
		// Expect '{' for 'else' block
		if parser.token != TOpenBraces {
			return parser.handleSyntaxError(fmt.Errorf(errExpectedOpenBraces, parser.lexeme))
		}
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
	}

	// Recursive call for CMD_ELIF_REST
	if err := parser.cmdElifRest(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	return nil
}

// <CMD_FOR> :
// <CMD_FOR> ::= '{' <CMDS> '}' <CONDITION> 'for'
func (parser *Parser) cmdFor() error {
	errSalt := "Parser.cmdFor"
	parser.accumulateRule("<CMD_FOR> ::= '{' <CMDS> '}' <CONDITION> 'for'")

	// Expect 'for'
	if parser.token != TFor {
		return parser.handleSyntaxError(fmt.Errorf("expected 'for', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect <CONDITION>
	if err := parser.condition(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect '}'
	if parser.token != TCloseBraces {
		return parser.handleSyntaxError(fmt.Errorf("expected '}', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect <CMDS>
	if err := parser.cmds(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect '{'
	if parser.token != TOpenBraces {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedOpenBraces, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	return nil
}

// <CMD_INTEGRATE> :
//
// <CMD_INTEGRATE> ::= <E> 'Integrate'
func (parser *Parser) cmdIntegrate() error {
	errSalt := "Parser.cmdIntegrate"
	parser.accumulateRule("<CMD_INTEGRATE> ::= <E> 'Integrate'")

	// Expect 'Integrate'
	if parser.token != TIntegrate {
		return parser.handleSyntaxError(fmt.Errorf("expected 'Integrate', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect <E>
	if err := parser.e(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	return nil
}

// <CMD_DECLARATION> :
//
// <CMD_DECLARATION> ::= <E> '=:' <TYPE> ':' <VAR>
func (parser *Parser) cmdDeclaration() error {
	errSalt := "Parser.cmdDeclaration"
	parser.accumulateRule("<CMD_DECLARATION> ::= <E> '=:' <TYPE> ':' <VAR>")

	// Expect <VAR>
	if err := parser.varToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect ':'
	if parser.token != TColon {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedColon, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect <TYPE>
	if err := parser.typeToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect '=:'
	if parser.token != TDeclarationOperator {
		return parser.handleSyntaxError(fmt.Errorf("expected '=:', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect <E>
	if err := parser.e(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	return nil
}

// <CMD_ASSIGNMENT> :
//
// <CMD_ASSIGNMENT> ::= <E> '=' <VAR>
func (parser *Parser) cmdAssignment() error {
	errSalt := "Parser.cmdAssignment"
	parser.accumulateRule("<CMD_ASSIGNMENT> ::= <E> '=' <VAR>")

	// Expect <VAR>
	if err := parser.varToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect '='
	if parser.token != TAttributionOperator {
		return parser.handleSyntaxError(fmt.Errorf("expected '=', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect <E>
	if err := parser.e(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	return nil
}

// <CMD_RECEIVE> :
//
// <CMD_RECEIVE> ::= '(' <VAR> ')' 'Receive'
func (parser *Parser) cmdReceive() error {
	errSalt := "Parser.cmdReceive"
	parser.accumulateRule("<CMD_RECEIVE> ::= '(' <VAR> ')' 'Receive'")

	// Expect 'Receive'
	if parser.token != TReceive {
		return parser.handleSyntaxError(fmt.Errorf("expected 'Receive', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect ')'
	if parser.token != TCloseParentheses {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedCloseParenthesis, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect <VAR>
	if err := parser.varToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect '('
	if parser.token != TOpenParentheses {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedOpenParenthesis, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	return nil
}

// <CMD_SEND> :
//
// <CMD_SEND> ::= '(' <E> ')' 'Send'
func (parser *Parser) cmdSend() error {
	parser.accumulateRule("<CMD_SEND> ::= '(' <E> ')' 'Send'")

	// Expect TSend (first, since lexing is bottom-up, right-to-left)
	if parser.token != TSend {
		return parser.handleSyntaxError(fmt.Errorf("expected 'Send', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect TCloseParentheses
	if parser.token != TCloseParentheses {
		return parser.handleSyntaxError(fmt.Errorf("expected ')', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Parse <E>
	if err := parser.e(); err != nil {
		return err
	}

	// Expect TOpenParentheses
	if parser.token != TOpenParentheses {
		return parser.handleSyntaxError(fmt.Errorf("expected '(', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	return nil
}

// <CONDITION> :
//
// <CONDITION> ::= <E> '>' <E>
// <CONDITION> ::= <E> '>=' <E>
// <CONDITION> ::= <E> '<>' <E>
// <CONDITION> ::= <E> '<=' <E>
// <CONDITION> ::= <E> '<' <E>
// <CONDITION> ::= <E> '==' <E>
func (parser *Parser) condition() error {
	errSalt := "Parser.condition"
	parser.accumulateRule("<CONDITION> ::= <E> '>' <E> | <E> '>=' <E> | <E> '<>' <E> | <E> '<=' <E> | <E> '<' <E> | <E> '==' <E>")

	// All conditions are of the form <E> OPERATOR <E>
	// Parse the second <E> (rightmost) first
	if err := parser.e(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect a comparison operator
	if parser.token != TGreaterThanOperator &&
		parser.token != TGreaterEqualOperator &&
		parser.token != TLessThanOperator &&
		parser.token != TLessEqualOperator &&
		parser.token != TNotEqualOperator &&
		parser.token != TEqualOperator {
		return parser.handleSyntaxError(fmt.Errorf("expected a comparison operator, got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Parse the first <E> (leftmost)
	if err := parser.e(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	return nil
}

// <E> :
// <E> ::= <E_REST> <T>
func (parser *Parser) e() error {
	errSalt := "Parser.e"
	parser.accumulateRule("<E> ::= <T> <E_REST>")

	if err := parser.t(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}
	return parser.eRest()
}

// eRest :
// <E_REST> ::= <E_REST> '+' <T>
// <E_REST> ::= <E_REST> '-' <T>
// <E_REST> ::= ε
func (parser *Parser) eRest() error {
	errSalt := "Parser.eRest"

	switch parser.token {
	case TAdditionOperator, TSubtractionOperator:
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
		if err := parser.t(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
		return parser.eRest() // recursive continuation
	default:
		// ε-production matched — stop parsing this rule
		parser.accumulateRule("<E_REST> ::= ε")
		return nil
	}
}

// <T> :
func (parser *Parser) t() error {
	errSalt := "Parser.t"
	parser.accumulateRule("<T> ::= <F> <T_REST>")

	if err := parser.f(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	return parser.tRest()
}

// <T_REST>
func (parser *Parser) tRest() error {
	errSalt := "Parser.tRest"

	switch parser.token {
	case TMultiplicationOperator, TDivisionOperator, TModuleOperator:
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
		if err := parser.f(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
		return parser.tRest() // recursive continuation
	default:
		// ε-production matched — stop
		parser.accumulateRule("<T_REST> ::= ε")
		return nil
	}
}

// <F>
func (parser *Parser) f() error {
	errSalt := "Parser.f"
	parser.accumulateRule("<F> ::= -<F> | <X>")

	if parser.token == TSubtractionOperator {
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
		if err := parser.f(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
	} else {
		if err := parser.x(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
	}
	return nil
}

// <X> :
func (parser *Parser) x() error {
	errSalt := "Parser.x"
	parser.accumulateRule("<X> ::= '(' <E> ')' | [0-9]+('.'[0-9]+) | <STRING> | <NIL> | <VAR> | '(' <PARAMETERS_CALL> ')' <ID>")

	switch parser.token {

	// Case: STRING literal
	case TDoubleQuote:
		parser.displayToken()
		return parser.advanceToken()

	// Case: NIL
	case TNil:
		parser.displayToken()
		return parser.advanceToken()

	// Case: numeric literal (integer or float)
	case TGear, TTensor:
		parser.displayToken()
		return parser.advanceToken()

	// Case: identifier (variable or function call)
	case TId:
		parser.displayToken()
		return parser.advanceToken()

	// Case: open parentheses — could be (E) or (PARAMS_CALL) ID
	case TCloseParentheses:
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}

		// Try to parse parametersCall (supports multiple expressions)
		if err := parser.parametersCall(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}

		// Expect matching opening parenthesis
		if parser.token != TOpenParentheses {
			return parser.handleSyntaxError(fmt.Errorf("expected '(', got %s", parser.lexeme))
		}
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}

		return nil
	}

	// If no valid rule matches, return error
	return parser.handleSyntaxError(fmt.Errorf("unexpected token in <X>: %s", parser.lexeme))
}

func (parser *Parser) nilToken() error {
	parser.accumulateRule("<NIL> :: 'Nil'")

	if parser.token != TNil {
		return parser.handleSyntaxError(fmt.Errorf("expected 'Nil', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}
	return nil
}

// <STRING> :
//
// <STRING> ::= '"' <TEXT_WITH_NUMBERS> '"'
func (parser *Parser) stringToken() error {
	parser.accumulateRule("<STRING> ::= '\"' <TEXT_WITH_NUMBERS> '\"'")

	// The lexer identifies the entire string literal (including quotes) as TDoubleQuote.
	// So, the parser just needs to consume the TDoubleQuote
	if parser.token != TDoubleQuote {
		return parser.handleSyntaxError(fmt.Errorf("expected a string literal, got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}
	return nil
}

// <VAR> :
//
// <VAR> ::= <ID>
func (parser *Parser) varToken() error {
	errSalt := "Parser.varToken"
	parser.accumulateRule("<VAR> ::= <ID>")

	if err := parser.id(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}
	return nil
}

// <ID> :
//
// <ID> ::= (([A-Z]|[a-z])+(_|[0-9])*)+
func (parser *Parser) id() error {
	parser.accumulateRule("<ID> ::= (([A-Z]|[a-z])+(_|[0-9])*)+")
	if parser.token != TId {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedIdentifier, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err // Propagation of error
	}
	return nil
}

// <PARAMETERS_DECL> :
//
// <PARAMETERS_CALL> ::= <E>
// <PARAMETERS_CALL> ::= <EXTRA_PARAMETERS_CALL> <E>
// <EXTRA_PARAMETERS_CALL> ::= <E> ','
// <EXTRA_PARAMETERS_CALL> ::= <EXTRA_PARAMETERS_CALL> <E> ','
func (parser *Parser) parametersDecl() error {
	errSalt := "Parser.parametersDecl"
	parser.accumulateRule("<PARAMETERS> ::= <EXTRA_PARAMETERS> <TYPE> ':' <ID> | <TYPE> ':' <ID>")

	// Expect ID (rightmost identifier in the parameter list)
	if parser.token != TId {
		return parser.handleSyntaxError(fmt.Errorf("expected parameter ID, got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect ':'
	if parser.token != TColon {
		return parser.handleSyntaxError(fmt.Errorf("expected ':', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect <TYPE>
	if err := parser.typeToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Loop to check for extra parametersDecl (reverse order)
	for parser.token == TComma {
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}

		// Expect ID
		if parser.token != TId {
			return parser.handleSyntaxError(fmt.Errorf("expected parameter ID, got %s", parser.lexeme))
		}
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}

		// Expect ':'
		if parser.token != TColon {
			return parser.handleSyntaxError(fmt.Errorf("expected ':', got %s", parser.lexeme))
		}
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}

		// Expect <TYPE>
		if err := parser.typeToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
	}

	return nil
}

// <PARAMETERS_CALL>
//
// <PARAMETERS_CALL> ::= <E>
// <PARAMETERS_CALL> ::= <EXTRA_PARAMETERS_CALL> <E>
// <EXTRA_PARAMETERS_CALL> ::= <E> ',' | <EXTRA_PARAMETERS_CALL> <E> ','
func (parser *Parser) parametersCall() error {
	errSalt := "Parser.parametersCall"
	parser.accumulateRule("<PARAMETERS_CALL> ::= <EXTRA_PARAMETERS_CALL> <E> | <E>")

	// Parse rightmost expression (last param)
	if err := parser.e(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Repeatedly handle comma-separated expressions
	for parser.token == TComma {
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}

		if err := parser.e(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
	}

	return nil
}

// <EXTRA_PARAMETERS> :
//
// <EXTRA_PARAMETERS> ::= ',' <ID> ':' <TYPE>
// <EXTRA_PARAMETERS> ::= ',' <ID> ':' <TYPE> ',' <PARAMETERS>
func (parser *Parser) extraParameters() error {
	errSalt := "Parser.extraParameters"
	parser.accumulateRule("<EXTRA_PARAMETERS> ::= ',' <ID> ':' <TYPE> | ',' <ID> ':' <TYPE> ',' <PARAMETERS>")

	// If there's a comma, it's a recursive call or a single extra parameter.
	if parser.token == TComma {
		// Consume the ','
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}

		// Check for <TYPE>
		if err := parser.typeToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}

		// Then ':'
		if parser.token != TColon {
			return parser.handleSyntaxError(fmt.Errorf(errExpectedColon, parser.lexeme))
		}
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}

		// Then <ID>
		if err := parser.id(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}

		// After matching `ID : TYPE`, if the next token is also a `TComma`, then recursively call `extraParameters`.
		if parser.token == TComma {
			// Consume the comma for the recursive call
			parser.displayToken()
			if err := parser.advanceToken(); err != nil {
				return log.SyntaxErrorf(errSalt, err)
			}
			if err := parser.parametersDecl(); err != nil {
				return log.SyntaxErrorf(errSalt, err)
			}
		}

	} else {
		return parser.handleSyntaxError(fmt.Errorf("expected ',', got %s for extraParameters", parser.lexeme))
	}
	return nil
}

// advanceToken :
// Advances the lexer to the next token and updates the parser's state.
//
// Fails if the lexer fails to get the next token.
func (parser *Parser) advanceToken() error {
	errSalt := "Parser.advanceToken"
	if parser.debug {
		log.LogDebug("Advancing token...")
	}

	token, err := parser.lexer.NextToken()
	if err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	parser.token = token
	parser.lexeme = parser.lexer.GetLexeme()

	return nil
}

// displayToken :
// Displays the current token and lexeme.
func (parser *Parser) displayToken() {
	if parser.debug {
		parser.lexer.DisplayToken()
	}
}

// handleSyntaxError :
// Records a syntax error and sets the error message.
//
// Returns a new error of type ErrSyntax.
func (parser *Parser) handleSyntaxError(err error) error {
	if parser.errorMessage == nil {
		parser.errorMessage = log.SyntaxErrorf(log.SyntaxError,
			fmt.Errorf("%s at %s", err.Error(), parser.lexer.DisplayPos()))
		//log.LogError(parser.errorMessage)
	}
	return parser.errorMessage
}

// accumulateRule :
// Accumulates the recognized grammar rule for output.
func (parser *Parser) accumulateRule(rule string) {
	if parser.debug {
		log.LogDebug(fmt.Sprintf("Recognized rule: %s\n", rule))
		parser.recognizedRules.WriteString(fmt.Sprintf("Recognized rule: %s\n", rule))
	}
}

// ShowRecognizedRules :
// Displays all recognized grammar rules.
func (parser *Parser) ShowRecognizedRules() {
	if parser.debug {
		log.LogDebug("Recognized Grammar Rules:")
	}
	fmt.Println(parser.recognizedRules.String())
}

// Fail :
// Checks if the Parser failed during execution.
func (parser *Parser) Fail() error {
	return parser.errorMessage
}

// Close :
// Closes the input and output files used by the parser's lexer.
func (parser *Parser) Close() {
	parser.lexer.Close("input")
	parser.lexer.Close("output")
}

// WriteOutput :
// Writes the lexer's identified tokens to the output file.
func (parser *Parser) WriteOutput() error {
	return parser.lexer.WriteOutput()
}
