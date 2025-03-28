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

	TConstruct   = 1
	TArchitect   = 2
	TIntegrate   = 3
	TComma       = 4
	TColon       = 5
	TSingleQuote = 6
	TDoubleQuote = 7

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
	TAndOperator            = 312
	TOrOperator             = 313
	TNotOperator            = 314
	TDeclarationOperator    = 315
	TAttributionOperator    = 316

	//	 Type tokens

	TNil       = 401
	TGear      = 402
	TTensor    = 403
	TState     = 404
	TMonodrone = 405
	TOmnidrone = 406
	TTypeName  = 407
	TId        = 408

	// Built-in functions

	TSend    = 501
	TReceive = 502

	//	 Control tokens

	TInputEnd = 601
	TLexError = 602
	TNilValue = 603
)

// Constants for Token symbols
const (
	//	 Construction tokens

	Construct    = "CONSTRUCT"
	Architect    = "ARCHITECT"
	Integrate    = "INTEGRATE"
	StringLexeme = "STRING"

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

	// Built-in functions

	Send    = "SEND"
	Receive = "RECEIVE"
)

// Constants for unique-symbol tokens
const (
	// Construction tokens

	Comma       = ','
	Colon       = ':'
	DoubleQuote = '"'
	SingleQuote = '\''

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
	OutputString    = "T_STRING"

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
	OutputTypeName  = "T_TYPE"
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

	// Built-in functions

	OutputSend    = "T_SEND"
	OutputReceive = "T_RECEIVE"
)

// Lexical struct to hold Lexical analyzer state
type Lexical struct {
	InputFile        *os.File
	Lines            []string
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
	CommentBlock     bool
}

func NewLexical(inputFile, outputFile *os.File) Lexical {
	lex := Lexical{
		InputFile:     inputFile,
		OutputFile:    outputFile,
		Lines:         make([]string, 0),
		CurrentLine:   0,
		CurrentColumn: 0,
		Pointer:       0,
		InputLine:     "",
		Token:         TNilValue,
		ErrorMessage:  "",
	}
	return lex
}

func (lex *Lexical) ReadLines() error {
	scanner := bufio.NewScanner(lex.InputFile)
	for scanner.Scan() {
		lex.Lines = append(lex.Lines, scanner.Text())
	}
	return scanner.Err()
}

