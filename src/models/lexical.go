package models

import (
	"bufio"
	"fmt"
	custom_errors2 "mechanus-compiler/src/error"
	"os"
	"strings"
)

// Lexical :
// This is the structure responsible for making the lexical analysis of the source file. It checks for unrecognized
// lexemes and, if it finds one, it returns an error code.
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

// NewLexical :
// Initializes a new Lexical instance with the provided input and output files. It also sets up various initial values
// for the lexer.
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

// ReadLines :
// Reads all lines from source file and stores them inside lex.Lines
func (lex *Lexical) ReadLines() error {
	scanner := bufio.NewScanner(lex.InputFile)

	for scanner.Scan() {
		lex.Lines = append(lex.Lines, scanner.Text())
	}

	err := scanner.Err()

	if err == nil {
		lex.CurrentLine = len(lex.Lines) - 1
		lex.InputLine = lex.Lines[lex.CurrentLine]
		lex.CurrentColumn = len(lex.InputLine)
		lex.Pointer = lex.CurrentColumn
	}

	return err
}

// MoveLookAhead :
// Moves the pointer to the next character in the current line. If the end of the line is reached, it loads the next
// line.
func (lex *Lexical) MoveLookAhead() error {
	// end of line reached
	lex.Pointer--
	if lex.Pointer < 0 {
		err := lex.nextLine()

		if err != nil {
			return err
		}

		if len(lex.InputLine) >= 1 {
			lex.LookAhead = rune(lex.InputLine[lex.Pointer])
		} else {
			err := lex.MoveLookAhead()
			if err != nil {
				return err
			}
		}
	} else {
		lex.CurrentColumn = lex.Pointer + 1
		lex.LookAhead = rune(lex.InputLine[lex.Pointer])
	}
	return nil
}

func (lex *Lexical) nextLine() error {
	lex.CurrentLine--
	if lex.CurrentLine >= 0 {
		lex.InputLine = lex.Lines[lex.CurrentLine]
		lex.Pointer = len(lex.InputLine) - 1
		return nil
	} else {
		custom_errors2.Log(custom_errors2.EndOfFileReached, nil, custom_errors2.InfoLevel)
		return fmt.Errorf(custom_errors2.EndOfFileReached)
	}
}

