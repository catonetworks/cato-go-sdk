package cato_go_sdk

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestTraceIDFromResponseHeader(t *testing.T) {
	t.Parallel()
	h := http.Header{}
	h.Set("Trace_id", "abc-123")
	if got := traceIDFromResponseHeader(h); got != "abc-123" {
		t.Fatalf("Trace_id: got %q", got)
	}

	h2 := http.Header{}
	h2.Set("trace-id", "lower-456")
	if got := traceIDFromResponseHeader(h2); got != "lower-456" {
		t.Fatalf("trace-id alias: got %q", got)
	}
}

func TestAPIErrorIncludesRequestBody(t *testing.T) {
	t.Parallel()
	err := &APIError{
		Err:         errors.New("boom"),
		TraceID:     "trace-1",
		RequestBody: `{"query":"q","variables":{"id":"1"}}`,
	}

	msg := err.Error()
	if !strings.Contains(msg, "traceID: trace-1") {
		t.Fatalf("missing trace ID in message: %q", msg)
	}
	if !strings.Contains(msg, `requestBody: {"query":"q","variables":{"id":"1"}}`) {
		t.Fatalf("missing request body in message: %q", msg)
	}
	if got := RequestBodyFromError(err); got != `{"query":"q","variables":{"id":"1"}}` {
		t.Fatalf("unexpected request body from helper: %q", got)
	}
}
