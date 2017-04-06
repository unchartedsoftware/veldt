package veldt

var (
	debugLog Logger
	infoLog  Logger
	warnLog  Logger
	errorLog Logger
)

// Logger represents a logger interface for tracing internal operations.
type Logger interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}

// SetDebugLogger sets the debug level logger.
func SetDebugLogger(log Logger) {
	debugLog = log
}

// SetInfoLogger sets the debug level logger.
func SetInfoLogger(log Logger) {
	infoLog = log
}

// SetWarnLogger sets the warn level logger.
func SetWarnLogger(log Logger) {
	warnLog = log
}

// SetErrorLogger sets the error level logger.
func SetErrorLogger(log Logger) {
	errorLog = log
}

// Debugf logs to the debug log.
func Debugf(format string, args ...interface{}) {
	if debugLog != nil {
		debugLog.Debugf(format, args...)
	}
}

// Infof logs to the info log.
func Infof(format string, args ...interface{}) {
	if infoLog != nil {
		infoLog.Infof(format, args...)
	}
}

// Warnf logs to the warn log.
func Warnf(format string, args ...interface{}) {
	if warnLog != nil {
		warnLog.Warnf(format, args...)
	}
}

// Errorf logs to the err log.
func Errorf(format string, args ...interface{}) {
	if errorLog != nil {
		errorLog.Errorf(format, args...)
	}
}
