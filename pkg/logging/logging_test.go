package logging

import (
	"bytes"
	"context"
	"net/url"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()

	type args struct {
		debug  bool
		client bool
	}
	cases := []struct {
		name string
		args args
	}{
		{
			name: "client debug logger",
			args: args{debug: true, client: true},
		},
		{
			name: "client production logger",
			args: args{debug: false, client: true},
		},
		{
			name: "server debug logger",
			args: args{debug: true, client: false},
		},
		{
			name: "client production logger",
			args: args{debug: false, client: false},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var logger *zap.SugaredLogger
			if tc.args.client {
				logger = NewClientLogger(tc.args.debug)
			} else {
				logger = NewServerLogger(tc.args.debug)
			}

			if logger == nil {
				t.Fatal("Returned a nil logger")
			}

			if tc.args.debug {
				if !logger.Desugar().Core().Enabled(zap.DebugLevel) {
					t.Error("Debug logger not enabled at debug level")
				}
			} else {
				if !logger.Desugar().Core().Enabled(zap.InfoLevel) {
					t.Error("Production logger not enabled at info level")
				}
			}
		})
	}
}

func TestFromContext(t *testing.T) {
	t.Parallel()

	// Positive case; context has a logger.
	l1 := NewClientLogger(false)
	if l1 == nil {
		t.Fatal("NewClientLogger returned a nil logger")
	}

	ctx := WithLogger(context.Background(), l1)
	if ctx == nil {
		t.Fatal("WithLogger returned a nil context")
	}

	l2 := FromContext(ctx)
	if l2 == nil {
		t.Fatal("FromContext returned a nil logger")
	}

	// Comparing memory addresses should be sufficient to determine equality.
	if l1 != l2 {
		t.Error("FromContext returned a different logger")
	}

	// Context has no logger, returned logger should still not be nil.
	ctx = context.Background()
	l3 := FromContext(ctx)
	if l3 == nil {
		t.Fatal("FromContext returned a nil logger")
	}
}

// This is for testing purposes only. Stole it from:
//   - https://github.com/uber-go/zap/blob/v1.22.0/sink.go#L61
type nopCloserSink struct{ zapcore.WriteSyncer }

func (nopCloserSink) Close() error { return nil }

func TestVisualLevelEncoder(t *testing.T) {
	t.Parallel()

	// Custom memory sink to capture output.
	buf := bytes.NewBuffer(nil)
	memFactory := func(u *url.URL) (zap.Sink, error) {
		return nopCloserSink{zapcore.AddSync(buf)}, nil
	}

	err := zap.RegisterSink("mem", memFactory)
	if err != nil {
		t.Fatalf("register new zap sink: %v", zap.Error(err))
	}

	// Create a logger with the custom sink and the VisualLevelEncoder that we want to test.
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = VisualLevelEncoder
	config.OutputPaths = []string{"mem://"}

	logger, err := config.Build()
	if err != nil {
		t.Fatalf("build logger: %v", err)
	}

	// Log a message at each level and check the buffer for the expected output.
	// Debug
	logger.Debug("test")
	if !strings.Contains(buf.String(), symbolDebug) {
		t.Fatal("expected debug symbol in log message")
	}
	buf.Reset()

	// Info
	logger.Info("test")
	if !strings.Contains(buf.String(), symbolInfo) {
		t.Fatal("expected info symbol in log message")
	}
	buf.Reset()

	// Warning
	logger.Warn("test")
	if !strings.Contains(buf.String(), symbolWarn) {
		t.Fatal("expected warn symbol in log message")
	}
	buf.Reset()

	// Error
	logger.Error("test")
	if !strings.Contains(buf.String(), symbolError) {
		t.Fatal("expected error symbol in log message")
	}
	buf.Reset()
}
