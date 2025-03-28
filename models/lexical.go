package models

import (
	"bufio"
	"fmt"
	custom_errors "mechanus-compiler/error"
	"os"
	"strings"
)

// Constants for Token values
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

// Constants for Token symbols
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

// Lexical struct to hold Lexical analyzer state
type Lexical struct {
	InputFile        *os.File
	RdInput          *bufio.Reader
	OutputFile       *os.File
	LookAhead        rune
	Token            int
	Lexeme           string
	Pointer          int
	InputLine        string
	CurrentLine      int
	CurrentColumn    int
	ErrorMessage     string
	IdentifiedTokens strings.Builder
}

func NewLexical(inputFile, outputFile *os.File) Lexical {
	lex := Lexical{
		InputFile:     inputFile,
		OutputFile:    outputFile,
		CurrentLine:   0,
		CurrentColumn: 0,
		Pointer:       0,
		InputLine:     "",
		Token:         TNilValue,
		ErrorMessage:  "",
	}
	return lex
}

func (lex *Lexical) movelookAhead() error {
	if lex.Pointer+1 > len(lex.InputLine) {
		lex.CurrentLine++
		lex.Pointer = 0
		line, err := lex.RdInput.ReadString('\n')
		if err != nil {
			lex.LookAhead = EOF
			return err
		}
		lex.InputLine = line
		lex.LookAhead = rune(lex.InputLine[lex.Pointer])
	} else {
		lex.LookAhead = rune(lex.InputLine[lex.Pointer])
	}
	if lex.LookAhead >= 'a' && lex.LookAhead <= 'z' {
		lex.LookAhead = lex.LookAhead - 'a' + 'A'
	}
	lex.Pointer++
	lex.CurrentColumn = lex.Pointer + 1
	return nil
}

func (lex *Lexical) nextToken() error {
	var sbLexeme strings.Builder

	for lex.LookAhead == ' ' || lex.LookAhead == '\t' || lex.LookAhead == '\n' || lex.LookAhead == '\r' {
		lex.movelookAhead()
	}

	if lex.LookAhead >= 'A' && lex.LookAhead <= 'Z' {
		lex.alphabeticalCharacter(sbLexeme)
	} else if lex.LookAhead >= '0' && lex.LookAhead <= '9' {
		lex.numericalCharacter(sbLexeme)
	} else {
		lex.symbolCharacter(sbLexeme)
	}

	lex.Lexeme = sbLexeme.String()
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

func (lex *Lexical) alphabeticalCharacter(sbLexeme strings.Builder) {
	sbLexeme.WriteRune(lex.LookAhead)
	lex.movelookAhead()

	for (lex.LookAhead >= 'A' && lex.LookAhead <= 'Z') || (lex.LookAhead >= '0' && lex.LookAhead <= '9') || lex.LookAhead == '_' {
		sbLexeme.WriteRune(lex.LookAhead)
		lex.movelookAhead()
	}

	lex.Lexeme = sbLexeme.String()

	switch strings.ToUpper(lex.Lexeme) {
	case "CONSTRUCT":
		lex.Token = TConstruct
	case "VARIABLES":
		lex.Token = TVariable
	case "IF":
		lex.Token = TIf
	case "ELSE":
		lex.Token = TElse
	case "ELSE_IF":
		lex.Token = TElif
	case "FOR":
		lex.Token = TFor
	case "BREAK":
		lex.Token = TBreak
	default:
		lex.Token = TId
	}
}

func (lex *Lexical) numericalCharacter(sbLexeme strings.Builder) {
	sbLexeme.WriteRune(lex.LookAhead)
	lex.movelookAhead()
	for lex.LookAhead >= '0' && lex.LookAhead <= '9' {
		sbLexeme.WriteRune(lex.LookAhead)
		lex.movelookAhead()
	}
	lex.Token = TInteger
}

func (lex *Lexical) symbolCharacter(sbLexeme strings.Builder) {
	switch lex.LookAhead {
	case '(':
		lex.Token = TOpenParentheses
	case ')':
		lex.Token = TCloseParentheses
	case '{':
		lex.Token = TOpenBraces
	case '}':
		lex.Token = TCloseBraces
	case ',':
		lex.Token = TComma
	case '+':
		lex.Token = TAdditionOperator
	case '-':
		lex.Token = TSubtractionOperator
	case '*':
		lex.Token = TMultiplicationOperator
	case '/':
		lex.Token = TDivisionOperator
	case '%':
		lex.Token = TModuleOperator
	case '<':
		lex.Token = TLessThanOperator
	case '>':
		lex.Token = TGreaterThanOperator
	case '=':
		lex.Token = TEqualOperator
	default:
		lex.Token = TLexError
		lex.ErrorMessage = fmt.Sprintf("Erro Léxico na linha: %d\nReconhecido ao atingir a coluna: %d\nLinha do Erro: <%s>\nToken desconhecido: %c", lex.CurrentLine, lex.CurrentColumn, lex.InputLine, lex.LookAhead)
	}
	sbLexeme.WriteRune(lex.LookAhead)
	lex.movelookAhead()
}

// ---------------------------------------------------------------------------------------------------------------------

func (lex *Lexical) displayToken() {
	var tokenLexeme string
	switch lex.Token {
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
	fmt.Println(tokenLexeme + " ( " + lex.Lexeme + " )")
	lex.storeTokens(tokenLexeme + " ( " + lex.Lexeme + " )")
}

func (lex *Lexical) Close(file string) {
	custom_errors.Log(fmt.Sprintf("closing %s file", file), nil, custom_errors.InfoLevel)

	switch file {
	case "input":
		err := lex.InputFile.Close()
		if err != nil {
			custom_errors.Log(custom_errors.FileCloseError, &err, custom_errors.ErrorLevel)
			return
		}
	case "output":
		err := lex.OutputFile.Close()
		if err != nil {
			custom_errors.Log(custom_errors.FileCloseError, &err, custom_errors.ErrorLevel)
			return
		}
	}

	custom_errors.Log(custom_errors.FileCloseSuccess, nil, custom_errors.InfoLevel)
}

func (lex *Lexical) writeOutput() error {
	if lex.OutputFile == nil {

		return fmt.Errorf(custom_errors.UninitializedFile)
	}
	file, err := os.Create("output.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(lex.IdentifiedTokens.String())
	if err != nil {
		return err
	}
	fmt.Println("Arquivo Salvo: output.txt")
	return nil
}

func (lex *Lexical) showTokens() {
	fmt.Println("Identified Tokens (Token/Lexeme):")
	fmt.Println(lex.IdentifiedTokens.String())
}

func (lex *Lexical) storeTokens(tokenIdentificado string) {
	lex.IdentifiedTokens.WriteString(tokenIdentificado)
	lex.IdentifiedTokens.WriteString("\n")
}
