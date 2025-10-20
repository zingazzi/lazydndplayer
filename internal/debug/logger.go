// internal/debug/logger.go
package debug

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	logFile   *os.File
	logMutex  sync.Mutex
	isEnabled bool
)

// Enable enables debug logging to file
func Enable(filename string) error {
	var err error
	logFile, err = os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	isEnabled = true
	Log("Debug logging enabled")
	return nil
}

// IsEnabled returns whether debug logging is enabled
func IsEnabled() bool {
	return isEnabled
}

// Log writes a debug message to the log file with timestamp
func Log(format string, args ...interface{}) {
	if !isEnabled {
		return
	}

	logMutex.Lock()
	defer logMutex.Unlock()

	timestamp := time.Now().Format("15:04:05.000")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] %s\n", timestamp, message)

	if logFile != nil {
		logFile.WriteString(logLine)
		logFile.Sync() // Flush immediately for real-time debugging
	}
}

// Close closes the log file
func Close() {
	if logFile != nil {
		Log("Debug logging disabled")
		logFile.Close()
		logFile = nil
	}
	isEnabled = false
}
