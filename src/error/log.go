package custom_errors

import (
	"fmt"
	"log"
)

const (
	InfoLevel    = "info"
	WarningLevel = "warning"
	ErrorLevel   = "error"
	SuccessLevel = "success"
)

const (
	defaultColor = "\033[0m"
	infoColor    = "\033[36m"
	errorColor   = "\033[91m"
	warningColor = "\033[93m"
	successColor = "\033[32m"
)

// Log :
// Logs a message to stdout.
func Log(message string, level string) {
	switch level {
	case InfoLevel:
		log.Println(fmt.Sprintf("%sinfo: %s %s", infoColor, message, defaultColor))
	case WarningLevel:
		log.Println(fmt.Sprintf("%swarning: %s %s", warningColor, message, defaultColor))
	case ErrorLevel:
		log.Println(fmt.Sprintf("%serror: %s %s", errorColor, message, defaultColor))
	case SuccessLevel:
		log.Println(fmt.Sprintf("%ssuccess: %s %s", successColor, message, defaultColor))
	default:
		log.Println(fmt.Sprintf("%sinvalid log level%s '%s'%s -> %s", errorColor, warningColor, level, defaultColor, message))
	}
}

// EnrichError :
// Enriches an error with extra context.
func EnrichError(err error, msg string) error {
	return fmt.Errorf("%s: %v", msg, err)
}
