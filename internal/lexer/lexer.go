package lexer

import (
	"bufio"
	"fmt"
	"mechanus-compiler/internal/compiler_error"
	"os"
	"strings"
)

// Lexer :
// This is the structure responsible for making the lexical analysis of the source file. It checks for unrecognized
// lexemes and, if it finds one, it returns an error code.
type Lexer struct {
	logger           *compiler_error.Logger
	inputFile        *os.File
	lines            []string
	outputFile       *os.File
	lookAhead        rune
	token            int
	lexeme           string
	pointer          int
	inputLine        string
	currentLine      int
	currentColumn    int
	errorMessage     error
	identifiedTokens strings.Builder
	commentBlock     bool
}

//**********************************************************************************************************************
// Public controllers
//**********************************************************************************************************************

// NewLexer :
// Initializes a new Lexer instance with the provided input and output files. It also sets up various initial values
// for the lexer.
//
// Fails if it is not possible to read the source file.
func NewLexer(inputFile, outputFile *os.File, debug bool) (Lexer, error) {
	// Initialize the logger. Log to Stderr. Set level based on the debug flag.
	logLevel := compiler_error.LevelInfo
	if debug {
		logLevel = compiler_error.LevelDebug
	}
	lg := compiler_error.New(os.Stderr, logLevel)

	// Initialize the structure
	lex := Lexer{
		logger:        lg,
		inputFile:     inputFile,
		outputFile:    outputFile,
		lines:         make([]string, 0),
		currentLine:   0,
		currentColumn: 0,
		pointer:       0,
		inputLine:     "",
		token:         TNilValue,
		errorMessage:  nil,
	}

	// Read the source file
	if err := lex.readLines(); err != nil {
		err = compiler_error.FileErrorf("NewLexer", err)
		// Use the new logger to log the error
		lex.logger.Error(err, map[string]any{"source": "NewLexer"})
		return Lexer{}, err
	}

	// Collect the first lexeme
	if err := lex.moveLookAhead(); err != nil {
		err = compiler_error.FileErrorf("NewLexer", err)
		// Use the new logger to log the error
		lex.logger.Error(err, map[string]any{"source": "NewLexer"})
		return Lexer{}, err
	}

	return lex, nil
}

// NextToken :
// Advances the lexer to the next token, checking for separators, alphabetical characters, numerical characters, string
// literals, or symbols.
func (lex *Lexer) NextToken() (int, error) {
	errSalt := "Lexer.NextToken"

	// Check if lex.lookAhead is inside a comment block
	if lex.commentBlock {
		if err := lex.skipComment(); err != nil {
			err = compiler_error.LexerErrorf(errSalt, err)
			lex.logger.Error(err, map[string]any{"source": errSalt})
			return -1, err
		}
	} else {
		for lex.isSeparatorCharacter() {
			if err := lex.moveLookAhead(); err != nil {
				err = compiler_error.LexerErrorf(errSalt, err)
				return -1, err
			}
		}
	}

	err := lex.collectLexeme()

	if lex.token == TSingleLineComment {
		_, err = lex.NextToken()
	}

	if err != nil {
		err = compiler_error.LexerErrorf(errSalt, err)
		lex.logger.Error(err, map[string]any{"source": errSalt})
		return -1, err
	}

	lex.logger.Debug("Token processed", map[string]any{"tokenID": lex.token})

	return lex.token, nil
}

// WIP :
// Checks if Lexer should keep working.
func (lex *Lexer) WIP() bool {
	return lex.token != TInputEnd && lex.token != TLexError
}

// DisplayToken :
// Displays the current token and lexeme to the output.
func (lex *Lexer) DisplayToken() {
	var tokenLexeme string
	lex.lexeme = reverse(lex.lexeme)
	tokenLexeme = lex.identifyDisplayToken()
	fmt.Println(tokenLexeme + " ( " + lex.lexeme + " )")
	lex.storeTokens(tokenLexeme + " ( " + lex.lexeme + " )")
}

// GetToken :
// Returns the current Token ID.
func (lex *Lexer) GetToken() int {
	return lex.token
}

// GetLexeme :
// Returns the current lexeme.
func (lex *Lexer) GetLexeme() string {
	return lex.lexeme
}

// DisplayPos :
// Returns the current line and column in a formatted string.
func (lex *Lexer) DisplayPos() string {
	return fmt.Sprintf("Line: %d, Column: %d", lex.currentLine+1, lex.currentColumn+1)
}

// GetPos :
// Returns an array with the current line and column.
func (lex *Lexer) GetPos() []int {
	return []int{lex.currentLine, lex.currentColumn}
}

