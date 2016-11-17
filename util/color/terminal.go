package color

import (
	"os"
	"regexp"
	"runtime"

	"github.com/mattn/go-isatty"
)

const (
	// Black output color code.
	Black = "\033[30m"
	// Red output color code.
	Red = "\033[31m"
	// Green output color code.
	Green = "\033[32m"
	// Yellow output color code.
	Yellow = "\033[33m"
	// Blue output color code.
	Blue = "\033[34m"
	// Magenta output color code.
	Magenta = "\033[35m"
	// Cyan output color code.
	Cyan = "\033[36m"
	// White output color code.
	White = "\033[37m"
	// Reset output color code.
	Reset = "\033[0m"
)

var (
	// ColorTerminal represents whether or not the terminal supports color.
	ColorTerminal = isatty.IsTerminal(os.Stdout.Fd()) && (runtime.GOOS != "windows")
	// regex to find any color patterns
	colorPattern = regexp.MustCompile(`\x1b\[\d{0,2}m`)
)

// RemoveColor will remove any color formatting from the provided text.
func RemoveColor(text string) string {
	return colorPattern.ReplaceAllString(text, "")
}
