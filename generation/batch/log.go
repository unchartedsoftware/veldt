package batch

import (
	"github.com/unchartedsoftware/veldt"
)

// logging utilities for the batch package, to provide clean and easy logging
// with simple visual tags with which to pick batch messages out of the log as
// a whole

var (
	debugLog veldt.Logger
	infoLog  veldt.Logger
	warnLog  veldt.Logger
	errorLog veldt.Logger
)

const (
	// Teal "BATCH" log prefix
	// Code derived from github.com/mgutz/ansi, but I wanted it as a const, so
	// couldn't se that directly
	preLog = "\033[1;38;5;6mBATCH\033[0m: "
)

// SetDebugLogger sets the debug level logger for the batch package
func SetDebugLogger(log veldt.Logger) {
	debugLog = log
}

// SetInfoLogger sets the info level logger for the batch package
func SetInfoLogger(log veldt.Logger) {
	infoLog = log
}

// SetWarnLogger sets the info level logger for the batch package
func SetWarnLogger(log veldt.Logger) {
	warnLog = log
}

// SetErrorLogger sets the info level logger for the batch package
func SetErrorLogger(log veldt.Logger) {
	errorLog = log
}

func getLogger(level int) veldt.Logger {
	if veldt.Error == level {
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
	} else if veldt.Warn == level {
		if nil == warnLog {
			if nil == infoLog {
				return debugLog
			}
			return infoLog
		}
		return warnLog
	} else if veldt.Info == level {
		if nil == infoLog {
			return debugLog
		}
		return infoLog
	} else if veldt.Debug == level {
		return debugLog
	}
	return nil
}

// Errorf logs to the error log
func Errorf(format string, args ...interface{}) {
	logger := getLogger(veldt.Error)
	if nil != logger {
		logger.Errorf(preLog+format, args...)
	} else {
		veldt.Errorf(preLog+format, args...)
	}
}

// Warnf logs to the warn log
func Warnf(format string, args ...interface{}) {
	logger := getLogger(veldt.Warn)
	if nil != logger {
		logger.Warnf(preLog+format, args...)
	} else {
		veldt.Warnf(preLog+format, args...)
	}
}

// Infof logs to the info log
func Infof(format string, args ...interface{}) {
	logger := getLogger(veldt.Info)
	if nil != logger {
		logger.Infof(preLog+format, args...)
	} else {
		veldt.Infof(preLog+format, args...)
	}
}

// Debugf logs to the debug log
func Debugf(format string, args ...interface{}) {
	logger := getLogger(veldt.Debug)
	if nil != logger {
		logger.Debugf(preLog+format, args...)
	} else {
		veldt.Debugf(preLog+format, args...)
	}
}
