package log

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/hibiken/asynq/internal/base"
	asynqcontext "github.com/hibiken/asynq/internal/context"
)

func TestStructuredLoggerDebugContext(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := NewStructuredLogger(slog.New(h), DebugLevel)

	logger.DebugContext(context.Background(), "hello", "component", "test")

	var rec map[string]any
	if err := json.Unmarshal(buf.Bytes(), &rec); err != nil {
		t.Fatalf("invalid json: %v; got %q", err, buf.String())
	}
	if rec["msg"] != "hello" {
		t.Errorf("msg = %v, want hello", rec["msg"])
	}
	if rec["component"] != "test" {
		t.Errorf("component = %v, want test", rec["component"])
	}
}

func TestStructuredLoggerRespectsLevel(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := NewStructuredLogger(slog.New(h), InfoLevel)

	logger.DebugContext(context.Background(), "hidden")
	if buf.Len() != 0 {
		t.Errorf("DebugContext should be suppressed at InfoLevel, got %q", buf.String())
	}
}

func TestStructuredLoggerTaskAttrs(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := NewStructuredLogger(slog.New(h), DebugLevel)

	msg := &base.TaskMessage{
		ID:    "task-123",
		Queue: "default",
	}
	ctx, cancel := asynqcontext.New(context.Background(), msg, time.Now().Add(time.Minute))
	defer cancel()

	logger.InfoContext(ctx, "processing")

	var rec map[string]any
	if err := json.Unmarshal(buf.Bytes(), &rec); err != nil {
		t.Fatalf("invalid json: %v; got %q", err, buf.String())
	}
	if rec["task_id"] != "task-123" {
		t.Errorf("task_id = %v, want task-123", rec["task_id"])
	}
	if rec["queue"] != "default" {
		t.Errorf("queue = %v, want default", rec["queue"])
	}
}

func TestLegacyLoggerContextFallback(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(newBase(&buf))
	logger.SetLevel(DebugLevel)

	logger.InfoContext(context.Background(), "hello", "key", "value")

	got := buf.String()
	if !strings.Contains(got, "INFO: hello key=value") {
		t.Errorf("got %q, want message with key=value", got)
	}
}
