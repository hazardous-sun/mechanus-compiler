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
	errExpectedCloseBraces      = "expected '}', got %s"
	errExpectedOpenBraces       = "expected '{', got %s"
	errExpectedOpenParenthesis  = "expected '(', got %s"
	errExpectedCloseParenthesis = "expected ')', got %s"
	errExpectedIdentifier       = "expected an identifier, got %s"
	errExpectedColon            = "expected ':', got %s"
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
		return log.SyntaxErrorf(errSalt, err)
	}

	return nil
}

// body :
// <BODY> ::= '{' <CMDS> '}' '(' <PARAMETERS> ')' <ID> 'Architect'
// <BODY> ::= '{' <CMDS> '}' <TYPE> '(' <PARAMETERS> ')' <ID> 'Architect'
// <BODY> ::= <BODY_REST>
func (parser *Parser) body() error {
	errSalt := "Parser.body"
	parser.accumulateRule("<BODY> ::= '{' <CMDS> '}' '(' <PARAMETERS> ')' <ID> 'Architect' | '{' <CMDS> '}' <TYPE> '(' <PARAMETERS> ')' <ID> 'Architect' | <BODY_REST>")

	// The rule is written bottom-up, right-to-left. So we need to look for 'Architect' first.
	// Since BODY_REST is epsilon, we should check for the fixed part first.

	// Check for BODY_REST (epsilon)
	if parser.token == TCloseBraces || parser.token == TId { // Assuming these tokens indicate the end of a BODY or the start of a new one.
		parser.accumulateRule("<BODY> ::= ε (from <BODY_REST>)")
		return parser.bodyRest()
	}

	// It means it must start with 'Architect'
	if parser.token != TArchitect {
		return parser.handleSyntaxError(fmt.Errorf("expected 'Architect' or start of a new body, got %s", parser.lexeme))
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

	// Expect ')'
	if parser.token != TCloseParentheses {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedCloseBraces, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Expect <PARAMETERS> (not every function will have parameters, so the call may fail, and that's okay)
	_ = parser.parameters()

	// Expect '('
	if parser.token != TOpenParentheses {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedOpenParenthesis, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Check if there is a TYPE before the opening parenthesis (second <BODY> rule)
	if parser.token == TTypeName {
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
		parser.accumulateRule("<BODY> ::= '{' <CMDS> '}' <TYPE> '(' <PARAMETERS> ')' <ID> 'Architect'")
	} else {
		parser.accumulateRule("<BODY> ::= '{' <CMDS> '}' '(' <PARAMETERS> ')' <ID> 'Architect'")
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

// bodyRest :
// <BODY_REST> ::= '{' <CMDS> '}' '(' <PARAMETERS> ')' <ID> 'Architect' <BODY_REST>
// <BODY_REST> ::= '{' <CMDS> '}' <TYPE> '(' <PARAMETERS> ')' <ID> 'Architect' <BODY_REST>
// <BODY_REST> ::= ε
func (parser *Parser) bodyRest() error {
	errSalt := "Parser.bodyRest"
	parser.accumulateRule("<BODY_REST> ::= '{' <CMDS> '}' '(' <PARAMETERS> ')' <ID> 'Architect' <BODY_REST> | '{' <CMDS> '}' <TYPE> '(' <PARAMETERS> ')' <ID> 'Architect' <BODY_REST> | ε")

	// If the next token is 'Architect', we have a BODY_REST, otherwise it's epsilon.
	if parser.token == TArchitect {
		if err := parser.body(); err != nil { // Re-use the body method, as the structure is similar
			return log.SyntaxErrorf(errSalt, err)
		}
		// After parsing a body, recursively call bodyRest to handle multiple bodies.
		if err := parser.bodyRest(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
	} else {
		parser.accumulateRule("<BODY_REST> ::= ε")
	}
	return nil
}

// typeToken :
// <TYPE> ::= 'Nil' | 'Gear' | 'Tensor' | 'State' | 'Monodrone' | 'Omnidrone'
func (parser *Parser) typeToken() error {
	parser.accumulateRule("<TYPE> ::= 'Nil' | 'Gear' | 'Tensor' | 'State' | 'Monodrone' | 'Omnidrone'")
	if parser.token != TNilValue && parser.token != TGear && parser.token != TTensor &&
		parser.token != TState && parser.token != TMonodrone && parser.token != TOmnidrone {
		return parser.handleSyntaxError(fmt.Errorf("expected a Type keyword, got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err // Propagation of error
	}
	return nil
}

// cmds :
// <CMDS> ::= <CMD> <CMDS_REST>
func (parser *Parser) cmds() error {
	errSalt := "Parser.cmds"
	parser.accumulateRule("<CMDS> ::= <CMD> <CMDS_REST>")

	if err := parser.cmd(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}
	if err := parser.cmdsRest(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}
	return nil
}

// cmdsRest :
// <CMDS_REST> ::= '\n' <CMDS> | ε
func (parser *Parser) cmdsRest() error {
	errSalt := "Parser.cmdsRest"
	parser.accumulateRule("<CMDS_REST> ::= '\\n' <CMDS> | ε")

	if parser.token == TNewLine {
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
		if err := parser.cmds(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
	} else {
		parser.accumulateRule("<CMDS_REST> ::= ε")
	}
	return nil
}

// cmd :
// <CMD> ::= <CMD_IF> | <CMD_FOR> | <CMD_DECLARATION> | <CMD_ASSIGNMENT> | <CMD_RECEIVE> | <CMD_SEND>
func (parser *Parser) cmd() error {
	parser.accumulateRule("<CMD> ::= <CMD_IF> | <CMD_FOR> | <CMD_DECLARATION> | <CMD_ASSIGNMENT> | <CMD_RECEIVE> | <CMD_SEND>")

	switch parser.token {
	case TIf: // CMD_IF - looking for 'if' or 'else' or 'elif'
		return parser.cmdIf()
	case TFor: // CMD_FOR - looking for 'for'
		return parser.cmdFor()
	case TReceive: // CMD_RECEIVE - looking for 'Receive'
		return parser.cmdReceive()
	case TSend: // CMD_SEND - looking for 'Send'
		return parser.cmdSend()
	}

	// <CMD_DECLARATION>
	if err := parser.cmdDeclaration(); err == nil {
		return nil
	}

	// <CMD_ASSIGNMENT>
	if err := parser.cmdAssignment(); err == nil {
		return nil
	}

	return parser.handleSyntaxError(fmt.Errorf("expected a command, got %s", parser.lexeme))
}

// cmdIf :
// <CMD_IF> ::= '{' <CMDS> '}' 'if' <CONDITION>
// <CMD_IF> ::= '{' <CMDS> '}' 'else' '{' <CMDS> '}' 'if' <CONDITION>
// <CMD_IF> ::= <CMD_ELIF> '{' <CMDS> '}' 'if' <CONDITION>
func (parser *Parser) cmdIf() error {
	errSalt := "Parser.cmdIf"
	parser.accumulateRule("<CMD_IF> ::= '{' <CMDS> '}' 'if' <CONDITION> | '{' <CMDS> '}' 'else' '{' <CMDS> '}' 'if' <CONDITION> | <CMD_ELIF> '{' <CMDS> '}' 'if' <CONDITION>")

	// Due to right-to-left, bottom-to-top, check for 'if' and then the condition.
	// Then look for the commands and braces that follow.

	// All CMD_IF rules end with 'if' <CONDITION>
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

	// Now parse the preceding structure:
	// It can be '{' <CMDS> '}'
	// Or '{' <CMDS> '}' 'else' '{' <CMDS> '}'
	// Or <CMD_ELIF>

	// If the next token is '}', it's the end of a command block.
	// We need to look backward for the structure.

	// This is where the right-to-left parsing becomes critical.
	// The grammar specifies:
	// <CMD_IF> ::= '{' <CMDS> '}' 'if' <CONDITION>
	// <CMD_IF> ::= '{' <CMDS> '}' 'else' '{' <CMDS> '}' 'if' <CONDITION>
	// <CMD_IF> ::= <CMD_ELIF> '{' <CMDS> '}' 'if' <CONDITION>

	// We have already consumed 'if' <CONDITION>. Now we need to handle the parts before it.

	// Scenario 1: <CMD_ELIF>
	if parser.token == TElif {
		if err := parser.cmdElif(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
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

	return nil
}

// cmdElif :
// <CMD_ELIF> ::= '{' <CMDS> '}' 'elif' <CONDITION>
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

// cmdElifRest :
// <CMD_ELIF_REST> ::= '{' <CMDS> '}' 'elif' <CONDITION> <CMD_ELIF_REST>
// <CMD_ELIF_REST> ::= '{' <CMDS> '}' 'else' '{' <CMDS> '}' 'elif' <CONDITION> <CMD_ELIF_REST>
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

// cmdFor :
// <CMD_FOR> ::= '{' <CMDS> '}' 'for' <CONDITION>
func (parser *Parser) cmdFor() error {
	errSalt := "Parser.cmdFor"
	parser.accumulateRule("<CMD_FOR> ::= '{' <CMDS> '}' 'for' <CONDITION>")

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

// cmdDeclaration :
// <CMD_DECLARATION> ::= <E> '=:' <VAR>
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

// cmdAssignment :
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

// cmdReceive :
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

// cmdSend :
// <CMD_SEND> ::= '(' <E> ')' 'Send'
func (parser *Parser) cmdSend() error {
	errSalt := "Parser.cmdSend"
	parser.accumulateRule("<CMD_SEND> ::= '(' <E> ')' 'Send'")

	// Expect 'Send'
	if parser.token != TSend {
		return parser.handleSyntaxError(fmt.Errorf("expected 'Send', got %s", parser.lexeme))
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

	// Expect <E>
	if err := parser.e(); err != nil {
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

// condition :
// <CONDITION> ::= <E> '>' <E> | <E> '>=' <E> | <E> '<>' <E> | <E> '<=' <E> | <E> '<' <E> | <E> '==' <E>
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

// e :
// <E> ::= <T> <E_REST>
// <E_REST> ::= '+' <T> <E_REST> | '-' <T> <E_REST> | ε
func (parser *Parser) e() error {
	errSalt := "Parser.e"
	parser.accumulateRule("<E> ::= <T> <E_REST>")

	// For right-to-left, left-recursive rules like <E> ::= <E> + <T> (or <T> <E_REST>),
	// we will parse the E_REST part first if it's there.
	// Since the grammar is given as <E> ::= <T> <E_REST>, we parse <E_REST> first, then <T>.
	// This implies a non-standard recursive descent for typical left-recursive rules, but
	// given the right-to-left instruction, it makes sense.

	if err := parser.eRest(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	if err := parser.t(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	return nil
}

// eRest :
// <E_REST> ::= '+' <T> <E_REST> | '-' <T> <E_REST> | ε
func (parser *Parser) eRest() error {
	errSalt := "Parser.eRest"
	parser.accumulateRule("<E_REST> ::= '+' <T> <E_REST> | '-' <T> <E_REST> | ε")

	// If the current token is '+' or '-', then it's not epsilon.
	if parser.token == TAdditionOperator || parser.token == TSubtractionOperator {
		// Recursive call for E_REST first
		if err := parser.eRest(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}

		// Then <T>
		if err := parser.t(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}

		// Then the operator
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
	} else {
		parser.accumulateRule("<E_REST> ::= ε")
	}
	return nil
}

// t :
// <T> ::= <F> <T_REST>
// <T_REST> ::= '*' <F> <T_REST> | '/' <F> <T_REST> | '%' <F> <T_REST> | ε
func (parser *Parser) t() error {
	errSalt := "Parser.t"
	parser.accumulateRule("<T> ::= <F> <T_REST>")

	if err := parser.tRest(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	if err := parser.f(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	return nil
}

// tRest :
// <T_REST> ::= '*' <F> <T_REST> | '/' <F> <T_REST> | '%' <F> <T_REST> | ε
func (parser *Parser) tRest() error {
	errSalt := "Parser.tRest"
	parser.accumulateRule("<T_REST> ::= '*' <F> <T_REST> | '/' <F> <T_REST> | '%' <F> <T_REST> | ε")

	if parser.token == TMultiplicationOperator || parser.token == TDivisionOperator || parser.token == TModuleOperator {
		// Recursive call for T_REST first
		if err := parser.tRest(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}

		// Then <F>
		if err := parser.f(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}

		// Then the operator
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
	} else {
		parser.accumulateRule("<T_REST> ::= ε")
	}
	return nil
}

// f :
// <F> ::= -<F> | <X>
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

// x :
// <X> ::= '(' <E> ')' | [0-9]+('.'[0-9]+) | <VAR> | <STRING>
func (parser *Parser) x() error {
	errSalt := "Parser.x"
	parser.accumulateRule("<X> ::= '(' <E> ')' | [0-9]+('.'[0-9]+) | <VAR> | <STRING>")

	if parser.token == TOpenParentheses { // '(' <E> ')'
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
		if err := parser.e(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
		if parser.token != TCloseParentheses {
			return parser.handleSyntaxError(fmt.Errorf(errExpectedCloseParenthesis, parser.lexeme))
		}
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
	} else if parser.token == TGear || parser.token == TTensor { // [0-9]+('.'[0-9]+) (numbers)
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
	} else if parser.token == TId { // <VAR> (which is <ID>)
		if err := parser.varToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
	} else if parser.token == TDoubleQuote { // <STRING> (Omnidrone token is used for strings)
		if err := parser.stringToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
	} else {
		return parser.handleSyntaxError(fmt.Errorf("expected '(', a number, an identifier, or a string, got %s", parser.lexeme))
	}
	return nil
}

// stringToken :
// <STRING> ::= '"' <TEXT_WITH_NUMBERS> '"'
// Note: TEXT_WITH_NUMBERS is implicitly handled by the Lexer as part of TOmnidrone.
func (parser *Parser) stringToken() error {
	parser.accumulateRule("<STRING> ::= '\"' <TEXT_WITH_NUMBERS> '\"'")

	// The lexer identifies the entire string literal (including quotes) as TOmnidrone.
	// So, we just need to consume the TOmnidrone token here.
	if parser.token != TDoubleQuote {
		return parser.handleSyntaxError(fmt.Errorf("expected a string literal, got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err // Propagation of error
	}
	return nil
}

// varToken :
// <VAR> ::= <ID>
func (parser *Parser) varToken() error {
	errSalt := "Parser.varToken"
	parser.accumulateRule("<VAR> ::= <ID>")

	if err := parser.id(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}
	return nil
}

// id :
// <ID> ::= (([A-Z]|[a-z])+(_|[0-9])*)+
func (parser *Parser) id() error {
	parser.accumulateRule("<ID> ::= (([A-Z]|[a-z])+(_|[0-9])*)+")
	if parser.token != TId {
		return parser.handleSyntaxError(fmt.Errorf("expected an Identifier, got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err // Propagation of error
	}
	return nil
}

// parameters :
// <PARAMETERS> ::= <ID> ':' <TYPE>
// <PARAMETERS> ::= <ID> ':' <TYPE> <EXTRA_PARAMETERS>
func (parser *Parser) parameters() error {
	errSalt := "Parser.parameters"
	parser.accumulateRule("<PARAMETERS> ::= <ID> ':' <TYPE> | <ID> ':' <TYPE> <EXTRA_PARAMETERS>")

	// Due to right-to-left:
	// First, try to match EXTRA_PARAMETERS if present.
	if parser.token == TComma { // Check for ',' which starts EXTRA_PARAMETERS
		if err := parser.extraParameters(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}
	}

	// Then, <TYPE>
	if err := parser.typeToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Then, ':'
	if parser.token != TColon {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedColon, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	// Then, <ID>
	if err := parser.id(); err != nil {
		return log.SyntaxErrorf(errSalt, err)
	}

	return nil
}

// extraParameters :
// <EXTRA_PARAMETERS> ::= ',' <ID> ':' <TYPE>
// <EXTRA_PARAMETERS> ::= ',' <ID> ':' <TYPE> ',' <PARAMETERS>
func (parser *Parser) extraParameters() error {
	errSalt := "Parser.extraParameters"
	parser.accumulateRule("<EXTRA_PARAMETERS> ::= ',' <ID> ':' <TYPE> | ',' <ID> ':' <TYPE> ',' <PARAMETERS>")

	// If there's a comma, it's a recursive call or a single extra parameter.
	if parser.token == TComma {
		// Handle the ',' <PARAMETERS> part first if present, then ',' <ID> ':' <TYPE>
		// This needs careful handling for right-to-left.

		// Consume the ','
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return log.SyntaxErrorf(errSalt, err)
		}

		// Check if it's ',' <PARAMETERS> or just ',' <ID> ':' <TYPE>
		// The grammar is ambiguous here for simple lookahead based on "right to left",
		// but since the original parsing was recursive descent, it implies a certain
		// backtracking or lookahead capability. Given the new instruction:
		// <EXTRA_PARAMETERS> ::= ',' <ID> ':' <TYPE> <EXTRA_PARAMETERS> (implicitly, if the next is also a comma)
		// <EXTRA_PARAMETERS> ::= ',' <ID> ':' <TYPE>

		// To be truly right-to-left:
		// First parse <PARAMETERS> if the next token is a comma.
		// Then parse <TYPE>.
		// Then parse ':'.
		// Then parse <ID>.
		// Then parse ','.

		// This recursive descent implementation will just match the next components.
		// It expects: <TYPE>, then ':', then <ID>, then (optionally) a recursive call to extraParameters followed by a comma.

		// Parse <PARAMETERS> if it's the recursive rule: ',' <ID> ':' <TYPE> ',' <PARAMETERS>
		// This is tricky with current token, as <PARAMETERS> starts with <ID>.
		// Assuming the token stream reflects the right-to-left order, the parser would have
		// already seen the innermost <PARAMETERS> if it was a recursive definition.

		// Let's assume the recursive definition implies: (..., ID : TYPE), (..., ID : TYPE)
		// So we look for ID : TYPE and then decide if another comma makes it recursive.

		// Due to "bottom to top, right to left":
		// <EXTRA_PARAMETERS> ::= ',' <ID> ':' <TYPE>
		// <EXTRA_PARAMETERS> ::= ',' <ID> ':' <TYPE> ',' <PARAMETERS>

		// This is likely: parse the ID, then :, then TYPE. If there's another comma after that, then
		// it's the recursive rule, and we call extraParameters again.

		// Try to match the elements of the rule from right to left
		// So first look for <TYPE>
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

		// After matching `ID : TYPE`, if the next token is also a `TComma`, then we recursively call `extraParameters`.
		if parser.token == TComma {
			// Consume the comma for the recursive call
			parser.displayToken()
			if err := parser.advanceToken(); err != nil {
				return log.SyntaxErrorf(errSalt, err)
			}
			if err := parser.parameters(); err != nil { // Call parameters, as EXTRA_PARAMETERS can have a full PARAMETERS on its right.
				return log.SyntaxErrorf(errSalt, err)
			}
		}

	} else {
		// This should not happen if extraParameters was called due to a comma.
		// This means that the outer `parameters` method needs to handle the choice
		// of calling `extraParameters` or not based on lookahead.
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
		log.LogError(parser.errorMessage)
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
