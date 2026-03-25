package cato_go_sdk

import (
	"errors"
	"fmt"
)

// APIError represents an API call failure and includes the API trace ID
// when the server returns one (Trace_id response header), for root cause analysis and support.
// Clients created with New attach the trace ID to errors from failed GraphQL calls automatically.
type APIError struct {
	Err         error
	TraceID     string
	RequestBody string
}

func (e *APIError) Error() string {
	msg := e.Err.Error()
	if e.TraceID != "" {
		msg = fmt.Sprintf("%s (traceID: %s)", msg, e.TraceID)
	}
	if e.RequestBody != "" {
		msg = fmt.Sprintf("%s (requestBody: %s)", msg, e.RequestBody)
	}
	return msg
}

func (e *APIError) Unwrap() error {
	return e.Err
}

// TraceIDFromError returns the API trace ID from err if it is or wraps an APIError.
// Use this when handling failures to log or display the trace ID for support.
func TraceIDFromError(err error) string {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.TraceID
	}
	return ""
}

// RequestBodyFromError returns the GraphQL request body from err if it is or wraps an APIError.
func RequestBodyFromError(err error) string {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.RequestBody
	}
	return ""
}
