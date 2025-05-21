package custom_errors

import "fmt"

const (
	LexicalSuccess   = "lexical analysis completed with no errors"
	LexicalError     = "lexical analysis completed with an error"
	IdentifiedTokens = "Identified Tokens (token/lexeme):"
)

// LexicalErrorf :
// Wraps an existing error with additional context and the ErrLexical type.
//
// Example usage:
// return LexicalErrorf("caller function", ErrSomething)
func LexicalErrorf(context string, err error) error {
	return fmt.Errorf("(%s) %s -> %w", ErrLexical, context, err)
}
