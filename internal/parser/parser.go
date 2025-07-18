package parser

import (
	"fmt"
	"mechanus-compiler/internal/compiler_error"
	"mechanus-compiler/internal/lexer"
	"mechanus-compiler/internal/logger"
	"os"
	"strings"
)

// Parser :
// This is the structure responsible for making the syntactical analysis of the source file. It checks for unrecognized
// syntaxes and, if it finds one, it returns an error code.
type Parser struct {
	logger          *logger.Logger
	debug           bool // Restored for controlling debug-specific output
	lexer           lexer.Lexer
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
	// Initialize the logger. Log to Stderr. Set level based on the debug flag.
	logLevel := logger.LevelInfo
	if debug {
		logLevel = logger.LevelDebug
	}
	lg := logger.New(os.Stderr, logLevel)

	// Initialize the Lexer
	lex, err := lexer.NewLexer(inputFile, outputFile, debug)
	if err != nil {
		// The lexer's constructor will have already logged the error.
		return Parser{}, err
	}

	// Initialize the structure
	parser := Parser{
		logger:       lg,
		debug:        debug, // Set the debug flag
		lexer:        lex,
		outputFile:   outputFile,
		token:        lexer.TNilValue,
		errorMessage: nil,
	}

	return parser, nil
}

// Run :
// Starts the syntactical analysis.
//
// Fails if the lexer fails or if a syntactical error is found.
func (parser *Parser) Run() error {
	if err := parser.advanceToken(); err != nil {
		// The lexer logs its own errors, so we just propagate the error up.
		return err
	}

	if err := parser.g(); err != nil {
		// Parsing functions log their own errors via handleSyntaxError.
		return err
	}

	parser.logger.Info(compiler_error.SyntaxSuccess, nil)
	return nil
}

