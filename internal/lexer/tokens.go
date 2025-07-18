package lexer

//**********************************************************************************************************************
// Token IDs
//**********************************************************************************************************************

const (
	//	 Construction tokens

	TConstruct = iota
	TArchitect
	TIntegrate
	TComma
	TColon
	TSingleQuote
	TDoubleQuote

	//	 Conditional and repetition tokens

	TIf
	TElse
	TElif
	TFor
	TDetach

	//	 Structure tokens

	TOpenParentheses
	TCloseParentheses
	TOpenBraces
	TCloseBraces
	TSingleLineComment
	TOpenMultilineComment
	TCloseMultilineComment
	TNewLine

	//	 Operator tokens

	TGreaterThanOperator
	TLessThanOperator
	TGreaterEqualOperator
	TLessEqualOperator
	TEqualOperator
	TNotEqualOperator
	TAdditionOperator
	TSubtractionOperator
	TMultiplicationOperator
	TDivisionOperator
	TModuleOperator
	TAndOperator
	TOrOperator
	TNotOperator
	TDeclarationOperator
	TAttributionOperator

	//	 Type tokens

	TNil
	TGear
	TTensor
	TState
	TMonodrone
	TOmnidrone
	TTypeName
	TId

	// Built-in functions

	TSend
	TReceive

	//	 Control tokens

	TInputEnd
	TLexError
	TNilValue
)

//**********************************************************************************************************************
// Token values
//**********************************************************************************************************************

// Keyword tokens
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

// Unique-symbol tokens
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

// Multi-symbol tokens
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
	DeclarationOperator  = "=:"
)

//**********************************************************************************************************************
// Token output values
//**********************************************************************************************************************

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
	OutputNewLine               = "T_NEW_LINE"

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
