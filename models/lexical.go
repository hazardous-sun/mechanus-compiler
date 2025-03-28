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

	TConstruct = 1
	TArchitect = 2
	TIntegrate = 3
	TComma     = 4
	TColon     = 5

	//	 Conditional and repetition tokens

	TIf     = 101
	TElse   = 102
	TElif   = 103
	TFor    = 104
	TDetach = 105

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
	TAttributionOperator    = 317

	//	 Type tokens

	TNil       = 401
	TGear      = 402
	TTensor    = 403
	TState     = 404
	TMonodrone = 405
	TOmnidrone = 406
	TId        = 407

	//	 Control tokens

	TInputEnd = 501
	TLexError = 502
	TNilValue = 503
)

// Constants for Token symbols
const (
	//	 Construction tokens

	Construct = "CONSTRUCT"
	Architect = "ARCHITECT"
	Integrate = "INTEGRATE"

	//	 Conditional and repetition tokens

	If     = "IF"
	Else   = "ELSE"
	Elif   = "ELIF"
	For    = "FOR"
	Detach = "DETACH"

	//	 Type tokens

	Nil       = "NIL"
	Gear      = "GEAR"
	Tensor    = "TENSOR"
	State     = "STATE"
	Monodrone = "MONODRONE"
	Omnidrone = "OMNIDRONE"
	Id        = "ID"
)

// Constants for unique-symbol tokens
const (
	// Construction tokens

	Comma = ','
	Colon = ':'

	//	 Structure tokens

	OpenParentheses  = '('
	CloseParentheses = ')'
	OpenBraces       = '{'
	CloseBraces      = '}'

	//	 Operators

	GreaterThanOperator    = '>'
	LessThanOperator       = '<'
	AdditionOperator       = '+'
	SubtractionOperator    = '-'
	MultiplicationOperator = '*'
	DivisionOperator       = '/'
	ModuleOperator         = '%'
	NotOperator            = '!'
	AttributionOperator    = '='
)

// Constants for multi-symbol tokens
const (
	// Structure tokens

	SingleLineComment     = "//"
	OpenMultilineComment  = "/*"
	CloseMultilineComment = "*/"

	// Operators

	GreaterEqualOperator = ">="
	LessEqualOperator    = "<="
	EqualOperator        = "=="
	NotEqualOperator     = "!="
	AndOperator          = "&&"
	OrOperator           = "||"
	DeclarationOperator  = ":="
)

// Constants for output
const (
	//   Construction tokens

	OutputConstruct = "T_CONSTRUCT"
	OutputArchitect = "T_ARCHITECT"
	OutputIntegrate = "T_INTEGRATE"
	OutputComma     = "T_COMMA"
	OutputColon     = "T_COLON"

	//   Conditional and repetition tokens

	OutputIf     = "T_IF"
	OutputElse   = "T_ELSE"
	OutputElif   = "T_ELIF"
	OutputFor    = "T_FOR"
	OutputDetach = "T_DETACH"

	//   Type tokens

	OutputNil       = "T_NIL"
	OutputGear      = "T_GEAR"
	OutputTensor    = "T_TENSOR"
	OutputState     = "T_STATE"
	OutputMonodrone = "T_MONODRONE"
	OutputOmnidrone = "T_OMNIDRONE"
	OutputId        = "T_ID"

	//   Structure tokens

	OutputOpenParentheses       = "T_OPEN_PARENTHESES"
	OutputCloseParentheses      = "T_CLOSE_PARENTHESES"
	OutputOpenBraces            = "T_OPEN_BRACES"
	OutputCloseBraces           = "T_CLOSE_BRACES"
	OutputSingleLineComment     = "T_SINGLE_LINE_COMMENT"
	OutputOpenMultilineComment  = "T_OPEN_MULTILINE_COMMENT"
	OutputCloseMultilineComment = "T_CLOSE_MULTILINE_COMMENT"

	//   Operators

	// --- Comparison operators

	OutputGreaterThanOperator  = "T_GREATER_THAN_OPERATOR"
	OutputGreaterEqualOperator = "T_GREATER_EQUAL_OPERATOR"
	OutputLessThanOperator     = "T_LESS_THAN_OPERATOR"
	OutputLessEqualOperator    = "T_LESS_EQUAL_OPERATOR"
	OutputEqualOperator        = "T_EQUAL_OPERATOR"
	OutputNotEqualOperator     = "T_NOT_EQUAL_OPERATOR"

	// --- Math operators

	OutputAdditionOperator       = "T_ADDITION_OPERATOR"
	OutputSubtractionOperator    = "T_SUBTRACTION_OPERATOR"
	OutputMultiplicationOperator = "T_MULTIPLICATION_OPERATOR"
	OutputDivisionOperator       = "T_DIVISION_OPERATOR"
	OutputModuleOperator         = "T_MODULE_OPERATOR"

	// --- Logical operators

	OutputNotOperator = "T_NOT_OPERATOR"
	OutputAndOperator = "T_AND_OPERATOR"
	OutputOrOperator  = "T_OR_OPERATOR"

	// --- Value operators

	OutputDeclarationOperator = "T_DECLARATION_OPERATOR"
	OutputAttributionOperator = "T_ATTRIBUTION_OPERATOR"
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
		RdInput:       bufio.NewReader(inputFile),
		CurrentLine:   0,
		CurrentColumn: 0,
		Pointer:       0,
		InputLine:     "",
		Token:         TNilValue,
		ErrorMessage:  "",
	}
	return lex
}

