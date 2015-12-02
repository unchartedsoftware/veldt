package log

import (
	"path"
	"strings"
	"fmt"
	"github.com/Sirupsen/logrus"
	"runtime"
)

const (
	// InfoLevel logging is for high granularity development logging events.
	InfoLevel = 1
	// DebugLevel logging is for development level logging and common events.
	DebugLevel = 2
	// WarnLevel logging is for unexpected and recoverable events.
	WarnLevel = 3
	// ErrorLevel logging is for unexpected and unrecoverable fatal events.
	ErrorLevel = 4
)

var log = getLogger()

func getLogger() *logrus.Logger {
	log := logrus.New()
	log.Formatter = new(PrettyFormatter)
	log.Level = logrus.DebugLevel
	return log
}

func retrieveCallInfo() string {
	pc, file, line, _ := runtime.Caller(2)
	// get package name
	parts := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	lastIndex := len(parts) - 1
	index := 3 // domain/company/root
	if index > lastIndex {
		index = lastIndex
	}
	// remove function
	parts[lastIndex] = strings.Split(parts[lastIndex], ".")[0]
	packageName := strings.Join(parts[index:], "/")
	// get file name
	_, fileName := path.Split(file)
	return fmt.Sprint(packageName, "/", fileName, ":", line)
}

// Infof logging is for high granularity development logging events.
func Infof(format string, args ...interface{}) {
	log.WithFields(logrus.Fields{
		"fileinfo": retrieveCallInfo(),
	}).Infof(format, args...)
}

// Debugf logging is for development level logging events.
func Debugf(format string, args ...interface{}) {
	log.WithFields(logrus.Fields{
		"fileinfo": retrieveCallInfo(),
	}).Debugf(format, args...)
}

// Warnf logging is for unexpected and recoverable events.
func Warnf(format string, args ...interface{}) {
	log.WithFields(logrus.Fields{
		"fileinfo": retrieveCallInfo(),
	}).Warnf(format, args...)
}

// Errorf level is for unexpected and unrecoverable fatal events.
func Errorf(format string, args ...interface{}) {
	log.WithFields(logrus.Fields{
		"fileinfo": retrieveCallInfo(),
	}).Errorf(format, args...)
}

// Info logging is for high granularity development logging events.
func Info(args ...interface{}) {
	log.WithFields(logrus.Fields{
		"fileinfo": retrieveCallInfo(),
	}).Info(args...)
}

// Debug logging is for development level logging events.
func Debug(args ...interface{}) {
	log.WithFields(logrus.Fields{
		"fileinfo": retrieveCallInfo(),
	}).Debug(args...)
}

// Warn logging is for unexpected and recoverable events.
func Warn(args ...interface{}) {
	log.WithFields(logrus.Fields{
		"fileinfo": retrieveCallInfo(),
	}).Warn(args...)
}

// Error level is for unexpected and unrecoverable fatal events.
func Error(args ...interface{}) {
	log.WithFields(logrus.Fields{
		"fileinfo": retrieveCallInfo(),
	}).Error(args...)
}

// SetLevel sets the current logging output level.
func SetLevel(level int) {
	switch level {
	case InfoLevel:
		log.Level = logrus.InfoLevel
	case DebugLevel:
		log.Level = logrus.DebugLevel
	case WarnLevel:
		log.Level = logrus.WarnLevel
	case ErrorLevel:
		log.Level = logrus.ErrorLevel
	}
}
