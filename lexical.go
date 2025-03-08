package main

import (
	"bufio"
	"errors"
	"fmt"
	custom_errors "mechanus-compiler/error"
	"os"
	"strings"
)

// Constants for token values
const (
	//	 Construction tokens

	TConstruct      = 1
	TDefineFunction = 2
	TReturn         = 3
	TComma          = 4
	TColom          = 5

	//	 Conditional and repetition tokens

	TIf    = 101
	TElse  = 102
	TElif  = 103
	TFor   = 104
	TBreak = 105

	//	 Structure tokens

	TOpenParentheses       = 201
	TCloseParentheses      = 202
	TOpenBraces            = 203
	TCloseBraces           = 204
	TSingleLineComment     = 205
	TOpenMultilineComment  = 206
	TCloseMultilineComment = 207

	//	 Operator tokens

	TGreaterThanOperator    = 301
	TLessThanOperator       = 302
	TGreaterEqualOperator   = 303
	TLessEqualOperator      = 304
	TEqualOperator          = 305
	TNotEqualOperator       = 306
	TAdditionOperator       = 307
	TSubtractionOperator    = 308
	TMultiplicationOperator = 309
	TDivisionOperator       = 310
	TModuleOperator         = 311
	TExponentiationOperator = 312
	TAndOperator            = 313
	TOrOperator             = 314
	TNotOperator            = 315
	TDeclarationOperator    = 316
	TAttributioOperator     = 317

	//	 Type tokens

	TInteger = 401
	TFLoat   = 402
	TNil     = 404
	TId      = 405

	//	 Control tokens

	TInputEnd = 501
	TLexError = 502
	TNilValue = 503
)

// Constants for token symbols
const (
	//	 Construction tokens

	Construct      = "Construct"
	DefineFunction = "Architect"
	Return         = "Integrate"
	Comma          = ","
	Colom          = ":"

	//	 Conditional and repetition tokens

	If    = "if"
	Else  = "else"
	Elif  = "elif"
	For   = "for"
	Break = "detach"

	//	 Structure tokens

	OpenParentheses       = "("
	CloseParentheses      = ")"
	OpenBraces            = "{"
	CloseBraces           = "}"
	SingleLineComment     = "//"
	OpenMultilineComment  = "/*"
	CloseMultilineComment = "*/"

	//	 Operators

	GreaterThan            = ">"
	LessThanOperator       = "<"
	GreaterEqualOperator   = ">="
	LessEqualOperator      = "<="
	EqualOperator          = "=="
	NotEqualOperator       = "!="
	AdditionOperator       = "+"
	SubtractionOperator    = "-"
	MultiplicationOperator = "*"
	DivisionOperator       = "/"
	ModuleOperator         = "%"
	ExponentiationOperator = "**"
	AndOperator            = "&&"
	OrOperator             = "||"
	NotOperator            = "!"
	DeclarationOperator    = ":="
	AttributioOperator     = "="

	//	 Type tokens

	Integer = "Gear"
	FLoat   = "Tensor"
	Nil     = "Nil"
	Id      = "Id"
)

// lexical struct to hold lexical analyzer state
type lexical struct {
	inputFile        *os.File
	rdInput          *bufio.Reader
	outputFile       *os.File
	lookAhead        rune
	token            int
	lexeme           string
	pointer          int
	inputLine        string
	currentLine      int
	currentColumn    int
	errorMessage     string
	identifiedTokens strings.Builder
}

func NewLexical(inputFile, outputFile *os.File) lexical {
	lex := lexical{
		inputFile:     inputFile,
		outputFile:    outputFile,
		currentLine:   0,
		currentColumn: 0,
		pointer:       0,
		inputLine:     "",
		token:         TNilValue,
		errorMessage:  "",
	}
	return lex
}

func (lex *lexical) GetToken() (string, error) {
	lex.movelookAhead()

	for lex.token != TInputEnd && lex.token != TLexError {
		lex.nextToken()
		lex.displayToken()
	}

	var err error = nil

	if lex.token == TLexError {
		err = errors.New(lex.errorMessage)
		custom_errors.Log(fmt.Sprintf("Lexical error: %s", lex.errorMessage), &err, custom_errors.ErrorLevel)
	} else {
		fmt.Println("Lexical analys completed with no errors")
	}

	lex.showTokens()
	lex.writeOutput()

	return lex.lexeme, err
}

