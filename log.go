package veldt

const (
	// Error indicates error-level logging
	Error LogLevel = 4
	// Warn indicates warn-level logging
	Warn LogLevel = 3
	// Info indicates info-level logging
	Info LogLevel = 2
	// Debug indicated debug-level logging
	Debug LogLevel = 1
)

var (
	logger Logger
	level  LogLevel
)

// LogLevel represents a logging level.
type LogLevel int

// Logger represents a logger interface for tracing internal operations.
type Logger interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}

// SetLogger sets the logger and logging level.
func SetLogger(lvl LogLevel, log Logger) {
	logger = log
	level = lvl
}

// Debugf logs to the debug log.
func Debugf(format string, args ...interface{}) {
	if logger != nil && level <= Debug {
		logger.Debugf(format, args...)
	}
}

// Infof logs to the info log.
func Infof(format string, args ...interface{}) {
	if logger != nil && level <= Info {
		logger.Infof(format, args...)
	}
}

// Warnf logs to the warn log.
func Warnf(format string, args ...interface{}) {
	if logger != nil && level <= Warn {
		logger.Warnf(format, args...)
	}
}

// Errorf logs to the err log.
func Errorf(format string, args ...interface{}) {
	if logger != nil && level <= Error {
		logger.Errorf(format, args...)
	}
}
