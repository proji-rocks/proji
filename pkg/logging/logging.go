package logging

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	symbolDebug = "🐛"
	symbolInfo  = "💡"
	symbolWarn  = "⚠️"
	symbolError = "🔥"
)

var levelToPrefix = map[zapcore.Level]string{
	zap.DebugLevel: symbolDebug,
	zap.InfoLevel:  symbolInfo,
	zap.WarnLevel:  symbolWarn,
	zap.ErrorLevel: symbolError,
}

// VisualLevelEncoder is an zapcore.Encoder that encodes a zapcore.Level to a human-readable prefix.
func VisualLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(levelToPrefix[level])
}

func productionCLIEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		FunctionKey:      zapcore.OmitKey,
		LevelKey:         "L",
		MessageKey:       "M",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      VisualLevelEncoder,
		ConsoleSeparator: " ",
	}
}

func newLogger(debug, server bool) *zap.SugaredLogger {
	var config zap.Config

	if debug {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()

		// Special case for client production logging. Servers should use the default production config - json encoding
		if !server {
			config.Encoding = "console"
			config.EncoderConfig = productionCLIEncoderConfig()
		}
	}

	config.DisableStacktrace = true // Handled by cockroach/errors.

	logger, err := config.Build()
	if err != nil {
		logger = zap.NewNop()
	}

	return logger.Sugar()
}

// NewClientLogger returns a new logger that's meant to be used by clients. It always uses a human-readable format.
func NewClientLogger(debug bool) *zap.SugaredLogger {
	return newLogger(debug, false)
}

// NewServerLogger returns a new logger that's meant to be used by servers. It uses structured logging when run in
// production, and a human-readable format when run in development.
func NewServerLogger(debug bool) *zap.SugaredLogger {
	return newLogger(debug, true)
}

func defaultLogger() *zap.SugaredLogger {
	logger := NewClientLogger(false) // Or use nop here?
	logger = logger.Named("default") // Make it identifiable. Good for debugging.

	return logger
}

// As recommended by 'revive' linter.
type contextKey string

var loggerKey contextKey = "logger"

// WithLogger returns a new context.Context with the given logger.
func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext returns the logger from the given context. If the context does not contain a logger, the default logger
// is returned. If the context is nil, the default logger is returned.
func FromContext(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return logger
	}

	return defaultLogger()
}
