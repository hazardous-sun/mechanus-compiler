package custom_errors

import (
	"fmt"
	"log"
)

const (
	DebugLevel   = "debug"
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
	case DebugLevel:
		log.Println(fmt.Sprintf("%sdebug: %s %s", infoColor, message, defaultColor))
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

// LogDebug :
// Logs an info.
func LogDebug(message string) {
	customLog(message, DebugLevel)
}

// LogWarning :
// Logs a warning.
func LogWarning(message string) {
	customLog(message, WarningLevel)
}

// LogError :
// Logs an error with proper unwrapping
func LogError(err error) {
	customLog(err.Error(), ErrorLevel)
}

// LogSuccess :
// Logs a success message.
func LogSuccess(message string) {
	customLog(message, SuccessLevel)
}
