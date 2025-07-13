package compiler_error

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
// Wraps an existing error with the ErrFile type to provide consistent error handling.
// This annotates the error with file-related context while preserving the original error chain.
//
// Example usage:
// return FileError(ErrSomething)
func FileError(err error) error {
	return fmt.Errorf("(%w) %w", ErrFile, err)
}

// FileErrorf :
// Wraps an existing error with additional context and the ErrFile type.
// This provides a structured way to annotate file-related errors while preserving
// the original error chain. The context string describes the error scenario.
//
// Example usage:
// return FileErrorf("caller function", ErrSomething)
func FileErrorf(context string, err error) error {
	return fmt.Errorf("(%s) %s -> %w", ErrFile, context, err)
}
