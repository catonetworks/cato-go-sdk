package cato_go_sdk

import (
	"context"

	"github.com/Yamashou/gqlgenc/clientv2"
	cato_models "github.com/catonetworks/cato-go-sdk/models"
)

const EntityLookupMinimalDocument = `query entityLookup ($accountID: ID!, $type: EntityType!, $limit: Int, $from: Int, $parent: EntityInput, $sort: [SortInput], $filters: [LookupFilterInput]) {
	entityLookup(accountID: $accountID, type: $type, limit: $limit, from: $from, parent: $parent, sort: $sort, filters: $filters) {
		items {
			entity {
				id
				name
				type
			}
			description
			helperFields
		}
		total
	}
}
`

func (c *Client) EntityLookupMinimal(ctx context.Context, accountID string, typeArg cato_models.EntityType, limit *int64, from *int64, parent *cato_models.EntityInput, sort []*cato_models.SortInput, filters []*cato_models.LookupFilterInput, interceptors ...clientv2.RequestInterceptor) (*EntityLookup, error) {
	vars := map[string]any{
		"accountID": accountID,
		"type":      typeArg,
		"limit":     limit,
		"from":      from,
		"parent":    parent,
		"sort":      sort,
		"filters":   filters,
	}

	var res EntityLookup
	if err := c.Client.Post(ctx, "entityLookup", EntityLookupMinimalDocument, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}
