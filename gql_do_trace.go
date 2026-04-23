package cato_go_sdk

//Implements the same flow as gqlgenc’s do (including gzip), runs the same parse/unmarshal logic as parseResponse / unmarshal, reads Trace_id (with fallbacks for canonical / trace-id style names), and on any parse-time error with a non-empty trace wraps it in APIError

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/Yamashou/gqlgenc/graphqljson"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const dumpDirVariable = "TF_API_DUMP_DIR"

// traceIDFromResponseHeader returns the server tracing ID from response headers.
// The Cato server sets "Trace_id" (see AbstractTracingPlugin#addTraceIdToResponse).
func traceIDFromResponseHeader(h http.Header) string {
	for _, key := range []string{"Trace_id", "Trace-Id"} {
		if v := h.Get(key); v != "" {
			return v
		}
	}
	for k, vv := range h {
		if strings.EqualFold(strings.ReplaceAll(k, "-", "_"), "trace_id") && len(vv) > 0 {
			return vv[0]
		}
	}
	return ""
}

// executeGQLWithTrace performs the GraphQL HTTP round trip like clientv2.Client.do,
// then attaches Trace_id from the response to any returned error.
func executeGQLWithTrace(ctx context.Context, gqlc *clientv2.Client, req *http.Request, _ *clientv2.GQLRequestInfo, res any) error {
	requestBody := requestBodyForError(req)

	resp, err := gqlc.Client.Do(req)
	if err != nil {
		return &APIError{
			Err:         fmt.Errorf("request failed: %w", err),
			RequestBody: requestBody,
		}
	}
	defer resp.Body.Close()

	traceID := traceIDFromResponseHeader(resp.Header)

	bodyReader := io.Reader(resp.Body)
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gr, gerr := gzip.NewReader(resp.Body)
		if gerr != nil {
			return fmt.Errorf("gzip decode failed: %w", gerr)
		}
		defer gr.Close()
		bodyReader = gr
	}

	body, err := io.ReadAll(bodyReader)
	recordCall(ctx, traceID, requestBody, string(body))
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	parseErr := parseGQLResponse(gqlc, body, resp.StatusCode, res)
	if parseErr == nil {
		return parseErr
	}

	var api *APIError
	if errors.As(parseErr, &api) && (api.TraceID != "" || api.RequestBody != "") {
		return parseErr
	}

	return &APIError{
		Err:         parseErr,
		TraceID:     traceID,
		RequestBody: requestBody,
	}
}

// gqlEnvelope mirrors clientv2.response.
type gqlEnvelope struct {
	Data   json.RawMessage `json:"data"`
	Errors json.RawMessage `json:"errors"`
}

func parseGQLResponse(gqlc *clientv2.Client, body []byte, httpCode int, result any) error {
	errResponse := &clientv2.ErrorResponse{}
	notOK := httpCode < 200 || 299 < httpCode
	if notOK {
		errResponse.NetworkError = &clientv2.HTTPError{
			Code:    httpCode,
			Message: fmt.Sprintf("Response body %s", string(body)),
		}
	}

	if err := gqlUnmarshalResponse(gqlc, body, result); err != nil {
		var gqlErr *clientv2.GqlErrorList
		if errors.As(err, &gqlErr) {
			errResponse.GqlErrors = &gqlErr.Errors
		} else if !notOK {
			return err
		}
	}

	if errResponse.HasErrors() {
		return errResponse
	}

	return nil
}

func gqlUnmarshalResponse(gqlc *clientv2.Client, data []byte, res any) error {
	resp := gqlEnvelope{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return fmt.Errorf("failed to decode data %s: %w", string(data), err)
	}

	var err error
	if len(resp.Errors) > 0 {
		err = &clientv2.GqlErrorList{}
		if e := json.Unmarshal(data, err); e != nil {
			return fmt.Errorf("faild to parse graphql errors. Response content %s - %w", string(data), e)
		}
		if !gqlc.ParseDataWhenErrors {
			return err
		}
	}

	if errData := graphqljson.UnmarshalData(resp.Data, res); errData != nil {
		if gqlc.ParseDataWhenErrors {
			return err
		}
		return fmt.Errorf("failed to decode data into response %s: %w", string(data), errData)
	}

	return err
}

func requestBodyForError(req *http.Request) string {
	if req == nil {
		return ""
	}

	if req.GetBody != nil {
		rc, err := req.GetBody()
		if err != nil {
			return ""
		}
		defer rc.Close()

		b, err := io.ReadAll(rc)
		if err != nil {
			return ""
		}
		return string(b)
	}

	if req.Body == nil {
		return ""
	}

	b, err := io.ReadAll(req.Body)
	if err != nil {
		return ""
	}
	req.Body = io.NopCloser(bytes.NewReader(b))
	return string(b)
}

// recordCall logs the request and response, and optionally dumps them to a file if TF_API_DUMP_DIR is set.
func recordCall(ctx context.Context, traceID, requestBody, responseBody string) {
	tflog.Debug(ctx, "API Call", map[string]any{"request": requestBody, "response": responseBody, "trace_id": traceID})
	dumpDir := os.Getenv(dumpDirVariable)
	if dumpDir == "" {
		return
	}
	if err := os.MkdirAll(dumpDir, 0o755); err != nil {
		log.Printf("failed to create API call dump directory '%s': %v", dumpDir, err)
		return
	}

	operationName := getOperationName(requestBody)
	filename := fmt.Sprintf("%s_%s.txt", time.Now().Format("20060102_150405.000"), operationName)
	appendToFile(filepath.Join(dumpDir, filename), traceID, requestBody, responseBody)
}

// getOperationName extracts the "operationName" field from the GraphQL request body for use in log filenames.
func getOperationName(requestBody string) string {
	var req struct {
		OperationName string `json:"operationName"`
	}
	if err := json.Unmarshal([]byte(requestBody), &req); err != nil || req.OperationName == "" {
		return "unknown"
	}
	return req.OperationName
}

// appendToFile opens (or creates) a file and appends the given details to it.
func appendToFile(filename, traceID, request, response string) {
	// Open file in append mode, create if not exists, write-only
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("failed to open file: %v", err))
	}
	defer file.Close()

	// Write string to file
	ts := time.Now().Format(time.DateTime)
	if _, err := file.WriteString(fmt.Sprintf("traceID: %s  [%s]\n%s\n%s", traceID, ts, request, response)); err != nil {
		panic(fmt.Sprintf("failed to write to file: %v", err))
	}
}
