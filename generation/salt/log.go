package salt



import (
	"github.com/unchartedsoftware/plog"
)

// logging utilities for the salt package

const (
	// Yellow "SALT" log prefix
	// Code derived from github.com/mgutz/ansi, but I wanted it as a const, so
	// couldn't se that directly
	preLog = "\033[1;38;5;3mSALT\033[0m: "
	// And codes to make the message red, similarly as constants
	preMsg = "\033[1;97;3m"
	postMsg = "\033[0m"
)

func saltErrorf (format string, args ...interface{}) {
	log.SetLevel(log.WarnLevel)
	log.Errorf(preLog+format, args...)
}
func saltWarnf (format string, args ...interface{}) {
	log.SetLevel(log.WarnLevel)
	log.Warnf(preLog+format, args...)
}
func saltInfof (format string, args ...interface{}) {
	log.SetLevel(log.WarnLevel)
	log.Infof(preLog+format, args...)
}
func saltDebugf (format string, args ...interface{}) {
	log.SetLevel(log.WarnLevel)
	log.Debugf(preLog+format, args...)
}