func (lex *Lexical) MovelookAhead() {
	if lex.Pointer+1 > len(lex.InputLine) {
		lex.CurrentLine++
		lex.Pointer = 0
		line, err := lex.RdInput.ReadString('\n')
		if err != nil {
			lex.LookAhead = TInputEnd
			return
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
	return
}

func (lex *Lexical) NextToken() {
	var sbLexeme strings.Builder

	for lex.LookAhead == ' ' || lex.LookAhead == '\t' || lex.LookAhead == '\n' || lex.LookAhead == '\r' {
		lex.MovelookAhead()
	}

	if lex.LookAhead >= 'A' && lex.LookAhead <= 'Z' {
		lex.alphabeticalCharacter(sbLexeme)
	} else if lex.LookAhead >= '0' && lex.LookAhead <= '9' {
		lex.numericalCharacter(sbLexeme)
	} else {
		lex.uniqueSymbolCharacter(sbLexeme)
	}

	lex.Lexeme = sbLexeme.String()
}

// ---------------------------------------------------------------------------------------------------------------------

func (lex *Lexical) alphabeticalCharacter(sbLexeme strings.Builder) {
	sbLexeme.WriteRune(lex.LookAhead)
	lex.MovelookAhead()

	for (lex.LookAhead >= 'A' && lex.LookAhead <= 'Z') || (lex.LookAhead >= '0' && lex.LookAhead <= '9') || lex.LookAhead == '_' {
		sbLexeme.WriteRune(lex.LookAhead)
		lex.MovelookAhead()
	}

	lex.Lexeme = sbLexeme.String()

	switch strings.ToUpper(lex.Lexeme) {
	// Construction tokens
	case Construct:
		lex.Token = TConstruct
	case Architect:
		lex.Token = TArchitect
	case Integrate:
		lex.Token = TIntegrate
	// Conditional and repetition tokens
	case If:
		lex.Token = TIf
	case Else:
		lex.Token = TElse
	case Elif:
		lex.Token = TElif
	case For:
		lex.Token = TFor
	case Detach:
		lex.Token = TDetach
	case Nil:
		lex.Token = TNil
	// Types
	case Gear:
		lex.Token = TGear
	case Tensor:
		lex.Token = TTensor
	case State:
		lex.Token = TState
	case Monodrone:
		lex.Token = TMonodrone
	case Omnidrone:
		lex.Token = TOmnidrone
	default:
		lex.Token = TId
	}
}

func (lex *Lexical) numericalCharacter(sbLexeme strings.Builder) {
	sbLexeme.WriteRune(lex.LookAhead)
	lex.MovelookAhead()
	for lex.LookAhead >= '0' && lex.LookAhead <= '9' {
		sbLexeme.WriteRune(lex.LookAhead)
		lex.MovelookAhead()
	}
	lex.Token = TGear
}

func (lex *Lexical) uniqueSymbolCharacter(sbLexeme strings.Builder) {
	switch lex.LookAhead {
	// Construction tokens
	case Comma:
		lex.Token = TComma
	case Colon:
		lex.Token = TColon
	// Structure tokens
	case OpenParentheses:
		lex.Token = TOpenParentheses
	case CloseParentheses:
		lex.Token = TCloseParentheses
	case OpenBraces:
		lex.Token = TOpenBraces
	case CloseBraces:
		lex.Token = TCloseBraces
	// Operators
	case GreaterThanOperator:
		lex.Token = TGreaterThanOperator
	case LessThanOperator:
		lex.Token = TLessThanOperator
	case AdditionOperator:
		lex.Token = TAdditionOperator
	case SubtractionOperator:
		lex.Token = TSubtractionOperator
	case MultiplicationOperator:
		lex.Token = TMultiplicationOperator
	case DivisionOperator:
		lex.Token = TDivisionOperator
	case ModuleOperator:
		lex.Token = TModuleOperator
	case NotOperator:
		lex.Token = TNotOperator
	case AttributionOperator:
		lex.Token = TAttributionOperator
	default:
		lex.Token = TLexError
		lex.ErrorMessage = fmt.Sprintf("Lexical error on line: %d\nRecognized upon reaching column: %d\nError line: <%s>\nUnknown token: %c", lex.CurrentLine, lex.CurrentColumn, lex.InputLine, lex.LookAhead)
	}
	sbLexeme.WriteRune(lex.LookAhead)
	lex.MovelookAhead()
}

func (lex *Lexical) multiSymbolCharacter(sbLexeme strings.Builder) {
	sbLexeme.WriteRune(lex.LookAhead)
	lex.MovelookAhead()

	for lex.LookAhead >= '&' && lex.LookAhead <= '/' {
		sbLexeme.WriteRune(lex.LookAhead)
		lex.MovelookAhead()
	}

	lex.Lexeme = sbLexeme.String()

	switch lex.Lexeme {
	// Construction tokens
	case SingleLineComment:
		lex.Token = TSingleLineComment
	case OpenMultilineComment:
		lex.Token = TOpenMultilineComment
	case CloseMultilineComment:
		lex.Token = TCloseMultilineComment
	// Conditional and repetition tokens
	case GreaterEqualOperator:
		lex.Token = TGreaterEqualOperator
	case LessEqualOperator:
		lex.Token = TLessEqualOperator
	case EqualOperator:
		lex.Token = TEqualOperator
	case NotEqualOperator:
		lex.Token = TNotEqualOperator
	case AndOperator:
		lex.Token = TAndOperator
	case OrOperator:
		lex.Token = TOrOperator
	case DeclarationOperator:
		lex.Token = TDeclarationOperator
	default:
		lex.Token = TLexError
		lex.ErrorMessage = fmt.Sprintf("Lexical error on line: %d\nRecognized upon reaching column: %d\nError line: <%s>\nUnknown token: %s", lex.CurrentLine, lex.CurrentColumn, lex.InputLine, lex.Lexeme)
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func (lex *Lexical) DisplayToken() {
	var tokenLexeme string
	switch lex.Token {
	// Construction tokens
	case TConstruct:
		tokenLexeme = OutputConstruct
	case TArchitect:
		tokenLexeme = OutputArchitect
	case TIntegrate:
		tokenLexeme = OutputIntegrate
	case TComma:
		tokenLexeme = OutputComma
	case TColon:
		tokenLexeme = OutputColon
	// Conditional and repetition
	case TIf:
		tokenLexeme = OutputIf
	case TElse:
		tokenLexeme = OutputElse
	case TElif:
		tokenLexeme = OutputElif
	case TFor:
		tokenLexeme = OutputFor
	case TDetach:
		tokenLexeme = OutputDetach
	// Type
	case TNil:
		tokenLexeme = OutputNil
	case TGear:
		tokenLexeme = OutputGear
	case TTensor:
		tokenLexeme = OutputTensor
	case TState:
		tokenLexeme = OutputState
	case TMonodrone:
		tokenLexeme = OutputMonodrone
	case TOmnidrone:
		tokenLexeme = OutputTensor
	case TId:
		tokenLexeme = OutputId
	// Structure
	case TOpenParentheses:
		tokenLexeme = OutputOpenParentheses
	case TCloseParentheses:
		tokenLexeme = OutputCloseParentheses
	case TOpenBraces:
		tokenLexeme = OutputOpenBraces
	case TCloseBraces:
		tokenLexeme = OutputCloseBraces
	case TSingleLineComment:
		tokenLexeme = OutputSingleLineComment
	case TOpenMultilineComment:
		tokenLexeme = OutputOpenMultilineComment
	case TCloseMultilineComment:
		tokenLexeme = OutputCloseMultilineComment
	// Operators
	case TGreaterThanOperator:
		tokenLexeme = OutputGreaterThanOperator
	case TGreaterEqualOperator:
		tokenLexeme = OutputGreaterEqualOperator
	case TLessThanOperator:
		tokenLexeme = OutputLessThanOperator
	case TLessEqualOperator:
		tokenLexeme = OutputLessEqualOperator
	case TNotEqualOperator:
		tokenLexeme = OutputNotEqualOperator
	case TAdditionOperator:
		tokenLexeme = OutputAdditionOperator
	case TSubtractionOperator:
		tokenLexeme = OutputSubtractionOperator
	case TMultiplicationOperator:
		tokenLexeme = OutputMultiplicationOperator
	case TDivisionOperator:
		tokenLexeme = OutputDivisionOperator
	case TModuleOperator:
		tokenLexeme = OutputModuleOperator
	case NotOperator:
		tokenLexeme = OutputNotOperator
	case TAndOperator:
		tokenLexeme = OutputAndOperator
	case TOrOperator:
		tokenLexeme = OutputOrOperator
	case TDeclarationOperator:
		tokenLexeme = OutputDeclarationOperator
	case TAttributionOperator:
		tokenLexeme = OutputAttributionOperator
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

func (lex *Lexical) WriteOutput() error {
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

func (lex *Lexical) ShowTokens() {
	fmt.Println("Identified Tokens (Token/Lexeme):")
	fmt.Println(lex.IdentifiedTokens.String())
}

func (lex *Lexical) storeTokens(tokenIdentificado string) {
	lex.IdentifiedTokens.WriteString(tokenIdentificado)
	lex.IdentifiedTokens.WriteString("\n")
}
