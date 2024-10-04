package cato_go_sdk

import (
	"context"
	"net/http"

	"github.com/Yamashou/gqlgenc/clientv2"
)

// New function as wrapper to client
func New(url string, token string, httpClient http.Client, headers ...string) (*Client, error) {
	catoClient := &Client{
		Client: clientv2.NewClient(&httpClient, url, nil,
			func(ctx context.Context, req *http.Request, gqlInfo *clientv2.GQLRequestInfo, res interface{}, next clientv2.RequestInterceptorFunc) error {
				req.Header.Set("x-api-key", token)

				if len(headers) != 0 && len(headers)%2 == 0 {
					for i := 0; i < len(headers); i++ {
						req.Header.Set(headers[i], headers[i+1])
						i++
					}
				}

				return next(ctx, req, gqlInfo, res)
			}),
	}

	return catoClient, nil
}
