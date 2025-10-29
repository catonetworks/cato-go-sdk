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
	url := os.Getenv("CATO_API_URL")
	siteId := os.Getenv("CATO_SITE_ID")

	if token == "" {
		fmt.Println("no token provided")
		os.Exit(1)
	}

	if accountId == "" {
		fmt.Println("no account id provided")
		os.Exit(1)
	}

	if siteId == "" {
		fmt.Println("no site id provided")
		os.Exit(1)
	}

	if url == "" {
		fmt.Println("no url provided")
		os.Exit(1)
	}

	catoClient, _ := cato.New(url, token, accountId, nil, nil)

	ctx := context.Background()

	policyChange, err := catoClient.SiteRemoveSite(ctx, siteId, accountId)

	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	queryResultJson, err := json.Marshal(policyChange)
	if err != nil {
		fmt.Println("SiteAddSocketSite error: ", err)
	}

	fmt.Println(string(queryResultJson))

}