// Close :
// Closes the specified file (either input or output).
func (lex *Lexer) Close(file string) {
	lex.logger.Debug("Closing file", map[string]any{"file_type": file})

	switch file {
	case "input":
		if err := lex.inputFile.Close(); err != nil {
			err = compiler_error.FileErrorf("Lexer.Close", err)
			lex.logger.Error(err, nil)
			return
		}
	case "output":
		if err := lex.outputFile.Close(); err != nil {
			err = compiler_error.FileErrorf("Lexer.Close", err)
			lex.logger.Error(err, nil)
			return
		}
	}

	lex.logger.Info(compiler_error.FileCloseSuccess, nil)
}

// Fail :
// Checks if the Lexer failed to reach EOF.
func (lex *Lexer) Fail() error {
	if lex.token == TLexError {
		return lex.errorMessage
	}
	return nil
}

// WriteOutput :
// Writes the identified tokens to the output file.
func (lex *Lexer) WriteOutput() error {
	errSalt := "(Lexer.WriteOutput)"

	if lex.outputFile == nil {
		err := compiler_error.FileErrorf(errSalt, fmt.Errorf(compiler_error.UninitializedFile))
		lex.logger.Error(err, nil)
		return err
	}

	file, err := os.Create("output.txt")

	if err != nil {
		err = compiler_error.FileErrorf(errSalt, err)
		lex.logger.Error(err, nil)
		return err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			err = compiler_error.FileErrorf(errSalt, err)
			lex.logger.Error(err, nil)
			return
		}
	}(file)

	if _, err = file.WriteString(lex.identifiedTokens.String()); err != nil {
		err = compiler_error.FileErrorf(errSalt, err)
		lex.logger.Error(err, nil)
		return err
	}

	lex.logger.Info(compiler_error.FileCreateSuccess, nil)
	return nil
}

// ShowTokens :
// Displays the list of identified tokens.
func (lex *Lexer) ShowTokens() {
	lex.logger.Debug(compiler_error.IdentifiedTokens, nil)
	fmt.Println(lex.identifiedTokens.String())
}

//**********************************************************************************************************************
// Internal controllers
//**********************************************************************************************************************

// ----- File handling -------------------------------------------------------------------------------------------------

