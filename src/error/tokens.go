package custom_errors

import "fmt"

const (
	InvalidMonodrone = "type Monodrone expects only 1 character"
)

// TokenErrorf :
// Creates a new token-related error
//
// Example usage:
// return TokenErrorf(InvalidMonodrone)
// return TokenErrorf("invalid token '%s'", token)
func TokenErrorf(format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", ErrToken, fmt.Sprintf(format, args...))
}
