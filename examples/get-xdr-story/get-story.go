package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	cato "github.com/catonetworks/cato-go-sdk"
)

func main() {
	token := os.Getenv("CATO_API_KEY")
	accountId := os.Getenv("CATO_ACCOUNT_ID")
	storyId := os.Getenv("CATO_XDR_STORY_ID")
	incidentId := os.Getenv("CATO_XDR_INCIDENT_ID")
	url := "https://api.catonetworks.com/api/v1/graphql2"

	if token == "" {
		fmt.Println("no token provided")
		os.Exit(1)
	}

	if accountId == "" {
		fmt.Println("no account id provided")
		os.Exit(1)
	}

	if storyId == "" {
		fmt.Println("no story id provided")
		os.Exit(1)
	}

	catoClient, _ := cato.New(url, token, nil)

	ctx := context.Background()

	queryResult, err := catoClient.XdrStory(ctx, &storyId, &incidentId, accountId)
	if err != nil {
		fmt.Println("policy query error: ", err)
		return
	}

	queryResultJson, err := json.Marshal(queryResult)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(queryResultJson))

}
