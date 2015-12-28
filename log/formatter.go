package log

import (
	"bytes"
	"fmt"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/mgutz/ansi"
)

var isTerminal = logrus.IsTerminal()

// PrettyFormatter is a formatter that meets the logrus.Formatter interface
type PrettyFormatter struct{}

// Format formats a logrus.Entry into an array of bytes
func (f *PrettyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	isColored := isTerminal && (runtime.GOOS != "windows")
	if isColored {
		printColored(b, entry)
	} else {
		printUncolored(b, entry)
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func getLevelString(level logrus.Level) string {
	switch level {
	case logrus.InfoLevel:
		return "  INFO   "
	case logrus.WarnLevel:
		return "  WARN   "
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return "  ERROR  "
	default:
		return "  DEBUG  "
	}
}

func getLevelColor(level logrus.Level) string {
	switch level {
	case logrus.InfoLevel:
		return ansi.Reset
	case logrus.WarnLevel:
		return ansi.Yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return ansi.Red
	default:
		return ansi.Blue
	}
}

func printColored(b *bytes.Buffer, entry *logrus.Entry) {
	levelText := getLevelString(entry.Level)
	levelColor := getLevelColor(entry.Level)
	// get fileinfo field
	fileinfo := entry.Data["fileinfo"]
	// write log message to buffer
	fmt.Fprintf(b, "%s[ %s ]%s %s[%s]%s %s %s(%s)%s",
		ansi.LightBlack,
		entry.Time.Format(time.Stamp),
		ansi.Reset,
		levelColor,
		levelText,
		ansi.Reset,
		entry.Message,
		ansi.Cyan,
		fileinfo,
		ansi.Reset)
}

func printUncolored(b *bytes.Buffer, entry *logrus.Entry) {
	levelText := getLevelString(entry.Level)
	// get fileinfo field
	fileinfo := entry.Data["fileinfo"]
	// write log message to buffer
	fmt.Fprintf(b, "[ %s ] [%s] %s (%s)",
		entry.Time.Format(time.Stamp),
		levelText,
		entry.Message,
		fileinfo)
}
