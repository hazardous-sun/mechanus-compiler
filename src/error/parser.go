package custom_errors

import "fmt"

const (
	SyntaxSuccess        = "syntax analysis completed with no errors"
	SyntaxError          = "syntax error"
	MissingConstructBody = "missing Construct body"
)

// SyntaxErrorf :
// Creates a new syntax error with formatted message.
//
// Example usage:
// return SyntaxErrorf("unexpected token '%s'", token)
func SyntaxErrorf(format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", ErrSyntax, fmt.Sprintf(format, args...))
}