func (lex *Lexical) MovelookAhead() error {
	// end of line reached
	if lex.Pointer+1 > len(lex.InputLine) {
		err := lex.nextLine()

		if err != nil {
			return err
		}

		if len(lex.InputLine) > 1 {
			lex.LookAhead = rune(lex.InputLine[lex.Pointer])
		} else {
			err := lex.MovelookAhead()
			if err != nil {
				return err
			}
		}
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

func (lex *Lexical) nextLine() error {
	lex.CurrentLine++
	lex.Pointer = 0
	if lex.CurrentLine < len(lex.Lines) {
		lex.InputLine = lex.Lines[lex.CurrentLine]
		return nil
	} else {
		custom_errors.Log(custom_errors.EndOfFileReached, nil, custom_errors.InfoLevel)
		return fmt.Errorf(custom_errors.EndOfFileReached)
	}
}

func (lex *Lexical) NextToken() error {
	var err error
	// Check if lex.LookAhead is inside a comment block
	if lex.CommentBlock {
		err = lex.skipComment()
	} else {
		for lex.isSeparatorCharacter() {
			err = lex.MovelookAhead()
			if err != nil {
				return err
			}
		}
	}

	if err != nil {
		return err
	}

	if lex.isAlphabeticalCharacter() {
		err = lex.alphabeticalCharacter()
	} else if lex.isNumericalCharacter() {
		err = lex.numericalCharacter()
	} else if lex.isString() {
		err = lex.quoteCharacters()
	} else {
		err = lex.symbolCharacter()
	}
	return err
}

func (lex *Lexical) symbolCharacter() error {
	temp := lex.LookAhead
	err := lex.MovelookAhead()
	if err != nil {
		return err
	}
	err = lex.multiSymbolCharacter(temp)
	if err != nil {
		return err
	}
	return nil
}

func (lex *Lexical) skipComment() error {
	for !lex.multilineCommentEnd() {
		err := lex.MovelookAhead()
		if err != nil {
			return err
		}
	}
	return nil
}

func (lex *Lexical) multilineCommentEnd() bool {
	// Checks that pointing to lex.Pointer+1 won't raise an index out of bound exception
	// AND
	// Checks if lex.LookAhead == '*'
	// AND
	// Checks if the current char + the next char == CloseMultilineComment
	if lex.Pointer+1 <= len(lex.InputLine) && lex.LookAhead == '*' {
		temp := fmt.Sprintf("%c%c", lex.LookAhead, lex.InputLine[lex.Pointer])
		if temp == CloseMultilineComment {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------------------------------------------------

func (lex *Lexical) isSeparatorCharacter() bool {
	return lex.LookAhead == ' ' || lex.LookAhead == '\t' || lex.LookAhead == '\n' || lex.LookAhead == '\r'
}

func (lex *Lexical) isAlphabeticalCharacter() bool {
	return (lex.LookAhead >= 'A' && lex.LookAhead <= 'Z') || (lex.LookAhead >= 'a' && lex.LookAhead <= 'z')
}

func (lex *Lexical) isNumericalCharacter() bool {
	return lex.LookAhead >= '0' && lex.LookAhead <= '9'
}

func (lex *Lexical) isString() bool {
	return lex.LookAhead == '"'
}

func (lex *Lexical) isMultiCharacterSymbol() bool {
	if matchesSingleCharSymbols(lex.LookAhead) {
		return false
	}
	return (lex.Pointer+1) < len(lex.InputLine) && (lex.InputLine[lex.Pointer+1] >= '&' && lex.InputLine[lex.Pointer+1] <= '/')
}

// ---------------------------------------------------------------------------------------------------------------------

func matchesSingleCharSymbols(lookAhead rune) bool {
	switch lookAhead {
	// Construction
	case Comma:
		return true
	case Colon:
		return true
	case DoubleQuote:
		return true
	case SingleQuote:
		return true
	// Structure
	case OpenParentheses:
		return true
	case CloseParentheses:
		return true
	case OpenBraces:
		return true
	case CloseBraces:
		return true
	default:
		return false
	}
}

func (lex *Lexical) alphabeticalCharacter() error {
	sbLexeme := strings.Builder{}
	sbLexeme.WriteRune(lex.LookAhead)
	err := lex.MovelookAhead()

	if err != nil {
		return err
	}

	for (lex.LookAhead >= 'A' && lex.LookAhead <= 'Z') || (lex.LookAhead >= '0' && lex.LookAhead <= '9') || lex.LookAhead == '_' {
		sbLexeme.WriteRune(lex.LookAhead)
		err = lex.MovelookAhead()
		if err != nil {
			return err
		}
	}

	lex.Lexeme = sbLexeme.String()

	switch lex.Lexeme {
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
		lex.Token = TTypeName
	// Types
	case Gear:
		lex.Token = TTypeName
	case Tensor:
		lex.Token = TTypeName
	case State:
		lex.Token = TTypeName
	case Monodrone:
		lex.Token = TTypeName
	case Omnidrone:
		lex.Token = TTypeName
	// Built-in functions
	case Send:
		lex.Token = TSend
	case Receive:
		lex.Token = TReceive
	default:
		lex.Token = TId
	}

	return nil
}

func (lex *Lexical) numericalCharacter() error {
	var err error
	sbLexeme := strings.Builder{}
	sbLexeme.WriteRune(lex.LookAhead)
	err = lex.MovelookAhead()

	if err != nil {
		return err
	}

	floatSeparatorFound := false

	for (lex.LookAhead >= '0' && lex.LookAhead <= '9') || (lex.LookAhead >= '.' && !floatSeparatorFound) {
		if lex.LookAhead == '.' {
			floatSeparatorFound = true
		}
		sbLexeme.WriteRune(lex.LookAhead)
		err = lex.MovelookAhead()
		if err != nil {
			return err
		}
	}

	lex.Lexeme = sbLexeme.String()

	if !floatSeparatorFound {
		lex.Token = TGear
	} else {
		lex.Token = TTensor
	}

	return err
}

func (lex *Lexical) uniqueSymbolCharacter(temp rune) {
	sbLexeme := strings.Builder{}
	sbLexeme.WriteRune(temp)

	switch temp {
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
	lex.Lexeme = sbLexeme.String()
}

func (lex *Lexical) multiSymbolCharacter(temp rune) error {
	var err error
	sbLexeme := strings.Builder{}
	sbLexeme.WriteRune(temp)

	for lex.LookAhead >= '&' && lex.LookAhead <= '/' {
		if lex.LookAhead == temp {
			sbLexeme.WriteRune(lex.LookAhead)
			err = lex.MovelookAhead()
		} else {
			sbLexeme.WriteRune(lex.LookAhead)
			err = lex.MovelookAhead()
			break
		}

		if err != nil {
			return err
		}
	}

	lex.Lexeme = sbLexeme.String()
	uniqueSymbol := false

	switch lex.Lexeme {
	// Construction tokens
	case SingleLineComment:
		lex.Token = TSingleLineComment
		// The lexical analyzer can jump to the next line because anything to the right of the single line comment
		// symbol, "//", should be ignored
		err = lex.nextLine()
	case OpenMultilineComment:
		lex.Token = TOpenMultilineComment
		lex.CommentBlock = true
	case CloseMultilineComment:
		lex.Token = TCloseMultilineComment
		lex.CommentBlock = false
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
		uniqueSymbol = true
		lex.uniqueSymbolCharacter(temp)
	}

	if err != nil {
		return err
	}

	if !uniqueSymbol {
		lex.Lexeme = sbLexeme.String()
	}

	return err
}

func (lex *Lexical) quoteCharacters() error {
	var err error
	char := lex.LookAhead
	sbLexeme := strings.Builder{}
	sbLexeme.WriteRune(lex.LookAhead)
	err = lex.MovelookAhead()

	if err != nil {
		return err
	}

	for lex.LookAhead != char {
		sbLexeme.WriteRune(lex.LookAhead)
		err = lex.MovelookAhead()

		if err != nil {
			return err
		}
	}

	sbLexeme.WriteRune(lex.LookAhead)
	err = lex.MovelookAhead()
	lex.Lexeme = sbLexeme.String()
	switch char {
	case DoubleQuote:
		lex.Token = TDoubleQuote
	case SingleQuote:
		lex.Token = TSingleQuote
	}
	return err
}

// ---------------------------------------------------------------------------------------------------------------------

func (lex *Lexical) DisplayToken() {
	var tokenLexeme string

	if lex.Token >= TConstruct && lex.Token < TIf {
		tokenLexeme = lex.displayConstructionToken()
	} else if lex.Token >= TIf && lex.Token < TOpenParentheses {
		tokenLexeme = lex.displayConditionalRepetitionToken()
	} else if lex.Token >= TOpenParentheses && lex.Token < TGreaterThanOperator {
		tokenLexeme = lex.displayStructureToken()
	} else if lex.Token >= TGreaterThanOperator && lex.Token <= TNil {
		tokenLexeme = lex.displayOperatorToken()
	} else if lex.Token >= TNil && lex.Token <= TSend {
		tokenLexeme = lex.displayTypeToken()
	} else {
		tokenLexeme = lex.displayFunctions()
	}

	fmt.Println(tokenLexeme + " ( " + lex.Lexeme + " )")
	lex.storeTokens(tokenLexeme + " ( " + lex.Lexeme + " )")
}

func (lex *Lexical) displayConstructionToken() string {
	switch lex.Token {
	// Construction tokens
	case TConstruct:
		return OutputConstruct
	case TArchitect:
		return OutputArchitect
	case TIntegrate:
		return OutputIntegrate
	case TComma:
		return OutputComma
	case TColon:
		return OutputColon
	case TSingleLineComment:
		return OutputMonodrone
	case TDoubleQuote:
		return OutputOmnidrone
	default:
		return "N/A"
	}
}

func (lex *Lexical) displayConditionalRepetitionToken() string {
	// Conditional and repetition
	switch lex.Token {
	case TIf:
		return OutputIf
	case TElse:
		return OutputElse
	case TElif:
		return OutputElif
	case TFor:
		return OutputFor
	case TDetach:
		return OutputDetach
	default:
		return "N/A"
	}
}

func (lex *Lexical) displayTypeToken() string {
	switch lex.Token {
	// Type
	case TNil:
		return OutputNil
	case TGear:
		return OutputGear
	case TTensor:
		return OutputTensor
	case TState:
		return OutputState
	case TMonodrone:
		return OutputMonodrone
	case TOmnidrone:
		return OutputOmnidrone
	case TTypeName:
		return OutputTypeName
	case TId:
		return OutputId
	default:
		return "N/A"
	}
}

func (lex *Lexical) displayStructureToken() string {
	switch lex.Token {
	// Structure
	case TOpenParentheses:
		return OutputOpenParentheses
	case TCloseParentheses:
		return OutputCloseParentheses
	case TOpenBraces:
		return OutputOpenBraces
	case TCloseBraces:
		return OutputCloseBraces
	case TSingleLineComment:
		return OutputSingleLineComment
	case TOpenMultilineComment:
		return OutputOpenMultilineComment
	case TCloseMultilineComment:
		return OutputCloseMultilineComment
	default:
		return "N/A"
	}
}

func (lex *Lexical) displayOperatorToken() string {
	switch lex.Token {
	// Operators
	case TGreaterThanOperator:
		return OutputGreaterThanOperator
	case TGreaterEqualOperator:
		return OutputGreaterEqualOperator
	case TLessThanOperator:
		return OutputLessThanOperator
	case TLessEqualOperator:
		return OutputLessEqualOperator
	case TEqualOperator:
		return OutputEqualOperator
	case TNotEqualOperator:
		return OutputNotEqualOperator
	case TAdditionOperator:
		return OutputAdditionOperator
	case TSubtractionOperator:
		return OutputSubtractionOperator
	case TMultiplicationOperator:
		return OutputMultiplicationOperator
	case TDivisionOperator:
		return OutputDivisionOperator
	case TModuleOperator:
		return OutputModuleOperator
	case TAndOperator:
		return OutputAndOperator
	case TOrOperator:
		return OutputOrOperator
	case TDeclarationOperator:
		return OutputDeclarationOperator
	case TAttributionOperator:
		return OutputAttributionOperator
	case TNotOperator:
		return OutputNotOperator
	default:
		return "N/A"
	}
}

func (lex *Lexical) displayFunctions() string {
	switch lex.Token {
	// Built-in functions
	case TSend:
		return OutputSend
	case TReceive:
		return OutputReceive
	default:
		return "N/A"
	}
}

// ---------------------------------------------------------------------------------------------------------------------

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

	custom_errors.Log(custom_errors.FileCloseSuccess, nil, custom_errors.SuccessLevel)
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

	custom_errors.Log(custom_errors.FileCreateSuccess, nil, custom_errors.SuccessLevel)
	return nil
}

func (lex *Lexical) ShowTokens() {
	custom_errors.Log(custom_errors.IdentifiedTokens, nil, custom_errors.SuccessLevel)
	fmt.Println(lex.IdentifiedTokens.String())
}

func (lex *Lexical) storeTokens(identifiedToken string) {
	lex.IdentifiedTokens.WriteString(identifiedToken)
	lex.IdentifiedTokens.WriteString("\n")
}
