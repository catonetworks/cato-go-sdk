package cato_go_sdk

import (
	"context"

	"github.com/gqlgo/gqlgenc/clientv2"
)

const EventsFeedIndexDocument = `query eventsFeed ($accountIDs: [ID!], $marker: String) {
	eventsFeed(accountIDs: $accountIDs, marker: $marker) {
		marker
		fetchedCount
	}
}
`

func (c *Client) EventsFeedIndex(ctx context.Context, accountIDs []string, marker *string, interceptors ...clientv2.RequestInterceptor) (*EventsFeed, error) {
	vars := map[string]any{
		"accountIDs": accountIDs,
		"marker":     marker,
	}

	var res EventsFeed
	if err := c.Client.Post(ctx, "eventsFeed", EventsFeedIndexDocument, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}
