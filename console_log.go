package veldt

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sync"
)

// SimpleConsoleLogger logs messages to the console as is, with no decoration
type SimpleConsoleLogger struct {
}

var (
	mu     = &sync.Mutex{}
	output = os.Stdout
)

func formatMessage(message string) []byte {
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "%s\n", message)
	return b.Bytes()
}

func logf(format string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	writer := bufio.NewWriter(output)
	defer writer.Flush()
	msg := fmt.Sprintf(format, args...)
	writer.Write(formatMessage(msg))
}

// Debugf logs a debug message to the console
func (l SimpleConsoleLogger) Debugf(format string, args ...interface{}) {
	logf(format, args...)
}

// Infof logs an info message to the console
func (l SimpleConsoleLogger) Infof(format string, args ...interface{}) {
	logf(format, args...)
}

// Warnf logs a warn message to the console
func (l SimpleConsoleLogger) Warnf(format string, args ...interface{}) {
	logf(format, args...)
}

// Errorf logs an error message to the console
func (l SimpleConsoleLogger) Errorf(format string, args ...interface{}) {
	logf(format, args...)
}
