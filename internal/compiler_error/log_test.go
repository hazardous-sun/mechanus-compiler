package compiler_error

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"
)

// TestLogger_Levels ensures that logs below the minimum level are suppressed.
func TestLogger_Levels(t *testing.T) {
	var buf bytes.Buffer

	// Create a logger that only logs INFO level and above.
	log := New(&buf, LevelInfo)

	// This debug message should be ignored.
	log.Debug("This is a debug message.", nil)
	if buf.Len() > 0 {
		t.Errorf("expected buffer to be empty after logging a debug message, but got: %s", buf.String())
	}

	// This info message should be logged.
	log.Info("This is an info message.", nil)
	if buf.Len() == 0 {
		t.Error("expected buffer to contain data after logging an info message, but it was empty")
	}

	// Check if the output contains the correct level string.
	if !strings.Contains(buf.String(), `"level":"INFO"`) {
		t.Errorf("expected log output to contain '\"level\":\"INFO\"', but got: %s", buf.String())
	}
}

// TestLogger_OutputFormat verifies the structure and content of the JSON output.
func TestLogger_OutputFormat(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, LevelDebug)

	// Define a custom error to log.
	testErr := errors.New("database connection failed")

	// Log an error with additional properties.
	log.Error(testErr, map[string]any{
		"userID":   123,
		"isActive": true,
	})

	// Define a struct to unmarshal the JSON log entry into.
	var entry struct {
		Level      string         `json:"level"`
		Message    string         `json:"message"`
		Properties map[string]any `json:"properties"`
		File       string         `json:"file"`
	}

	// Unmarshal the output from the buffer.
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to unmarshal log output: %v. Output was: %s", err, buf.String())
	}

	// Assert that the fields have the expected values.
	if entry.Level != "ERROR" {
		t.Errorf("expected level to be 'ERROR', but got '%s'", entry.Level)
	}

	if entry.Message != "an error occurred" {
		t.Errorf("expected message to be 'an error occurred', but got '%s'", entry.Message)
	}

	// Check the properties map for the original error string.
	if propErr, ok := entry.Properties["error"].(string); !ok || propErr != testErr.Error() {
		t.Errorf("expected properties to contain error '%s', but got '%v'", testErr.Error(), entry.Properties["error"])
	}

	// Check other custom properties. Note that JSON unmarshals numbers into float64 by default.
	if userID, ok := entry.Properties["userID"].(float64); !ok || userID != 123 {
		t.Errorf("expected userID to be 123, but got %v", entry.Properties["userID"])
	}

	// Check the file info.
	if !strings.Contains(entry.File, "log_test.go") {
		t.Errorf("expected file to be 'log_test.go', but got '%s'", entry.File)
	}
}

// TestLogger_Concurrency ensures that the logger is safe to use from multiple goroutines.
func TestLogger_Concurrency(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, LevelDebug)

	// Number of concurrent goroutines to spawn.
	numGoroutines := 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Start all goroutines to log simultaneously.
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			log.Info("Concurrent log message", nil)
		}()
	}

	// Wait for all logging to complete.
	wg.Wait()

	// The output should contain numGoroutines lines, each being a valid JSON object.
	// We split by newline and check each one.
	output := strings.TrimSpace(buf.String())
	lines := strings.Split(output, "\n")

	if len(lines) != numGoroutines {
		t.Errorf("expected %d log entries, but got %d", numGoroutines, len(lines))
	}

	for i, line := range lines {
		if !json.Valid([]byte(line)) {
			t.Errorf("log entry at line %d is not valid JSON: %s", i+1, line)
		}
	}
}

// logEntry is a helper struct for unmarshaling and verifying log output.
type logEntry struct {
	Level      string         `json:"level"`
	Message    string         `json:"message"`
	Properties map[string]any `json:"properties"`
	File       string         `json:"file"`
}

// TestLogger_NilError ensures that calling Error with a nil error produces no output.
func TestLogger_NilError(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, LevelDebug)

	// Calling Error with a nil error should be a no-op.
	log.Error(nil, map[string]any{"some_context": "value"})

	if buf.Len() != 0 {
		t.Errorf("expected buffer to be empty when logging a nil error, but got: %s", buf.String())
	}
}

// TestLogger_NilProperties tests that methods handle a nil properties map gracefully.
func TestLogger_NilProperties(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, LevelDebug)

	// Call Info with a nil map for properties.
	log.Info("test message with nil properties", nil)

	var entry logEntry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to unmarshal log output: %v", err)
	}

	// Because of `omitempty`, the 'properties' field should not be in the JSON output
	// or its value should be nil. A length check is a simple way to verify this.
	if len(entry.Properties) != 0 {
		t.Errorf("expected properties to be empty, but got: %v", entry.Properties)
	}
}

// TestLogger_ErrorWithNilProperties tests that the Error method correctly
// initializes a new properties map if one is not provided.
func TestLogger_ErrorWithNilProperties(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, LevelDebug)

	testErr := errors.New("internal server error")
	// Call Error with a nil properties map. The method should create one internally.
	log.Error(testErr, nil)

	var entry logEntry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to unmarshal log output: %v", err)
	}

	// The properties map should have been created and contain exactly one key: "error".
	if len(entry.Properties) != 1 {
		t.Errorf("expected exactly 1 property, but got %d", len(entry.Properties))
	}
	if errMsg, ok := entry.Properties["error"].(string); !ok || errMsg != testErr.Error() {
		t.Errorf("expected properties to contain the error string, but it didn't")
	}
}

// TestLogger_JsonMarshalFailure tests the logger's fallback mechanism
// when it fails to marshal the log entry to JSON.
func TestLogger_JsonMarshalFailure(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, LevelDebug)

	// Channels cannot be marshaled to JSON and will cause an error.
	unmarshalableProps := map[string]any{
		"bad_data": make(chan int),
	}

	log.Info("This message will fail to marshal", unmarshalableProps)

	output := buf.String()

	// Check if the output is the plain-text fallback message, not JSON.
	expectedFallback := "ERROR"
	if !strings.Contains(output, expectedFallback) {
		t.Errorf("expected fallback log to contain '%s', but got: %s", expectedFallback, output)
	}
	expectedError := "marshaling error"
	if !strings.Contains(output, expectedError) {
		t.Errorf("expected fallback log to contain '%s', but got: %s", expectedError, output)
	}
}
