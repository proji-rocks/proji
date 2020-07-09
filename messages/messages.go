//nolint:gochecknoglobals
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

// DisableColors disables colorful text output.
func DisableColors() {
	setNoColors()
	renderPrefixes()
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

// Infof prints a formatted information message.
func Infof(format string, args ...interface{}) {
	printPrefixedMessagef(defaultOutput, prefixInfo, format, args...)
}

// Successf prints a formatted success message.
func Successf(format string, args ...interface{}) {
	printPrefixedMessagef(defaultOutput, prefixSuccess, format, args...)
}

// Warningf prints a formatted warning message.
func Warningf(format string, args ...interface{}) {
	printPrefixedMessagef(defaultOutput, prefixWarning, format, args...)
}

// Errorf prints a formatted error message.
func Errorf(format string, err error, args ...interface{}) {
	format = formatErrorMessage(format, err)
	printPrefixedMessagef(defaultErrorOutput, prefixError, format, args...)
}

func formatErrorMessage(format string, err error) string {
	errorString := colorError("error") + "=" + err.Error()
	format = strings.TrimSpace(format)
	if len(format) < 1 {
		return errorString
	}
	return format + " " + errorString
}

func printPrefixedMessagef(out io.Writer, prefix, format string, args ...interface{}) {
	_, err := fmt.Fprintf(out, prefix+format+"\n", args...)
	if err != nil {
		log.Fatalf("failed to print message to %v.\n", out)
	}
}
