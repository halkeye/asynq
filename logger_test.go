package asynq

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestWithStructuredLogger(t *testing.T) {
	var buf bytes.Buffer
	slogger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	cfg := Config{LogLevel: DebugLevel}
	WithStructuredLogger(slogger)(&cfg)

	if cfg.StructuredLogger != slogger {
		t.Fatal("StructuredLogger not set on Config")
	}

	logger := newLogger(cfg.Logger, cfg.StructuredLogger, cfg.LogLevel)
	logger.InfoContext(context.Background(), "test", "component", "server")

	var rec map[string]any
	if err := json.Unmarshal(buf.Bytes(), &rec); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if rec["msg"] != "test" {
		t.Errorf("msg = %v, want test", rec["msg"])
	}
}

func TestNewLoggerLegacyPath(t *testing.T) {
	logger := newLogger(nil, nil, InfoLevel)
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}
	logger.Info("legacy path ok")
}

func TestWithStructuredLoggerForScheduler(t *testing.T) {
	slogger := slog.Default()
	opts := &SchedulerOpts{}
	WithStructuredLoggerForScheduler(slogger)(opts)
	if opts.StructuredLogger != slogger {
		t.Fatal("StructuredLogger not set on SchedulerOpts")
	}
}
