package gqls

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// HTTPClient represents an HTTP client with custom headers
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{client: &http.Client{Timeout: timeout}}
}

// GetIntrospection retrieves graphql introspection json from the specified URL
func (c HTTPClient) GetIntrospection(reqURL, apiCode string, undocumented bool) ([]byte, error) {
	const applicationJSON = "application/json"
	headers := map[string]string{
		"X-Api-Key":    apiCode,
		"Content-Type": applicationJSON,
		"Accept":       applicationJSON,
	}
	if undocumented {
		reqURL += "?with_undocumented=true"
	}

	query := []byte(`{"query":"query IntrospectionQuery { __schema { description } }","operationName":"IntrospectionQuery"}`)
	data, _, err := c.SendRequest(http.MethodPost, reqURL, headers, query)
	if err != nil {
		return data, err
	}
	return data, nil
}

// SendRequest sends a request with custom headers
func (c HTTPClient) SendRequest(method, reqURL string, headers map[string]string, inputData []byte,
) (data []byte, statusCode int, err error) {
	fmt.Fprintf(os.Stderr, "URL=%s\n", reqURL)
	var inputBody io.Reader = http.NoBody
	// prepare request data reader
	if inputData != nil {
		inputBody = bytes.NewBuffer(inputData)
	}

	// create request
	req, err := http.NewRequest(method, reqURL, inputBody)
	if err != nil {
		return nil, 0, fmt.Errorf("creating http request to '%s': %v", reqURL, err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// send request
	resp, err := c.client.Do(req) // #nosec G704
	if err != nil {
		return nil, 0, fmt.Errorf("request to '%s' failed: %v", reqURL, err)
	}
	if resp == nil {
		return nil, 0, fmt.Errorf("request to '%s' failed: no response", reqURL)
	}
	defer func() {
		if errClose := resp.Body.Close(); errClose != nil {
			log.Printf("failed to close response body: %v", errClose)
		}
	}()

	// read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("request to '%s' failed: cannot read body: %v", reqURL, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return respBody, resp.StatusCode, fmt.Errorf("request to '%s' failed: status-code: %d\n%s", reqURL, resp.StatusCode, string(respBody))
	}

	return respBody, resp.StatusCode, nil
}
