package citus

import (
	"github.com/unchartedsoftware/veldt"
)

var (
	logger veldt.Logger
	level  veldt.LogLevel
)

const (
	prefix = "CITUS: "
)

// Debugf logs to the debug log.
func Debugf(format string, args ...interface{}) {
	if logger != nil && level >= veldt.Debug {
		logger.Debugf(prefix+format, args...)
	} else {
		veldt.Debugf(prefix+format, args...)
	}
}

// Infof logs to the info log.
func Infof(format string, args ...interface{}) {
	if logger != nil && level >= veldt.Info {
		logger.Infof(prefix+format, args...)
	} else {
		veldt.Infof(prefix+format, args...)
	}
}

// Warnf logs to the warn log.
func Warnf(format string, args ...interface{}) {
	if logger != nil && level >= veldt.Warn {
		logger.Warnf(prefix+format, args...)
	} else {
		veldt.Warnf(prefix+format, args...)
	}
}

// Errorf logs to the err log.
func Errorf(format string, args ...interface{}) {
	if logger != nil && level >= veldt.Error {
		logger.Errorf(prefix+format, args...)
	} else {
		veldt.Errorf(prefix+format, args...)
	}
}
