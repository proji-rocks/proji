package messages

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gookit/color"
)

const (
	symbolInfo    = "• "
	symbolSuccess = "🗸 "
	symbolWarning = "⚡"
	symbolError   = "✗ "
)

var (
	defaultOutput      = os.Stdout
	defaultErrorOutput = os.Stderr

	defaultColorNone    = color.FgDefault.Render
	defaultColorInfo    = color.FgBlue.Render
	defaultColorSuccess = color.FgGreen.Render
	defaultColorWarning = color.FgYellow.Render
	defaultColorError   = color.FgRed.Render

	colorInfo    = defaultColorInfo
	colorSuccess = defaultColorSuccess
	colorWarning = defaultColorWarning
	colorError   = defaultColorError

	prefixInfo    = colorInfo(symbolInfo)
	prefixSuccess = colorSuccess(symbolSuccess)
	prefixWarning = colorWarning(symbolWarning)
	prefixError   = colorError(symbolError)
)

func EnableColors(disableColors bool) {
	if disableColors {
		setNoColors()
		renderPrefixes()
	}
}

func setNoColors() {
	colorInfo = defaultColorNone
	colorSuccess = defaultColorNone
	colorWarning = defaultColorNone
	colorError = defaultColorNone
}

func renderPrefixes() {
	prefixInfo = colorInfo(symbolInfo)
	prefixSuccess = colorSuccess(symbolSuccess)
	prefixWarning = colorWarning(symbolWarning)
	prefixError = colorError(symbolError)
}

// Info prints a formatted information message.
func Info(format string, args ...interface{}) {
	printPrefixedMessage(defaultOutput, prefixInfo, format, args...)
}

// Success prints a formatted success message.
func Success(format string, args ...interface{}) {
	printPrefixedMessage(defaultOutput, prefixSuccess, format, args...)
}

// Warning prints a formatted warning message.
func Warning(format string, args ...interface{}) {
	printPrefixedMessage(defaultOutput, prefixWarning, format, args...)
}

// Error prints a formatted error message.
func Error(format string, err error, args ...interface{}) {
	format = formatErrorMessage(format, err)
	printPrefixedMessage(defaultErrorOutput, prefixError, format, args...)
}

func formatErrorMessage(format string, err error) string {
	errorString := colorError("error") + "=" + err.Error()
	format = strings.TrimSpace(format)
	if len(format) < 1 {
		return errorString
	}
	return format + " " + errorString
}

func printPrefixedMessage(out io.Writer, prefix, format string, args ...interface{}) {
	_, err := fmt.Fprintf(out, prefix+format+"\n", args...)
	if err != nil {
		log.Fatalf("failed to print message to %v.\n", out)
	}
}
