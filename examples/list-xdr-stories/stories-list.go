package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"

	cato "github.com/catonetworks/cato-go-sdk"
	cato_models "github.com/catonetworks/cato-go-sdk/models"
)

func main() {
	token := os.Getenv("CATO_API_KEY")
	accountId := os.Getenv("CATO_ACCOUNT_ID")
	url := "https://api.catonetworks.com/api/v1/graphql2"

	if token == "" {
		fmt.Println("no token provided")
		os.Exit(1)
	}

	if accountId == "" {
		fmt.Println("no account id provided")
		os.Exit(1)
	}

	catoClient, _ := cato.New(url, token, accountId, nil, nil)

	ctx := context.Background()

	xdrStoriesInput := cato_models.StoryInput{
		Filter: []*cato_models.StoryFilterInput{},
		Paging: &cato_models.PagingInput{
			From:  0,
			Limit: 80,
		},
		Sort: []*cato_models.StorySortInput{},
	}
	loopDone := false
	totalCount := int64(math.MaxInt64)

	queryInitialResult, err := catoClient.XdrStoriesList(ctx, xdrStoriesInput, accountId)
	if err != nil {
		fmt.Println("XdrStoriesList initial query error: ", err)
		return
	}

	if queryInitialResult.Xdr.Stories.Paging.Total < xdrStoriesInput.Paging.Limit {
		loopDone = true
	}

	for !loopDone {

		if xdrStoriesInput.Paging.From+xdrStoriesInput.Paging.Limit >= totalCount {
			xdrStoriesInput.Paging.Limit = totalCount - xdrStoriesInput.Paging.From
			loopDone = true
			break
		} else {
			xdrStoriesInput.Paging.From += xdrStoriesInput.Paging.Limit
		}

		queryResult, err := catoClient.XdrStoriesList(ctx, xdrStoriesInput, accountId)
		if err != nil {
			fmt.Println("XdrStoriesList loop query error: ", err)
			return
		}

		totalCount = queryResult.Xdr.Stories.Paging.Total

		queryInitialResult.Xdr.Stories.Items = append(queryInitialResult.Xdr.Stories.Items, queryResult.GetXdr().GetStories().GetItems()...)
	}

	queryResultJson, err := json.Marshal(queryInitialResult)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(queryResultJson))

}
