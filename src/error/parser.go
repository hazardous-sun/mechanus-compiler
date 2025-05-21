package custom_errors

import "fmt"

const (
	SyntaxSuccess        = "syntax analysis completed with no errors"
	SyntaxError          = "syntax error"
	MissingConstructBody = "missing Construct body"
)

// SyntaxErrorf :
// Wraps an existing error with additional context and the ErrSyntax type.
//
// Example usage:
// return SyntaxErrorf("caller function", ErrSomething)
func SyntaxErrorf(context string, err error) error {
	return fmt.Errorf("(%s) %s -> %w", ErrSyntax, context, err)
}