// NextToken :
// Advances the lexer to the next token, checking for separators, alphabetical characters, numerical characters, string
// literals, or symbols.
func (lex *Lexical) NextToken() error {
	var err error
	// Check if lex.LookAhead is inside a comment block
	if lex.CommentBlock {
		err = lex.skipComment()
	} else {
		for lex.isSeparatorCharacter() {
			err = lex.MoveLookAhead()
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
	} else if lex.isQuotation() {
		err = lex.quoteCharacters()
	} else {
		err = lex.symbolCharacter()
	}
	return err
}

// Handles symbols like operators, delimiters, and comments.
func (lex *Lexical) symbolCharacter() error {
	temp := lex.LookAhead
	err := lex.MoveLookAhead()
	if err != nil {
		return err
	}
	err = lex.multiSymbolCharacter(temp)
	if err != nil {
		return err
	}
	return nil
}

// Skips over a comment block until the end of the comment is reached.
func (lex *Lexical) skipComment() error {
	for !lex.multilineCommentEnd() {
		err := lex.MoveLookAhead()
		if err != nil {
			return err
		}
	}
	return nil
}

// Checks if the current position marks the end of a multiline comment.
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

// Checks if the current character is a separator (e.g., space, tab, newline).
func (lex *Lexical) isSeparatorCharacter() bool {
	return lex.LookAhead == ' ' || lex.LookAhead == '\t' || lex.LookAhead == '\n' || lex.LookAhead == '\r'
}

// Checks if the current character is an alphabetical letter (A-Z or a-z).
func (lex *Lexical) isAlphabeticalCharacter() bool {
	return (lex.LookAhead >= 'A' && lex.LookAhead <= 'Z') || (lex.LookAhead >= 'a' && lex.LookAhead <= 'z')
}

// Checks if the current character is a numerical digit (0-9).
func (lex *Lexical) isNumericalCharacter() bool {
	return lex.LookAhead >= '0' && lex.LookAhead <= '9'
}

// Checks if the current character is a quote (single or double).
func (lex *Lexical) isQuotation() bool {
	return lex.LookAhead == SingleQuote || lex.LookAhead == DoubleQuote
}

// Checks if the current character could be part of a multi-character symbol like an operator.
func (lex *Lexical) isMultiCharacterSymbol() bool {
	if matchesSingleCharSymbols(lex.LookAhead) {
		return false
	}
	return (lex.Pointer+1) < len(lex.InputLine) && (lex.InputLine[lex.Pointer+1] >= '&' && lex.InputLine[lex.Pointer+1] <= '/')
}

func matchesSingleCharSymbols(lookAhead rune) bool {
	switch lookAhead {
	// Construction
	case Comma:
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

// ---------------------------------------------------------------------------------------------------------------------

// Processes alphabetical characters to form identifiers or keywords.
func (lex *Lexical) alphabeticalCharacter() error {
	sbLexeme := strings.Builder{}

	for (lex.LookAhead >= 'A' && lex.LookAhead <= 'Z') || (lex.LookAhead >= 'a' && lex.LookAhead <= 'z') || (lex.LookAhead >= '0' && lex.LookAhead <= '9') || lex.LookAhead == '_' {
		sbLexeme.WriteRune(lex.LookAhead)
		err := lex.MoveLookAhead()
		if err != nil {
			return err
		}
	}

	lex.Lexeme = sbLexeme.String()

	switch reverse(strings.ToUpper(lex.Lexeme)) {
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

// Processes numerical characters and determines the type (Gear or Tensor).
func (lex *Lexical) numericalCharacter() error {
	var err error
	sbLexeme := strings.Builder{}
	sbLexeme.WriteRune(lex.LookAhead)
	err = lex.MoveLookAhead()

	if err != nil {
		return err
	}

	floatSeparatorFound := false

	for (lex.LookAhead >= '0' && lex.LookAhead <= '9') || (lex.LookAhead >= '.' && !floatSeparatorFound) {
		if lex.LookAhead == '.' {
			floatSeparatorFound = true
		}
		sbLexeme.WriteRune(lex.LookAhead)
		err = lex.MoveLookAhead()
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

// Handles multi-character symbols like operators and comments.
func (lex *Lexical) multiSymbolCharacter(temp rune) error {
	var err error
	sbLexeme := strings.Builder{}
	sbLexeme.WriteRune(temp)

	uniqueSymbol := false

	if checkMultiSymbolMatch(temp, lex.LookAhead) {
		sbLexeme.WriteRune(lex.LookAhead)
		err = lex.MoveLookAhead()

		if err != nil {
			return err
		}
	}

	lex.Lexeme = sbLexeme.String()

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
		lex.uniqueSymbolCharacter(temp)
		uniqueSymbol = true
	}

	if err != nil {
		return err
	}

	if uniqueSymbol {
		lex.Lexeme = sbLexeme.String()
	}

	return err
}

func checkMultiSymbolMatch(char1, char2 rune) bool {
	sbLexeme := strings.Builder{}
	sbLexeme.WriteRune(char1)
	sbLexeme.WriteRune(char2)
	symbol := sbLexeme.String()
	switch symbol {
	case SingleLineComment:
		return true
	case OpenMultilineComment:
		return true
	case CloseMultilineComment:
		return true
	case GreaterEqualOperator:
		return true
	case LessEqualOperator:
		return true
	case EqualOperator:
		return true
	case NotEqualOperator:
		return true
	case AndOperator:
		return true
	case OrOperator:
		return true
	case DeclarationOperator:
		return true
	default:
		return false
	}
}

// Processes single-character symbols and maps them to their respective token types.
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

// Handles string literals, either single or double-quoted.
func (lex *Lexical) quoteCharacters() error {
	var err error
	charCount := 0
	char := lex.LookAhead
	if char == '\'' {
		charCount = 1
	}
	sbLexeme := strings.Builder{}
	sbLexeme.WriteRune(lex.LookAhead)
	err = lex.MoveLookAhead()

	if err != nil {
		return err
	}

	for lex.LookAhead != char {
		if char == '\'' && charCount > 1 {
			return fmt.Errorf(custom_errors2.InvalidMonodrone)
		}
		sbLexeme.WriteRune(lex.LookAhead)
		err = lex.MoveLookAhead()

		if err != nil {
			return err
		}

		charCount++
	}

	sbLexeme.WriteRune(lex.LookAhead)
	err = lex.MoveLookAhead()
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

// DisplayToken :
// Displays the current token and lexeme to the output.
func (lex *Lexical) DisplayToken() {
	var tokenLexeme string
	lex.Lexeme = reverse(lex.Lexeme)

	if lex.Token >= TConstruct && lex.Token < TIf {
		tokenLexeme = lex.displayConstructionToken()
	} else if lex.Token >= TIf && lex.Token < TOpenParentheses {
		tokenLexeme = lex.displayConditionalRepetitionToken()
	} else if lex.Token >= TOpenParentheses && lex.Token < TGreaterThanOperator {
		tokenLexeme = lex.displayStructureToken()
	} else if lex.Token >= TGreaterThanOperator && lex.Token <= TNil {
		tokenLexeme = lex.displayOperatorToken()
	} else if lex.Token >= TNil && lex.Token < TSend {
		tokenLexeme = lex.displayTypeToken()
	} else {
		tokenLexeme = lex.displayFunctions()
	}

	fmt.Println(tokenLexeme + " ( " + lex.Lexeme + " )")
	lex.storeTokens(tokenLexeme + " ( " + lex.Lexeme + " )")
}

// ---------------------------------------------------------------------------------------------------------------------

// Reverses a string. Used to output the correct lexeme
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ---------------------------------------------------------------------------------------------------------------------

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
	case TSingleQuote:
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

// Close :
// Closes the specified file (either input or output).
func (lex *Lexical) Close(file string) {
	custom_errors2.Log(fmt.Sprintf("closing %s file", file), nil, custom_errors2.InfoLevel)

	switch file {
	case "input":
		err := lex.InputFile.Close()
		if err != nil {
			custom_errors2.Log(custom_errors2.FileCloseError, &err, custom_errors2.ErrorLevel)
			return
		}
	case "output":
		err := lex.OutputFile.Close()
		if err != nil {
			custom_errors2.Log(custom_errors2.FileCloseError, &err, custom_errors2.ErrorLevel)
			return
		}
	}

	custom_errors2.Log(custom_errors2.FileCloseSuccess, nil, custom_errors2.SuccessLevel)
}

// WriteOutput :
// Writes the identified tokens to the output file.
func (lex *Lexical) WriteOutput() error {
	if lex.OutputFile == nil {

		return fmt.Errorf(custom_errors2.UninitializedFile)
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

	custom_errors2.Log(custom_errors2.FileCreateSuccess, nil, custom_errors2.SuccessLevel)
	return nil
}

// ShowTokens :
// Displays the list of identified tokens.
func (lex *Lexical) ShowTokens() {
	custom_errors2.Log(custom_errors2.IdentifiedTokens, nil, custom_errors2.SuccessLevel)
	fmt.Println(lex.IdentifiedTokens.String())
}

// Stores an identified token into the IdentifiedTokens builder.
func (lex *Lexical) storeTokens(identifiedToken string) {
	lex.IdentifiedTokens.WriteString(identifiedToken)
	lex.IdentifiedTokens.WriteString("\n")
}
