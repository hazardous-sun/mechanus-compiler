package compiler_error

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

// LogLevel defines the severity of the log entry.
type LogLevel int

// Defines all available log levels.
const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
)

// String returns the string representation of a log level.
func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarning:
		return "WARNING"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger represents an active logging object that generates lines of JSON output to an io.Writer.
type Logger struct {
	out      io.Writer
	minLevel LogLevel
	mu       sync.Mutex
}

// New creates a new Logger. The out variable sets the destination to which log data will be written.
// The minLevel parameter sets the minimum level for logs to be written.
func New(out io.Writer, minLevel LogLevel) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

// print is an internal method that marshals the log entry and writes it to the output.
func (l *Logger) print(level LogLevel, message string, properties map[string]any) {
	if level < l.minLevel {
		return
	}

	// Get file and line number of the caller.
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	entry := struct {
		Time       string         `json:"time"`
		Level      string         `json:"level"`
		Message    string         `json:"message"`
		Properties map[string]any `json:"properties,omitempty"`
		File       string         `json:"file,omitempty"`
	}{
		Time:       time.Now().UTC().Format(time.RFC3339),
		Level:      level.String(),
		Message:    message,
		Properties: properties,
		File:       fmt.Sprintf("%s:%d", file, line),
	}

	// Use a mutex to ensure that log writes from different goroutines don't get interleaved.
	l.mu.Lock()
	defer l.mu.Unlock()

	lineBytes, err := json.Marshal(entry)
	if err != nil {
		// If marshaling fails, fall back to a plain text representation.
		lineBytes = []byte(fmt.Sprintf("%s %s: marshaling error: %v", LevelError.String(), time.Now().UTC().Format(time.RFC3339), err))
	}

	l.out.Write(append(lineBytes, '\n'))
}

// --- Public Logging Methods ---

func (l *Logger) Debug(message string, properties map[string]any) {
	l.print(LevelDebug, message, properties)
}

func (l *Logger) Info(message string, properties map[string]any) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) Warning(message string, properties map[string]any) {
	l.print(LevelWarning, message, properties)
}

func (l *Logger) Error(err error, properties map[string]any) {
	if err == nil {
		return
	}
	if properties == nil {
		properties = make(map[string]any)
	}
	properties["error"] = err.Error()
	l.print(LevelError, "an error occurred", properties)
}

func (l *Logger) Fatal(err error, properties map[string]any) {
	if err == nil {
		return
	}
	if properties == nil {
		properties = make(map[string]any)
	}
	properties["error"] = err.Error()
	l.print(LevelFatal, "a fatal error occurred", properties)
	os.Exit(1)
}
