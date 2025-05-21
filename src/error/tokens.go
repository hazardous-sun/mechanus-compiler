package custom_errors

import "fmt"

const (
	InvalidMonodrone = "type Monodrone expects only 1 character"
)

// TokenErrorf :
// Wraps an existing error with additional context and the ErrToken type.
//
// Example usage:
// return TokenErrorf("caller function", ErrSomething)
func TokenErrorf(context string, err error) error {
	return fmt.Errorf("(%s) %s -> %w", ErrToken, context, err)
}
