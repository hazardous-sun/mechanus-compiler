package custom_errors

import "fmt"

const (
	LexerSuccess     = "lexical analysis completed with no errors"
	LexerError       = "lexical analysis completed with an error"
	IdentifiedTokens = "Identified Tokens (token/lexeme):"
)

// LexerErrorf :
// Wraps an existing error with additional context and the ErrLexical type.
//
// Example usage:
// return LexerErrorf("caller function", ErrSomething)
func LexerErrorf(context string, err error) error {
	return fmt.Errorf("(%s) %s -> %w", ErrLexical, context, err)
}
