package veldt

const (
	// Error indicates error-level logging
	Error int = iota
	// Warn indicates warn-level logging
	Warn
	// Info indicates info-level logging
	Info
	// Debug indicated debug-level logging
	Debug
)

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

func getLogger (level int) Logger {
	if Error == level {
		if nil == errorLog {
			if nil == warnLog {
				if nil == infoLog {
					return debugLog
				}
				return infoLog
			}
			return warnLog
		}
		return errorLog
	} else if Warn == level {
		if nil == warnLog {
			if nil == infoLog {
				return debugLog
			}
			return infoLog
		}
		return warnLog
	} else if Info == level {
		if nil == infoLog {
			return debugLog
		}
		return infoLog
	} else if Debug == level {
		return debugLog
	}
	return nil
}

// Debugf logs to the debug log.
func Debugf(format string, args ...interface{}) {
	logger := getLogger(Debug)
	if nil != logger {
		logger.Debugf(format, args...)
	}
}

// Infof logs to the info log.
func Infof(format string, args ...interface{}) {
	logger := getLogger(Info)
	if nil != logger {
		logger.Infof(format, args...)
	}
}

// Warnf logs to the warn log.
func Warnf(format string, args ...interface{}) {
	logger := getLogger(Warn)
	if nil != logger {
		logger.Warnf(format, args...)
	}
}

// Errorf logs to the err log.
func Errorf(format string, args ...interface{}) {
	logger := getLogger(Error)
	if nil != logger {
		logger.Errorf(format, args...)
	}
}
