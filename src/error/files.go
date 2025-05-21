package custom_errors

import "fmt"

const (
	NoSourceFile      = "no source file was provided"
	InvalidFileName   = "invalid file name"
	UninitializedFile = "uninitialized file"
	EmptyFile         = "empty file"
	FileCreateSuccess = "successfully created file"
	FileCreateError   = "unable to create file"
	FileOpenSuccess   = "successfully opened file"
	FileOpenError     = "unable to open file"
	FileCloseSuccess  = "successfully closed the file"
	FileCloseError    = "unable to close file"
	EndOfFileReached  = "end of file reached"
)

// FileError :
// Creates a new file-related error with a static message.
// This is the simpler version of FileErrorf for when no formatting is needed.
// The error will be wrapped with the ErrFile type for consistent error handling.
//
// Example usage:
// return FileError(InvalidFileName)
// return FileError(FileOpenError)
func FileError(msg string) error {
	return fmt.Errorf("%w: %s", ErrFile, msg)
}

// FileErrorf :
// Creates a new formatted file-related error with context.
// Wraps the error with ErrFile type while preserving the original error structure.
// Supports all standard fmt.Sprintf formatting verbs.
//
// Example usage:
// return FileErrorf(InvalidFileName)
// return FileErrorf("%s: %s", FileOpenError, filename)
// return FileErrorf("failed to read %s at offset %d", filename, offset)
func FileErrorf(format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", ErrFile, fmt.Sprintf(format, args...))
}
