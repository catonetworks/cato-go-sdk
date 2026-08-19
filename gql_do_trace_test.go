package cato_go_sdk

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/hashicorp/terraform-plugin-log/tflogtest"
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

func TestParseGQLResponseMarksMalformedResponse(t *testing.T) {
	t.Parallel()

	var result struct{}
	err, parseFailed := parseGQLResponse(&clientv2.Client{}, []byte("not-json"), http.StatusOK, &result)
	if err == nil {
		t.Fatal("expected malformed response error")
	}
	if !parseFailed {
		t.Fatal("expected malformed response to be marked as a parse failure")
	}
}

func TestParseGQLResponseReturnsDataDecodeErrorWhenParsingDataWithErrors(t *testing.T) {
	t.Parallel()

	gqlc := &clientv2.Client{ParseDataWhenErrors: true}
	var result struct {
		Value int `json:"value"`
	}

	err, parseFailed := parseGQLResponse(
		gqlc,
		[]byte(`{"data":{"value":"not-an-integer"}}`),
		http.StatusOK,
		&result,
	)
	if err == nil {
		t.Fatal("expected data decode error")
	}
	if !parseFailed {
		t.Fatal("expected data decode error to be marked as a parse failure")
	}
}

func TestParseGQLResponseMarksInt64OverflowFromAccountSnapshot(t *testing.T) {
	t.Parallel()

	var result struct {
		AccountSnapshot struct {
			Timestamp int64 `json:"timestamp"`
		} `json:"accountSnapshot"`
	}
	err, parseFailed := parseGQLResponse(
		&clientv2.Client{},
		[]byte(`{"data":{"accountSnapshot":{"timestamp":18446744073709551613}}}`),
		http.StatusOK,
		&result,
	)
	if err == nil {
		t.Fatal("expected int64 overflow error")
	}
	if !parseFailed {
		t.Fatal("expected int64 overflow to be marked as a parse failure")
	}
	if !strings.Contains(err.Error(), "cannot unmarshal number 18446744073709551613") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseGQLResponseDoesNotMarkGraphQLErrorAsParseFailure(t *testing.T) {
	t.Parallel()

	var result struct{}
	err, parseFailed := parseGQLResponse(
		&clientv2.Client{},
		[]byte(`{"errors":[{"message":"request rejected"}]}`),
		http.StatusOK,
		&result,
	)
	if err == nil {
		t.Fatal("expected GraphQL error")
	}
	if parseFailed {
		t.Fatal("valid GraphQL error must not be marked as a parse failure")
	}
}

func TestLogResponseParseFailureIncludesOperationAndRenderedResponse(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	ctx := tflogtest.RootLogger(context.Background(), &output)
	logResponseParseFailure(
		ctx,
		"trace-1",
		`{"operationName":"accountSnapshot"}`,
		[]byte(`{"data":{"message":"line one\n\t\t\"quoted\""}}`),
	)

	entries, err := tflogtest.MultilineJSONDecode(&output)
	if err != nil {
		t.Fatalf("decode logs: %v", err)
	}
	if len(entries) < 2 {
		t.Fatalf("expected parse failure and response logs, got %d entries", len(entries))
	}
	if got := entries[0]["operation_name"]; got != "accountSnapshot" {
		t.Fatalf("unexpected operation name: %v", got)
	}

	var messages strings.Builder
	for _, entry := range entries[1:] {
		message, ok := entry["@message"].(string)
		if !ok {
			t.Fatalf("log message is not a string: %v", entry["@message"])
		}
		messages.WriteString(message)
		messages.WriteByte('\n')
	}
	rendered := messages.String()
	if strings.Contains(rendered, `\n`) || strings.Contains(rendered, `\t`) || strings.Contains(rendered, `\"`) {
		t.Fatalf("response contains unrendered JSON escapes: %q", rendered)
	}
	if !strings.Contains(rendered, "\t\t\"quoted\"") {
		t.Fatalf("response does not contain rendered whitespace and quotes: %q", rendered)
	}
}