func (lex *lexical) movelookAhead() error {
	if lex.pointer+1 > len(lex.inputLine) {
		lex.currentLine++
		lex.pointer = 0
		line, err := lex.rdInput.ReadString('\n')
		if err != nil {
			lex.lookAhead = EOF
			return err
		}
		lex.inputLine = line
		lex.lookAhead = rune(lex.inputLine[lex.pointer])
	} else {
		lex.lookAhead = rune(lex.inputLine[lex.pointer])
	}
	if lex.lookAhead >= 'a' && lex.lookAhead <= 'z' {
		lex.lookAhead = lex.lookAhead - 'a' + 'A'
	}
	lex.pointer++
	lex.currentColumn = lex.pointer + 1
	return nil
}

func (lex *lexical) nextToken() error {
	var sbLexeme strings.Builder

	for lex.lookAhead == ' ' || lex.lookAhead == '\t' || lex.lookAhead == '\n' || lex.lookAhead == '\r' {
		lex.movelookAhead()
	}

	if lex.lookAhead >= 'A' && lex.lookAhead <= 'Z' {
		lex.alphabeticalCharacter(sbLexeme)
	} else if lex.lookAhead >= '0' && lex.lookAhead <= '9' {
		lex.numericalCharacter(sbLexeme)
	} else {
		lex.symbolCharacter(sbLexeme)
	}

	lex.lexeme = sbLexeme.String()
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

func (lex *lexical) alphabeticalCharacter(sbLexeme strings.Builder) {
	sbLexeme.WriteRune(lex.lookAhead)
	lex.movelookAhead()

	for (lex.lookAhead >= 'A' && lex.lookAhead <= 'Z') || (lex.lookAhead >= '0' && lex.lookAhead <= '9') || lex.lookAhead == '_' {
		sbLexeme.WriteRune(lex.lookAhead)
		lex.movelookAhead()
	}

	lex.lexeme = sbLexeme.String()

	switch strings.ToUpper(lex.lexeme) {
	case "CONSTRUCT":
		lex.token = TConstruct
	case "VARIABLES":
		lex.token = TVariable
	case "IF":
		lex.token = TIf
	case "ELSE":
		lex.token = TElse
	case "ELSE_IF":
		lex.token = TElif
	case "FOR":
		lex.token = TFor
	case "BREAK":
		lex.token = TBreak
	default:
		lex.token = TId
	}
}

func (lex *lexical) numericalCharacter(sbLexeme strings.Builder) {
	sbLexeme.WriteRune(lex.lookAhead)
	lex.movelookAhead()
	for lex.lookAhead >= '0' && lex.lookAhead <= '9' {
		sbLexeme.WriteRune(lex.lookAhead)
		lex.movelookAhead()
	}
	lex.token = TInteger
}

func (lex *lexical) symbolCharacter(sbLexeme strings.Builder) {
	switch lex.lookAhead {
	case '(':
		lex.token = TOpenParentheses
	case ')':
		lex.token = TCloseParentheses
	case '{':
		lex.token = TOpenBraces
	case '}':
		lex.token = TCloseBraces
	case ',':
		lex.token = TComma
	case '+':
		lex.token = TAdditionOperator
	case '-':
		lex.token = TSubtractionOperator
	case '*':
		lex.token = TMultiplicationOperator
	case '/':
		lex.token = TDivisionOperator
	case '%':
		lex.token = TModuleOperator
	case '<':
		lex.token = TLessThanOperator
	case '>':
		lex.token = TGreaterThanOperator
	case '=':
		lex.token = TEqualOperator
	default:
		lex.token = TLexError
		lex.errorMessage = fmt.Sprintf("Erro Léxico na linha: %d\nReconhecido ao atingir a coluna: %d\nLinha do Erro: <%s>\nToken desconhecido: %c", lex.currentLine, lex.currentColumn, lex.inputLine, lex.lookAhead)
	}
	sbLexeme.WriteRune(lex.lookAhead)
	lex.movelookAhead()
}

// ---------------------------------------------------------------------------------------------------------------------

func (lex *lexical) displayToken() {
	var tokenLexeme string
	switch lex.token {
	case TConstruct:
		tokenLexeme = "T_MODULE"
	case TVariable:
		tokenLexeme = "T_VARIABLE"
	case TComma:
		tokenLexeme = "T_COMMA"
	case TIf:
		tokenLexeme = "T_IF"
	case TElse:
		tokenLexeme = "T_ELSE"
	case TElif:
		tokenLexeme = "T_ELSE_IF"
	case TFor:
		tokenLexeme = "T_FOR"
	case TBreak:
		tokenLexeme = "T_BREAK"
	case TOpenParentheses:
		tokenLexeme = "T_OPEN_PARENTHESES"
	case TCloseParentheses:
		tokenLexeme = "T_CLOSE_PARENTHESES"
	case TGreaterThanOperator:
		tokenLexeme = "T_GREATER_THAN_OPERATOR"
	case TLessThanOperator:
		tokenLexeme = "T_LESS_THAN_OPERATOR"
	case TGreaterEqualOperator:
		tokenLexeme = "T_GREATER_EQUAL_OPERATOR"
	case TLessEqualOperator:
		tokenLexeme = "T_LESS_EQUAL_OPERATOR"
	case TEqualOperator:
		tokenLexeme = "T_EQUAL_OPERATOR"
	case TNotEqualOperator:
		tokenLexeme = "T_NOT_EQUAL_OPERATOR"
	case TAdditionOperator:
		tokenLexeme = "T_ADDITION_OPERATOR"
	case TSubtractionOperator:
		tokenLexeme = "T_SUBTRACTION_OPERATOR"
	case TMultiplicationOperator:
		tokenLexeme = "T_MULTIPLICATION_OPERATOR"
	case TDivisionOperator:
		tokenLexeme = "T_DIVISION_OPERATOR"
	case TModuleOperator:
		tokenLexeme = "T_MODULE_OPERATOR"
	case TExponentiationOperator:
		tokenLexeme = "T_EXPONENTIATION_OPERATOR"
	case TInteger:
		tokenLexeme = "T_INTEGER"
	case TFLoat:
		tokenLexeme = "T_FLOAT"
	case TId:
		tokenLexeme = "T_ID"
	case TInputEnd:
		tokenLexeme = "T_INPUT_END"
	case TLexError:
		tokenLexeme = "T_LEX_ERROR"
	case TNilValue:
		tokenLexeme = "T_NIL"
	default:
		tokenLexeme = "N/A"
	}
	fmt.Println(tokenLexeme + " ( " + lex.lexeme + " )")
	lex.storeTokens(tokenLexeme + " ( " + lex.lexeme + " )")
}

func (lex *lexical) close(file string) {
	custom_errors.Log(fmt.Sprintf("closing %s file", file), nil, custom_errors.InfoLevel)

	switch file {
	case "input":
		err := lex.inputFile.Close()
		if err != nil {
			custom_errors.Log(custom_errors.FileCloseError, &err, custom_errors.ErrorLevel)
			return
		}
	case "output":
		err := lex.outputFile.Close()
		if err != nil {
			custom_errors.Log(custom_errors.FileCloseError, &err, custom_errors.ErrorLevel)
			return
		}
	}

	custom_errors.Log(custom_errors.FileCloseSuccess, nil, custom_errors.InfoLevel)
}

func (lex *lexical) writeOutput() error {
	if lex.outputFile == nil {

		return fmt.Errorf(custom_errors.UninitializedFile)
	}
	file, err := os.Create("output.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(lex.identifiedTokens.String())
	if err != nil {
		return err
	}
	fmt.Println("Arquivo Salvo: output.txt")
	return nil
}

func (lex *lexical) showTokens() {
	fmt.Println("Identified Tokens (token/lexeme):")
	fmt.Println(lex.identifiedTokens.String())
}

func (lex *lexical) storeTokens(tokenIdentificado string) {
	lex.identifiedTokens.WriteString(tokenIdentificado)
	lex.identifiedTokens.WriteString("\n")
}
