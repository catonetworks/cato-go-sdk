package cato_go_sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/hashicorp/go-retryablehttp"
)

var (
	// below is used with the hashicorp http retry client and
	// used to replace out the default CheckRetry function

	// Default retry configuration
	defaultRetryWaitMin = 3 * time.Second
	defaultRetryWaitMax = 30 * time.Second
	defaultRetryMax     = 5

	matchRateLimit      = regexp.MustCompile(`ratelimit`)
	matchRateDashLimit  = regexp.MustCompile(`rate-limit`)
	matchRateSpaceLimit = regexp.MustCompile(`rate limit`)
)

type RespErrors struct {
	Errors        []Errors `json:"errors"`
	NetworkErrors []Errors `json:"networkErrors"`
	GraphQLErrors []Errors `json:"graphqlErrors"`
}
type Errors struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

// New function as wrapper to client
func New(url string, token string, accountId string, httpClient *http.Client, headers map[string]string) (*Client, error) {

	// if an HTTP client is not provided, leverage the retry-enabled HTTP client
	// which allows for built-in support for rate limit and exponential backoff/retryconfiguration profile
	if httpClient == nil {
		retryClient := retryablehttp.NewClient()
		retryClient.RetryMax = defaultRetryMax
		retryClient.RetryWaitMin = defaultRetryWaitMin
		retryClient.RetryWaitMax = defaultRetryWaitMax
		retryClient.CheckRetry = baseRetryPolicy
		httpClient = retryClient.StandardClient()
	}

	catoClient := &Client{
		Client: clientv2.NewClient(httpClient, url, nil,
			func(ctx context.Context, req *http.Request, gqlInfo *clientv2.GQLRequestInfo, res any, next clientv2.RequestInterceptorFunc) error {
				req.Header.Set("x-api-key", token)
				req.Header.Set("x-account-id", accountId)
				for key, val := range headers {
					req.Header.Set(key, val)
				}

				// set a user-agent if not set by client
				_, exists := headers["User-Agent"]
				if !exists {
					req.Header.Set("User-Agent", "Cato_Go_SDK/v0")
				}

				return next(ctx, req, gqlInfo, res)
			}),
	}

	return catoClient, nil
}

func baseRetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {

	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	// in the event of connection errors, we do not get a response.
	if resp == nil {
		return false, fmt.Errorf("connection error")
	}

	// 429 Too Many Requests is recoverable. Sometimes the server puts
	// a Retry-After response header to indicate when the server is
	// available to start processing request from client.
	if resp.StatusCode == http.StatusTooManyRequests {
		log.Println("received http 429 response as http.StatusTooManyRequests")
		return true, nil
	}

	// rate limit errors could be in three locations within a 200 response:
	// - errors
	// - networkErrors
	// - graphqlErrors
	// these three locations can contain an array list of error messages so we need to read and
	// replace the response body to check for rate limit errors
	respErrors := &RespErrors{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &respErrors)

	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	if len(respErrors.Errors) > 0 {
		for _, respError := range respErrors.Errors {
			if matchRateLimit.MatchString(respError.Message) {
				log.Println("matched errors matchRateLimit: ", respError.Message)
				return true, nil
			}
			if matchRateDashLimit.MatchString(respError.Message) {
				log.Println("matched errors matchRateDashLimit: ", respError.Message)
				return true, nil
			}
			if matchRateSpaceLimit.MatchString(respError.Message) {
				log.Println("matched errors matchRateSpaceLimit: ", respError.Message)
				return true, nil
			}
		}
	}

	if len(respErrors.NetworkErrors) > 0 {
		for _, respError := range respErrors.NetworkErrors {
			if matchRateLimit.MatchString(respError.Message) {
				log.Println("matched networkErrors matchRateLimit: ", respError.Message)
				return true, nil
			}
			if matchRateDashLimit.MatchString(respError.Message) {
				log.Println("matched networkErrors matchRateDashLimit: ", respError.Message)
				return true, nil
			}
			if matchRateSpaceLimit.MatchString(respError.Message) {
				log.Println("matched networkErrors matchRateSpaceLimit: ", respError.Message)
				return true, nil
			}
		}
	}

	if len(respErrors.GraphQLErrors) > 0 {
		for _, respError := range respErrors.GraphQLErrors {
			if matchRateLimit.MatchString(respError.Message) {
				log.Println("matched graphqlErrors matchRateLimit: ", respError.Message)
				return true, nil
			}
			if matchRateDashLimit.MatchString(respError.Message) {
				log.Println("matched graphqlErrors matchRateDashLimit: ", respError.Message)
				return true, nil
			}
			if matchRateSpaceLimit.MatchString(respError.Message) {
				log.Println("matched graphqlErrors matchRateSpaceLimit: ", respError.Message)
				return true, nil
			}
		}
	}

	return false, nil
}
