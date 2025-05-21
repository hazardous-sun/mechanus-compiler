package custom_errors

import "fmt"

const (
	LexicalSuccess   = "lexical analysis completed with no errors"
	LexicalError     = "lexical analysis completed with an error"
	IdentifiedTokens = "Identified Tokens (token/lexeme):"
)

// LexicalErrorf :
// Creates a new lexical error with formatted message.
//
// Example usage:
// return LexicalErrorf("unknown character '%c' at position %d", char, pos)
func LexicalErrorf(format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", ErrLexical, fmt.Sprintf(format, args...))
}
