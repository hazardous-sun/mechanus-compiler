package custom_errors

import (
	"fmt"
	"log"
)

const (
	FileOpenSuccess  = "successfully opened file"
	FileCloseSuccess = "successfully closed the file"
	FileOpenError    = "unable to open file"
	FileCloseError   = "unable to close file"
)

const (
	InfoLevel    = "info"
	WarningLevel = "warning"
	ErrorLevel   = "error"
)

const (
	defaultColor = "\033[0m"
	errorColor   = "\033[91m"
	warningColor = "\033[93m"
)

func Log(message string, err *error, level string) {
	switch level {
	case InfoLevel:
		log.Println(fmt.Sprintf("%sinfo: %s %s", defaultColor, message, defaultColor))
	case WarningLevel:
		log.Println(fmt.Sprintf("%swarning: %s -> %v %s", warningColor, message, *err, defaultColor))
	case ErrorLevel:
		log.Println(fmt.Sprintf("%serror: %s -> %v %s", errorColor, message, *err, defaultColor))
	default:
		log.Println("invalid log level")
	}
}
