package cato_go_sdk

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/Yamashou/gqlgenc/graphqljson"
)

type CatoError struct {
	msg     string
	err     error
	traceID string
}

func (e *CatoError) Error() string   { return e.msg }
func (e *CatoError) TraceID() string { return e.traceID }
func (e *CatoError) Unwrap() error   { return e.err }

func makeError(msg string, err error, traceID string) *CatoError {
	if traceID != "" {
		msg += "\nTrace-ID: " + traceID
	}
	return &CatoError{
		msg:     msg,
		err:     err,
		traceID: traceID,
	}
}

func (c *Client) do(_ context.Context, req *http.Request, _ *clientv2.GQLRequestInfo, res any) error {
	resp, err := c.Client.Client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("Content-Encoding") == "gzip" {
		resp.Body, err = gzip.NewReader(resp.Body)
		if err != nil {
			return fmt.Errorf("gzip decode failed: %w", err)
		}
	}

	traceID := resp.Header.Get("trace_id")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return makeError(fmt.Sprintf("failed to read response body: %v", err), err, traceID)
	}

	err = c.parseResponse(body, resp.StatusCode, res)
	if err != nil {
		return makeError(err.Error(), err, traceID)
	}
	return nil
}

func (c *Client) parseResponse(body []byte, httpCode int, result any) error {
	errResponse := &clientv2.ErrorResponse{}
	isOKCode := httpCode < 200 || 299 < httpCode
	if isOKCode {
		errResponse.NetworkError = &clientv2.HTTPError{
			Code:    httpCode,
			Message: fmt.Sprintf("Response body %s", string(body)),
		}
	}

	// some servers return a graphql error with a non OK http code, try anyway to parse the body
	if err := c.unmarshal(body, result); err != nil {
		var gqlErr *clientv2.GqlErrorList
		if errors.As(err, &gqlErr) {
			errResponse.GqlErrors = &gqlErr.Errors
		} else if !isOKCode {
			return err
		}
	}

	if errResponse.HasErrors() {
		return errResponse
	}

	return nil
}

// response is a GraphQL layer response from a handler.
type response struct {
	Data   json.RawMessage `json:"data"`
	Errors json.RawMessage `json:"errors"`
}

func (c *Client) unmarshal(data []byte, res any) error {
	resp := response{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return fmt.Errorf("failed to decode data %s: %w", string(data), err)
	}

	var err error
	if len(resp.Errors) > 0 {
		// try to parse standard graphql error
		err = &clientv2.GqlErrorList{}
		if e := json.Unmarshal(data, err); e != nil {
			return fmt.Errorf("faild to parse graphql errors. Response content %s - %w", string(data), e)
		}

		// if ParseDataWhenErrors is true, try to parse data as well
		if !c.Client.ParseDataWhenErrors {
			return err
		}
	}

	if errData := graphqljson.UnmarshalData(resp.Data, res); errData != nil {
		// if ParseDataWhenErrors is true, and we failed to unmarshal data, return the actual error
		if c.Client.ParseDataWhenErrors {
			return err
		}

		return fmt.Errorf("failed to decode data into response %s: %w", string(data), errData)
	}

	return err
}
