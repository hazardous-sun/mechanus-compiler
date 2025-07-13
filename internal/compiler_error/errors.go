package custom_errors

// AnalysisError :
// The base type for all custom errors in the compiler.
type AnalysisError string

func (e AnalysisError) Error() string { return string(e) }

// Error types
const (
	ErrFile    AnalysisError = "file error"
	ErrLexical AnalysisError = "lexical error"
	ErrSyntax  AnalysisError = "syntax error"
	ErrToken   AnalysisError = "token error"
)
