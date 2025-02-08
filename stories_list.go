package cato_go_sdk

import (
	"context"

	"github.com/Yamashou/gqlgenc/clientv2"
	cato_models "github.com/catonetworks/cato-go-sdk/models"
)

func (c *Client) XdrStoriesList(ctx context.Context, storyInput cato_models.StoryInput, accountID string, interceptors ...clientv2.RequestInterceptor) (*Xdr, error) {
	vars := map[string]any{
		"storyInput": storyInput,
		"accountID":  accountID,
	}

	var res Xdr
	if err := c.Client.Post(ctx, "xdr", XdrStoriesDocumentList, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}

const XdrStoriesDocumentList = `query xdr ( $storyInput:StoryInput! $accountID:ID! ) {
	xdr ( accountID:$accountID  ) {
		stories ( input:$storyInput ) {
			paging {
				from 
				limit 
				total 
			}
			items {
				id 
				accountId 
				accountName 
				updatedAt 
				createdAt 
				}
			}
		}
	}	

`
