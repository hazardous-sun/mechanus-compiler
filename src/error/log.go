package custom_errors

import (
	"errors"
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

// customLog :
// Logs a message to the appropriate channel.
func customLog(message string, level string) {
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

// LogInfo :
// Logs an info.
func LogInfo(message string) {
	customLog(message, InfoLevel)
}

// LogWarning :
// Logs a warning.
func LogWarning(message string) {
	customLog(message, WarningLevel)
}

// LogError :
// Logs an error with proper unwrapping
func LogError(err error) {
	var analysisErr AnalysisError
	if errors.As(err, &analysisErr) {
		// Handle our custom error types
		customLog(analysisErr.Error(), ErrorLevel)
	} else {
		// Handle standard errors
		customLog(err.Error(), ErrorLevel)
	}
}

// LogSuccess :
// Logs a success message.
func LogSuccess(message string) {
	customLog(message, SuccessLevel)
}

// EnrichError :
// Enriches an error with extra context.
func EnrichError(err error, context string) error {
	return fmt.Errorf("%s -> %w", context, err)
}
