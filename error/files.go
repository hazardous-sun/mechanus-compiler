package custom_errors

const (
	InvalidFIleName   = "invalid file name"
	UninitializedFile = "uninitialized file"
	EmptyFile         = "empty file"

	FileCreateSuccess = "successfully created file"
	FileCreateError   = "unable to create file"

	FileOpenSuccess = "successfully opened file"
	FileOpenError   = "unable to open file"

	FileCloseSuccess = "successfully closed the file"
	FileCloseError   = "unable to close file"

	EndOfFileReached = "end of file reached"
)