// Reads all lines from source file and stores them inside lex.lines
//
// Fails if it is not possible to read the source file, or if the source file is empty.
func (lex *Lexer) readLines() error {
	scanner := bufio.NewScanner(lex.inputFile)

	for scanner.Scan() {
		lex.lines = append(lex.lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		err = compiler_error.FileErrorf("Lexer.readLines", err)
		lex.logger.Error(err, nil)
		return err
	}

	if len(lex.lines) == 0 {
		err := compiler_error.FileErrorf("Lexer.readLines", fmt.Errorf(compiler_error.EmptyFile))
		lex.logger.Error(err, nil)
		return err
	}

	lex.currentLine = len(lex.lines) - 1
	lex.inputLine = lex.lines[lex.currentLine]
	lex.currentColumn = len(lex.inputLine)
	lex.pointer = lex.currentColumn
	return nil
}

// Moves the pointer to the next character in the current line. If the end of the line is reached, it loads the next
// line.
//
// Fails when reaching EOF.
func (lex *Lexer) moveLookAhead() error {
	// end of line reached
	lex.pointer--

	// Check if the end of the line (right to left) was reached
	if lex.pointer < 0 {
		// Move the cursor up one line
		err := lex.nextLine()

		// Check if EOF was reached
		if err != nil {
			err = compiler_error.FileErrorf("Lexer.moveLookAhead", err)
			return err
		}

		// Check if the current line is not empty
		if len(lex.inputLine) >= 1 {
			lex.lookAhead = rune(lex.inputLine[lex.pointer])
		} else { // Move to the next line if it is
			err := lex.moveLookAhead()

			// Check if EOF was reached
			if err != nil {
				return err
			}
		}

	} else { // If the end of the line was not reached, collect the next character
		lex.currentColumn = lex.pointer + 1
		lex.lookAhead = rune(lex.inputLine[lex.pointer])
	}
	return nil
}

// Moves the cursor one line up
func (lex *Lexer) nextLine() error {
	// Move up one line
	lex.currentLine--

	// Check if the top of the file was reached
	if lex.currentLine <= 0 {
		lex.logger.Debug(compiler_error.EndOfFileReached, nil)
		return compiler_error.FileError(fmt.Errorf(compiler_error.EndOfFileReached))
	}

	// Collect the content of the line
	lex.inputLine = lex.lines[lex.currentLine]
	lex.pointer = len(lex.inputLine) - 1
	return nil
}

// Skips over a comment block until the end of the comment is reached.
func (lex *Lexer) skipComment() error {
	for !lex.multilineCommentEnd() {
		if err := lex.moveLookAhead(); err != nil {
			err = compiler_error.LexerErrorf("Lexer.skipComment", err)
			lex.logger.Error(err, nil)
			return err
		}
	}
	return nil
}

// Checks if the current position marks the end of a multiline comment.
func (lex *Lexer) multilineCommentEnd() bool {
	// Checks that pointing to lex.pointer+1 won't raise an index out of bound exception
	// AND
	// Checks if lex.lookAhead == '*'
	// AND
	// Checks if the current char + the next char == CloseMultilineComment
	if lex.pointer+1 <= len(lex.inputLine) && lex.lookAhead == '*' {
		temp := fmt.Sprintf("%c%c", lex.lookAhead, lex.inputLine[lex.pointer])
		if temp == CloseMultilineComment {
			return true
		}
	}
	return false
}

// Reverses a string. Used to output the correct lexeme
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ----- Lexeme identifiers --------------------------------------------------------------------------------------------

// Collects the newly found lexeme
func (lex *Lexer) collectLexeme() error {
	var err error

	if lex.isAlphabeticalCharacter() {
		err = lex.alphabeticalCharacter()
	} else if lex.isNumericalCharacter() {
		err = lex.numericalCharacter()
	} else if lex.isQuotation() {
		err = lex.quoteCharacters()
	} else {
		err = lex.symbolCharacter()
	}

	if err != nil {
		err = compiler_error.LexerErrorf("Lexer.collectLexeme", err)
		lex.logger.Error(err, nil)
		return err
	}

	return nil
}

// Checks if the current character is a separator (e.g., space, tab, newline).
func (lex *Lexer) isSeparatorCharacter() bool {
	return lex.lookAhead == ' ' || lex.lookAhead == '\t' || lex.lookAhead == '\r'
}

// Checks if the current character is an alphabetical letter (A-Z or a-z).
func (lex *Lexer) isAlphabeticalCharacter() bool {
	return (lex.lookAhead >= 'A' && lex.lookAhead <= 'Z') || (lex.lookAhead >= 'a' && lex.lookAhead <= 'z')
}

// Checks if the current character is a numerical digit (0-9).
func (lex *Lexer) isNumericalCharacter() bool {
	return lex.lookAhead >= '0' && lex.lookAhead <= '9'
}

// Checks if the current character is a quote (single or double).
func (lex *Lexer) isQuotation() bool {
	return lex.lookAhead == SingleQuote || lex.lookAhead == DoubleQuote
}

// Checks if the current character could be part of a multi-character symbol like an operator.
func (lex *Lexer) isMultiCharacterSymbol() bool {
	if matchesSingleCharSymbols(lex.lookAhead) {
		return false
	}
	return (lex.pointer+1) < len(lex.inputLine) && (lex.inputLine[lex.pointer+1] >= '&' && lex.inputLine[lex.pointer+1] <= '/')
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

// ----- Lexeme token ID generators ------------------------------------------------------------------------------------

// Processes alphabetical characters to form identifiers or keywords.
func (lex *Lexer) alphabeticalCharacter() error {
	sbLexeme := strings.Builder{}

	for (lex.lookAhead >= 'A' && lex.lookAhead <= 'Z') || (lex.lookAhead >= 'a' && lex.lookAhead <= 'z') || (lex.lookAhead >= '0' && lex.lookAhead <= '9') || lex.lookAhead == '_' {
		sbLexeme.WriteRune(lex.lookAhead)
		if err := lex.moveLookAhead(); err != nil {
			err = compiler_error.LexerErrorf("Lexer.alphabeticalCharacter", err)
			lex.logger.Error(err, nil)
			return err
		}
	}

	lex.lexeme = sbLexeme.String()

	switch reverse(strings.ToUpper(lex.lexeme)) {
	// Construction tokens
	case Construct:
		lex.token = TConstruct
	case Architect:
		lex.token = TArchitect
	case Integrate:
		lex.token = TIntegrate
	// Conditional and repetition tokens
	case If:
		lex.token = TIf
	case Else:
		lex.token = TElse
	case Elif:
		lex.token = TElif
	case For:
		lex.token = TFor
	case Detach:
		lex.token = TDetach
	case Nil:
		lex.token = TNil
	// Types
	case Gear:
		lex.token = TGear
	case Tensor:
		lex.token = TTensor
	case State:
		lex.token = TState
	case Monodrone:
		lex.token = TMonodrone
	case Omnidrone:
		lex.token = TOmnidrone
	// Built-in functions
	case Send:
		lex.token = TSend
	case Receive:
		lex.token = TReceive
	default:
		lex.token = TId
	}

	return nil
}

// Processes numerical characters and determines the type (Gear or Tensor).
func (lex *Lexer) numericalCharacter() error {
	sbLexeme := strings.Builder{}
	sbLexeme.WriteRune(lex.lookAhead)

	if err := lex.moveLookAhead(); err != nil {
		err = compiler_error.LexerErrorf("Lexer.numericalCharacter", err)
		lex.logger.Error(err, nil)
		return err
	}

	floatSeparatorFound := false

	for (lex.lookAhead >= '0' && lex.lookAhead <= '9') || (lex.lookAhead >= '.' && !floatSeparatorFound) {
		if lex.lookAhead == '.' {
			floatSeparatorFound = true
		}

		sbLexeme.WriteRune(lex.lookAhead)

		if err := lex.moveLookAhead(); err != nil {
			err = compiler_error.LexerErrorf("Lexer.numericalCharacter", err)
			lex.logger.Error(err, nil)
			return err
		}
	}

	lex.lexeme = sbLexeme.String()

	if !floatSeparatorFound {
		lex.token = TGear
	} else {
		lex.token = TTensor
	}

	return nil
}

// Handles symbols like operators, delimiters, and comments.
func (lex *Lexer) symbolCharacter() error {
	temp := lex.lookAhead

	if err := lex.moveLookAhead(); err != nil {
		err = compiler_error.LexerErrorf("Lexer.symbolCharacter", err)
		lex.logger.Error(err, nil)
		return err
	}

	if err := lex.multiSymbolCharacter(temp); err != nil {
		err = compiler_error.LexerErrorf("Lexer.symbolCharacter", err)
		lex.logger.Error(err, nil)
		return err
	}

	return nil
}

// Handles multi-character symbols like operators and comments.
func (lex *Lexer) multiSymbolCharacter(temp rune) error {
	sbLexeme := strings.Builder{}
	sbLexeme.WriteRune(temp)

	uniqueSymbol := false

	if checkMultiSymbolMatch(temp, lex.lookAhead) {
		sbLexeme.WriteRune(lex.lookAhead)

		if err := lex.moveLookAhead(); err != nil {
			err = compiler_error.LexerErrorf("Lexer.multiSymbolCharacter", err)
			lex.logger.Error(err, nil)
			return err
		}
	}

	var err error
	lex.lexeme = sbLexeme.String()

	switch reverse(lex.lexeme) {
	// Construction tokens
	case SingleLineComment:
		lex.token = TSingleLineComment
		// The lexical analyzer can jump to the next line because anything to the right of the single line comment
		// symbol, "//", should be ignored
		err = lex.nextLine()
	case OpenMultilineComment:
		lex.token = TOpenMultilineComment
		lex.commentBlock = true
	case CloseMultilineComment:
		lex.token = TCloseMultilineComment
		lex.commentBlock = false
	// Conditional and repetition tokens
	case GreaterEqualOperator:
		lex.token = TGreaterEqualOperator
	case LessEqualOperator:
		lex.token = TLessEqualOperator
	case EqualOperator:
		lex.token = TEqualOperator
	case NotEqualOperator:
		lex.token = TNotEqualOperator
	case AndOperator:
		lex.token = TAndOperator
	case OrOperator:
		lex.token = TOrOperator
	case DeclarationOperator:
		lex.token = TDeclarationOperator
	default:
		lex.uniqueSymbolCharacter(temp)
		uniqueSymbol = true
	}

	if err != nil {
		err = compiler_error.LexerErrorf("Lexer.multiSymbolCharacter", err)
		lex.logger.Error(err, nil)
		return err
	}

	if uniqueSymbol {
		lex.lexeme = sbLexeme.String()
	}

	return nil
}

func checkMultiSymbolMatch(char1, char2 rune) bool {
	sbLexeme := strings.Builder{}
	sbLexeme.WriteRune(char2)
	sbLexeme.WriteRune(char1)
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
func (lex *Lexer) uniqueSymbolCharacter(temp rune) {
	sbLexeme := strings.Builder{}
	sbLexeme.WriteRune(temp)

	switch temp {
	// Construction tokens
	case Comma:
		lex.token = TComma
	case Colon:
		lex.token = TColon
	// Structure tokens
	case OpenParentheses:
		lex.token = TOpenParentheses
	case CloseParentheses:
		lex.token = TCloseParentheses
	case OpenBraces:
		lex.token = TOpenBraces
	case CloseBraces:
		lex.token = TCloseBraces
	// Operators
	case GreaterThanOperator:
		lex.token = TGreaterThanOperator
	case LessThanOperator:
		lex.token = TLessThanOperator
	case AdditionOperator:
		lex.token = TAdditionOperator
	case SubtractionOperator:
		lex.token = TSubtractionOperator
	case MultiplicationOperator:
		lex.token = TMultiplicationOperator
	case DivisionOperator:
		lex.token = TDivisionOperator
	case ModuleOperator:
		lex.token = TModuleOperator
	case NotOperator:
		lex.token = TNotOperator
	case AttributionOperator:
		lex.token = TAttributionOperator
	default:
		lex.token = TLexError
		lex.errorMessage = fmt.Errorf("Lexer error on line: %d\nRecognized upon reaching column: %d\nError line: <%s>\nUnknown token: %c", lex.currentLine, lex.currentColumn, lex.inputLine, lex.lookAhead)
	}
	lex.lexeme = sbLexeme.String()
}

// Handles string literals, either single or double-quoted.
func (lex *Lexer) quoteCharacters() error {
	charCount := 0
	char := lex.lookAhead
	if char == '\'' {
		charCount = 1
	}
	sbLexeme := strings.Builder{}
	sbLexeme.WriteRune(lex.lookAhead)
	errSalt := "(Lexer.quoteCharacters)"

	if err := lex.moveLookAhead(); err != nil {
		err = compiler_error.LexerErrorf(errSalt, err)
		lex.logger.Error(err, nil)
		return err
	}

	for lex.lookAhead != char {
		if char == '\'' && charCount > 1 {
			return fmt.Errorf(compiler_error.InvalidMonodrone)
		}

		sbLexeme.WriteRune(lex.lookAhead)

		if err := lex.moveLookAhead(); err != nil {
			err = compiler_error.LexerErrorf(errSalt, err)
			lex.logger.Error(err, nil)
			return err
		}

		charCount++
	}

	sbLexeme.WriteRune(lex.lookAhead)

	if err := lex.moveLookAhead(); err != nil {
		err = compiler_error.LexerErrorf("Lexer.quoteCharacters", err)
		lex.logger.Error(err, nil)
		return err
	}

	lex.lexeme = sbLexeme.String()
	switch char {
	case DoubleQuote:
		lex.token = TDoubleQuote
	case SingleQuote:
		lex.token = TSingleQuote
	}
	return nil
}

// ----- Display methods -----------------------------------------------------------------------------------------------

func (lex *Lexer) identifyDisplayToken() string {
	if lex.token >= TConstruct && lex.token < TIf {
		return lex.displayConstructionToken()
	} else if lex.token >= TIf && lex.token < TOpenParentheses {
		return lex.displayConditionalRepetitionToken()
	} else if lex.token >= TOpenParentheses && lex.token < TGreaterThanOperator {
		return lex.displayStructureToken()
	} else if lex.token >= TGreaterThanOperator && lex.token <= TNil {
		return lex.displayOperatorToken()
	} else if lex.token >= TNil && lex.token < TSend {
		return lex.displayTypeToken()
	} else {
		return lex.displayFunctions()
	}
}

func (lex *Lexer) displayConstructionToken() string {
	switch lex.token {
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

func (lex *Lexer) displayConditionalRepetitionToken() string {
	// Conditional and repetition
	switch lex.token {
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

func (lex *Lexer) displayTypeToken() string {
	switch lex.token {
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

func (lex *Lexer) displayStructureToken() string {
	switch lex.token {
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
	case TNewLine:
		return OutputNewLine
	default:
		return "N/A"
	}
}

func (lex *Lexer) displayOperatorToken() string {
	switch lex.token {
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

func (lex *Lexer) displayFunctions() string {
	switch lex.token {
	// Built-in functions
	case TSend:
		return OutputSend
	case TReceive:
		return OutputReceive
	default:
		return "N/A"
	}
}

// ----- Helper methods ------------------------------------------------------------------------------------------------

// Stores an identified token into the identifiedTokens builder.
func (lex *Lexer) storeTokens(identifiedToken string) {
	lex.identifiedTokens.WriteString(identifiedToken)
	lex.identifiedTokens.WriteString("\n")
}