// g :
// <G> ::= '{' <BODY> '}' <ID> 'Construct'
func (parser *Parser) g() error {
	parser.accumulateRule("<G> ::= '{' <BODY> '}' <ID> 'Construct'")

	// Expect 'Construct'
	if parser.token != lexer.TConstruct {
		return parser.handleSyntaxError(fmt.Errorf("expected 'Construct', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect <ID>
	if parser.token != lexer.TId {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedIdentifier, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect '}'
	if parser.token != lexer.TCloseBraces {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedCloseBraces, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect <BODY>
	if err := parser.body(); err != nil {
		return err
	}

	// Expect '{'
	if parser.token != lexer.TOpenBraces {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedOpenBraces, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		if strings.Contains(err.Error(), compiler_error.EndOfFileReached) {
			return nil
		}
		return err
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
	parser.accumulateRule("<BODY> ::= <BODY_REST> '{' <CMDS> '}' '(' <PARAMETERS> ')' <ID> 'Architect' | ...")

	// 1. Expect 'Architect'
	if parser.token != lexer.TArchitect {
		return parser.handleSyntaxError(fmt.Errorf("expected 'Architect', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// 2. Expect <ID>
	if parser.token != lexer.TId {
		return parser.handleSyntaxError(fmt.Errorf("expected ID after Architect, got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// 3. Expect ')'
	if parser.token != lexer.TCloseParentheses {
		return parser.handleSyntaxError(fmt.Errorf("expected ')', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// 4. Optionally parse <PARAMETERS> (may be ε)
	if parser.token != lexer.TOpenParentheses {
		_ = parser.parametersDecl() // fail silently if no parametersDecl
	}

	// 5. Expect ')'
	if parser.token != lexer.TOpenParentheses {
		return parser.handleSyntaxError(fmt.Errorf("expected '(', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// 6. Expect '}'
	if parser.token != lexer.TCloseBraces {
		return parser.handleSyntaxError(fmt.Errorf("expected '}', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// 7. Parse <CMDS>
	if err := parser.cmds(); err != nil {
		return err
	}

	// 8. Expect '{'
	if parser.token != lexer.TOpenBraces {
		return parser.handleSyntaxError(fmt.Errorf("expected '{', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		if strings.Contains(err.Error(), compiler_error.EndOfFileReached) {
			return nil
		}
		return err
	}

	// 9. Recursively parse any additional Architect bodies
	if err := parser.bodyRest(); err != nil {
		if strings.Contains(err.Error(), compiler_error.EndOfFileReached) {
			return nil
		}
		return err
	}

	return nil
}

// <BODY_REST> :
//
// <BODY_REST> ::= <BODY_REST> '{' <CMDS> '}' '(' <PARAMETERS_DECL> ')' <ID> 'Architect'
// <BODY_REST> ::= <BODY_REST> '{' <CMDS> '}' <TYPE> '(' <PARAMETERS_DECL> ')' <ID> 'Architect'
// <BODY_REST> ::= ε
func (parser *Parser) bodyRest() error {
	parser.accumulateRule("<BODY_REST> ::= <BODY_REST> '{' <CMDS> '}' '(' <PARAMETERS> ')' <ID> 'Architect' | ... | ε")

	// 1. Base case: ε
	if parser.token == lexer.TCloseBraces || parser.token == lexer.TInputEnd {
		parser.accumulateRule("<BODY_REST> ::= ε")
		return nil
	}

	// 2. Expect 'Architect'
	if parser.token != lexer.TArchitect {
		return parser.handleSyntaxError(fmt.Errorf("expected 'Architect', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// 3. Expect <ID>
	if parser.token != lexer.TId {
		return parser.handleSyntaxError(fmt.Errorf("expected ID after Architect, got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// 4. Expect ')'
	if parser.token != lexer.TCloseParentheses {
		return parser.handleSyntaxError(fmt.Errorf("expected ')', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// 5. Optionally parse <PARAMETERS>
	if parser.token != lexer.TOpenParentheses {
		_ = parser.parametersDecl() // fails silently if epsilon
	}

	// 6. Expect ')'
	if parser.token != lexer.TOpenParentheses {
		return parser.handleSyntaxError(fmt.Errorf("expected '(', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// 7. Expect '}'
	if parser.token != lexer.TCloseBraces {
		return parser.handleSyntaxError(fmt.Errorf("expected '}', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// 8. Parse <CMDS>
	if err := parser.cmds(); err != nil {
		return err
	}

	// 9. Expect '{'
	if parser.token != lexer.TOpenBraces {
		return parser.handleSyntaxError(fmt.Errorf("expected '{', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
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
	if parser.token != lexer.TNil && parser.token != lexer.TGear && parser.token != lexer.TTensor &&
		parser.token != lexer.TState && parser.token != lexer.TMonodrone && parser.token != lexer.TOmnidrone {
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
		for parser.token == lexer.TNewLine {
			parser.displayToken()
			if err := parser.advanceToken(); err != nil {
				return err
			}
		}

		// At this level, hitting '{' means the parser is done with <CMDS>
		if parser.token == lexer.TOpenBraces {
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

	if parser.token == lexer.TNewLine {
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

	if parser.token == lexer.TCloseParentheses {
		parser.accumulateRule("<CMD_CALL> ::= '(' <PARAMETERS_CALL> ')' <ID>")

		// Parse parameters
		if err := parser.parametersCall(); err != nil {
			return err
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
	parser.accumulateRule("<CMD_IF> ::= '{' <CMDS> '}' 'if' <CONDITION> | '{' <CMDS> '}' 'else' '{' <CMDS> '}' 'if' <CONDITION> | <CMD_ELIF> '{' <CMDS> '}' 'if' <CONDITION>")

	// Expect 'if'
	if parser.token != lexer.TIf {
		return parser.handleSyntaxError(fmt.Errorf("expected 'if', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect <CONDITION>
	if err := parser.condition(); err != nil {
		return err
	}

	// Expect '}'
	if parser.token != lexer.TCloseBraces {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedCloseBraces, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect <CMDS>
	if err := parser.cmds(); err != nil {
		return err
	}

	// Expect '{'
	if parser.token != lexer.TOpenBraces {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedOpenBraces, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Check for 'elif'
	for parser.token == lexer.TElif {
		if err := parser.cmdElif(); err != nil {
			return err
		}
	}

	// Check for 'else'
	if parser.token == lexer.TElse {
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return err
		}
		// Expect '}' after 'else' block
		if parser.token != lexer.TCloseBraces {
			return parser.handleSyntaxError(fmt.Errorf(errExpectedCloseBraces, parser.lexeme))
		}
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return err
		}
		// Expect <CMDS> for 'else' block
		if err := parser.cmds(); err != nil {
			return err
		}
		// Expect '{' for 'else' block
		if parser.token != lexer.TOpenBraces {
			return parser.handleSyntaxError(fmt.Errorf(errExpectedOpenBraces, parser.lexeme))
		}
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return err
		}
	}

	return nil
}

// <CMD_ELIF> :
//
// <CMD_ELIF> ::= '{' <CMDS> '}' <CONDITION> 'elif'
// <CMD_ELIF> ::= <CMD_ELIF_REST>
func (parser *Parser) cmdElif() error {
	parser.accumulateRule("<CMD_ELIF> ::= '{' <CMDS> '}' 'elif' <CONDITION> | <CMD_ELIF_REST>")

	// If the current token is TElif, it's a direct elif. Otherwise, it must be CMD_ELIF_REST.
	if parser.token != lexer.TElif {
		return parser.cmdElifRest()
	}

	// Expect 'elif'
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect <CONDITION>
	if err := parser.condition(); err != nil {
		return err
	}

	// Expect '}'
	if parser.token != lexer.TCloseBraces {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedCloseBraces, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect <CMDS>
	if err := parser.cmds(); err != nil {
		return err
	}

	// Expect '{'
	if parser.token != lexer.TOpenBraces {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedOpenBraces, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	return nil
}

// <CMD_ELIF_REST> :
//
// <CMD_ELIF_REST> ::= <CMD_ELIF_REST> '{' <CMDS> '}' <CONDITION> 'elif'
// <CMD_ELIF_REST> ::= '{' <CMDS> '}' 'else' <CMD_ELIF_REST> '{' <CMDS> '}' <CONDITION> 'elif'
// <CMD_ELIF_REST> ::= ε
func (parser *Parser) cmdElifRest() error {
	parser.accumulateRule("<CMD_ELIF_REST> ::= '{' <CMDS> '}' 'elif' <CONDITION> <CMD_ELIF_REST> | '{' <CMDS> '}' 'else' '{' <CMDS> '}' 'elif' <CONDITION> <CMD_ELIF_REST> | ε")

	// If the next token is not 'elif', it's epsilon.
	if parser.token != lexer.TElif {
		parser.accumulateRule("<CMD_ELIF_REST> ::= ε")
		return nil
	}

	// Expect 'elif'
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect <CONDITION>
	if err := parser.condition(); err != nil {
		return err
	}

	// Expect '}'
	if parser.token != lexer.TCloseBraces {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedCloseBraces, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect <CMDS>
	if err := parser.cmds(); err != nil {
		return err
	}

	// Expect '{'
	if parser.token != lexer.TOpenBraces {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedOpenBraces, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Check for 'else'
	if parser.token == lexer.TElse {
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return err
		}
		// Expect '}' after 'else' block
		if parser.token != lexer.TCloseBraces {
			return parser.handleSyntaxError(fmt.Errorf(errExpectedCloseBraces, parser.lexeme))
		}
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return err
		}
		// Expect <CMDS> for 'else' block
		if err := parser.cmds(); err != nil {
			return err
		}
		// Expect '{' for 'else' block
		if parser.token != lexer.TOpenBraces {
			return parser.handleSyntaxError(fmt.Errorf(errExpectedOpenBraces, parser.lexeme))
		}
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return err
		}
	}

	// Recursive call for CMD_ELIF_REST
	if err := parser.cmdElifRest(); err != nil {
		return err
	}

	return nil
}

// <CMD_FOR> :
// <CMD_FOR> ::= '{' <CMDS> '}' <CONDITION> 'for'
func (parser *Parser) cmdFor() error {
	parser.accumulateRule("<CMD_FOR> ::= '{' <CMDS> '}' <CONDITION> 'for'")

	// Expect 'for'
	if parser.token != lexer.TFor {
		return parser.handleSyntaxError(fmt.Errorf("expected 'for', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect <CONDITION>
	if err := parser.condition(); err != nil {
		return err
	}

	// Expect '}'
	if parser.token != lexer.TCloseBraces {
		return parser.handleSyntaxError(fmt.Errorf("expected '}', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect <CMDS>
	if err := parser.cmds(); err != nil {
		return err
	}

	// Expect '{'
	if parser.token != lexer.TOpenBraces {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedOpenBraces, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	return nil
}

// <CMD_INTEGRATE> :
//
// <CMD_INTEGRATE> ::= <E> 'Integrate'
func (parser *Parser) cmdIntegrate() error {
	parser.accumulateRule("<CMD_INTEGRATE> ::= <E> 'Integrate'")

	// Expect 'Integrate'
	if parser.token != lexer.TIntegrate {
		return parser.handleSyntaxError(fmt.Errorf("expected 'Integrate', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect <E>
	if err := parser.e(); err != nil {
		return err
	}

	return nil
}

// <CMD_DECLARATION> :
//
// <CMD_DECLARATION> ::= <E> '=:' <TYPE> ':' <VAR>
func (parser *Parser) cmdDeclaration() error {
	parser.accumulateRule("<CMD_DECLARATION> ::= <E> '=:' <TYPE> ':' <VAR>")

	// Expect <VAR>
	if err := parser.varToken(); err != nil {
		return err
	}

	// Expect ':'
	if parser.token != lexer.TColon {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedColon, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect <TYPE>
	if err := parser.typeToken(); err != nil {
		return err
	}

	// Expect '=:'
	if parser.token != lexer.TDeclarationOperator {
		return parser.handleSyntaxError(fmt.Errorf("expected '=:', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect <E>
	if err := parser.e(); err != nil {
		return err
	}

	return nil
}

// <CMD_ASSIGNMENT> :
//
// <CMD_ASSIGNMENT> ::= <E> '=' <VAR>
// <CMD_ASSIGNMENT> :
//
// <CMD_ASSIGNMENT> ::= <E> '=' <VAR>
func (parser *Parser) cmdAssignment() error {
	parser.accumulateRule("<CMD_ASSIGNMENT> ::= <E> '=' <VAR>")

	// Expect <VAR>
	if err := parser.varToken(); err != nil {
		return err
	}

	// Expect '='
	if parser.token != lexer.TAttributionOperator {
		return parser.handleSyntaxError(fmt.Errorf("expected '=', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect <E>
	if err := parser.e(); err != nil {
		return err
	}

	return nil
}

// <CMD_RECEIVE> :
//
// <CMD_RECEIVE> ::= '(' <VAR> ')' 'Receive'
func (parser *Parser) cmdReceive() error {
	parser.accumulateRule("<CMD_RECEIVE> ::= '(' <VAR> ')' 'Receive'")

	// Expect 'Receive'
	if parser.token != lexer.TReceive {
		return parser.handleSyntaxError(fmt.Errorf("expected 'Receive', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect ')'
	if parser.token != lexer.TCloseParentheses {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedCloseParenthesis, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect <VAR>
	if err := parser.varToken(); err != nil {
		return err
	}

	// Expect ')'
	if parser.token != lexer.TOpenParentheses {
		return parser.handleSyntaxError(fmt.Errorf(errExpectedOpenParenthesis, parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	return nil
}

// <CMD_SEND> :
//
// <CMD_SEND> ::= '(' <E> ')' 'Send'
func (parser *Parser) cmdSend() error {
	parser.accumulateRule("<CMD_SEND> ::= '(' <E> ')' 'Send'")

	// Expect TSend (first, since lexing is bottom-up, right-to-left)
	if parser.token != lexer.TSend {
		return parser.handleSyntaxError(fmt.Errorf("expected 'Send', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect TCloseParentheses
	if parser.token != lexer.TCloseParentheses {
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
	if parser.token != lexer.TOpenParentheses {
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
	parser.accumulateRule("<CONDITION> ::= <E> '>' <E> | <E> '>=' <E> | <E> '<>' <E> | <E> '<=' <E> | <E> '<' <E> | <E> '==' <E>")

	// All conditions are of the form <E> OPERATOR <E>
	// Parse the second <E> (rightmost) first
	if err := parser.e(); err != nil {
		return err
	}

	// Expect a comparison operator
	if parser.token != lexer.TGreaterThanOperator &&
		parser.token != lexer.TGreaterEqualOperator &&
		parser.token != lexer.TLessThanOperator &&
		parser.token != lexer.TLessEqualOperator &&
		parser.token != lexer.TNotEqualOperator &&
		parser.token != lexer.TEqualOperator {
		return parser.handleSyntaxError(fmt.Errorf("expected a comparison operator, got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Parse the first <E> (leftmost)
	if err := parser.e(); err != nil {
		return err
	}

	return nil
}

// <E> :
// <E> ::= <E_REST> <T>
func (parser *Parser) e() error {
	parser.accumulateRule("<E> ::= <T> <E_REST>")

	if err := parser.t(); err != nil {
		return err
	}
	return parser.eRest()
}

// eRest :
// <E_REST> ::= <E_REST> '+' <T>
// <E_REST> ::= <E_REST> '-' <T>
// <E_REST> ::= ε
func (parser *Parser) eRest() error {
	switch parser.token {
	case lexer.TAdditionOperator, lexer.TSubtractionOperator:
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return err
		}
		if err := parser.t(); err != nil {
			return err
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
	parser.accumulateRule("<T> ::= <F> <T_REST>")

	if err := parser.f(); err != nil {
		return err
	}

	return parser.tRest()
}

// <T_REST>
func (parser *Parser) tRest() error {
	switch parser.token {
	case lexer.TMultiplicationOperator, lexer.TDivisionOperator, lexer.TModuleOperator:
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return err
		}
		if err := parser.f(); err != nil {
			return err
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
	parser.accumulateRule("<F> ::= -<F> | <X>")

	if parser.token == lexer.TSubtractionOperator {
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return err
		}
		if err := parser.f(); err != nil {
			return err
		}
	} else {
		if err := parser.x(); err != nil {
			return err
		}
	}
	return nil
}

// <X> :
func (parser *Parser) x() error {
	parser.accumulateRule("<X> ::= '(' <E> ')' | [0-9]+('.'[0-9]+) | <STRING> | <NIL> | <VAR> | '(' <PARAMETERS_CALL> ')' <ID>")

	switch parser.token {

	// Case: STRING literal
	case lexer.TDoubleQuote:
		parser.displayToken()
		return parser.advanceToken()

	// Case: NIL
	case lexer.TNil:
		parser.displayToken()
		return parser.advanceToken()

	// Case: numeric literal (integer or float)
	case lexer.TGear, lexer.TTensor:
		parser.displayToken()
		return parser.advanceToken()

	// Case: identifier (variable or function call)
	case lexer.TId:
		parser.displayToken()
		return parser.advanceToken()

	// Case: open parentheses — could be (E) or (PARAMS_CALL) ID
	case lexer.TCloseParentheses:
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return err
		}

		// Try to parse parametersCall (supports multiple expressions)
		if err := parser.parametersCall(); err != nil {
			return err
		}

		// Expect matching opening parenthesis
		if parser.token != lexer.TOpenParentheses {
			return parser.handleSyntaxError(fmt.Errorf("expected '(', got %s", parser.lexeme))
		}
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return err
		}

		return nil
	}

	// If no valid rule matches, return error
	return parser.handleSyntaxError(fmt.Errorf("unexpected token in <X>: %s", parser.lexeme))
}

func (parser *Parser) nilToken() error {
	parser.accumulateRule("<NIL> :: 'Nil'")

	if parser.token != lexer.TNil {
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
	if parser.token != lexer.TDoubleQuote {
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
	parser.accumulateRule("<VAR> ::= <ID>")

	if err := parser.id(); err != nil {
		return err
	}
	return nil
}

// <ID> :
//
// <ID> ::= (([A-Z]|[a-z])+(_|[0-9])*)+
func (parser *Parser) id() error {
	parser.accumulateRule("<ID> ::= (([A-Z]|[a-z])+(_|[0-9])*)+")
	if parser.token != lexer.TId {
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
	parser.accumulateRule("<PARAMETERS> ::= <EXTRA_PARAMETERS> <TYPE> ':' <ID> | <TYPE> ':' <ID>")

	// Expect ID (rightmost identifier in the parameter list)
	if parser.token != lexer.TId {
		return parser.handleSyntaxError(fmt.Errorf("expected parameter ID, got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect ':'
	if parser.token != lexer.TColon {
		return parser.handleSyntaxError(fmt.Errorf("expected ':', got %s", parser.lexeme))
	}
	parser.displayToken()
	if err := parser.advanceToken(); err != nil {
		return err
	}

	// Expect <TYPE>
	if err := parser.typeToken(); err != nil {
		return err
	}

	// Loop to check for extra parametersDecl (reverse order)
	for parser.token == lexer.TComma {
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return err
		}

		// Expect ID
		if parser.token != lexer.TId {
			return parser.handleSyntaxError(fmt.Errorf("expected parameter ID, got %s", parser.lexeme))
		}
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return err
		}

		// Expect ':'
		if parser.token != lexer.TColon {
			return parser.handleSyntaxError(fmt.Errorf("expected ':', got %s", parser.lexeme))
		}
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return err
		}

		// Expect <TYPE>
		if err := parser.typeToken(); err != nil {
			return err
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
	parser.accumulateRule("<PARAMETERS_CALL> ::= <EXTRA_PARAMETERS_CALL> <E> | <E>")

	// Parse rightmost expression (last param)
	if err := parser.e(); err != nil {
		return err
	}

	// Repeatedly handle comma-separated expressions
	for parser.token == lexer.TComma {
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return err
		}

		if err := parser.e(); err != nil {
			return err
		}
	}

	return nil
}

// <EXTRA_PARAMETERS> :
//
// <EXTRA_PARAMETERS> ::= ',' <ID> ':' <TYPE>
// <EXTRA_PARAMETERS> ::= ',' <ID> ':' <TYPE> ',' <PARAMETERS>
func (parser *Parser) extraParameters() error {
	parser.accumulateRule("<EXTRA_PARAMETERS> ::= ',' <ID> ':' <TYPE> | ',' <ID> ':' <TYPE> ',' <PARAMETERS>")

	// If there's a comma, it's a recursive call or a single extra parameter.
	if parser.token == lexer.TComma {
		// Consume the ','
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return err
		}

		// Check for <TYPE>
		if err := parser.typeToken(); err != nil {
			return err
		}

		// Then ':'
		if parser.token != lexer.TColon {
			return parser.handleSyntaxError(fmt.Errorf(errExpectedColon, parser.lexeme))
		}
		parser.displayToken()
		if err := parser.advanceToken(); err != nil {
			return err
		}

		// Then <ID>
		if err := parser.id(); err != nil {
			return err
		}

		// After matching `ID : TYPE`, if the next token is also a `TComma`, then recursively call `extraParameters`.
		if parser.token == lexer.TComma {
			// Consume the comma for the recursive call
			parser.displayToken()
			if err := parser.advanceToken(); err != nil {
				return err
			}
			if err := parser.parametersDecl(); err != nil {
				return err
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
	parser.logger.Debug("Advancing token...", nil)

	token, err := parser.lexer.NextToken()
	if err != nil {
		// The lexer logs its own errors. We just propagate it.
		return err
	}

	parser.token = token
	parser.lexeme = parser.lexer.GetLexeme()

	return nil
}

// displayToken :
// Displays the current token and lexeme if debug mode is enabled.
func (parser *Parser) displayToken() {
	if parser.debug {
		parser.lexer.DisplayToken()
	}
}

// handleSyntaxError :
// Records and logs a syntax error, ensuring it is only logged once.
//
// Returns a new error of type ErrSyntax.
func (parser *Parser) handleSyntaxError(err error) error {
	if parser.errorMessage == nil {
		// Create the detailed error.
		detailedErr := fmt.Errorf("%s at %s", err.Error(), parser.lexer.DisplayPos())
		parser.errorMessage = compiler_error.SyntaxErrorf(compiler_error.SyntaxError, detailedErr)

		// Log the structured error.
		parser.logger.Error(parser.errorMessage, map[string]any{
			"position": parser.lexer.DisplayPos(),
			"lexeme":   parser.lexeme,
		})
	}
	return parser.errorMessage
}

// accumulateRule :
// Accumulates the recognized grammar rule for debug output.
func (parser *Parser) accumulateRule(rule string) {
	parser.logger.Debug("Recognized rule", map[string]any{"rule": rule})
	if parser.debug {
		parser.recognizedRules.WriteString(fmt.Sprintf("Recognized rule: %s\n", rule))
	}
}

// ShowRecognizedRules :
// Displays all recognized grammar rules.
func (parser *Parser) ShowRecognizedRules() {
	if parser.debug {
		parser.logger.Debug("Showing all recognized grammar rules.", nil)
		fmt.Println(parser.recognizedRules.String())
	}
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
