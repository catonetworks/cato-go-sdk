package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	cato "github.com/catonetworks/cato-go-sdk"
	cato_models "github.com/catonetworks/cato-go-sdk/models"
)

func main() {
	token := os.Getenv("CATO_API_KEY")
	accountId := os.Getenv("CATO_ACCOUNT_ID")
	// siteId := os.Getenv("CATO_SITE_BGP_ID")
	url := "https://api.catonetworks.com/api/v1/graphql2"

	if token == "" {
		fmt.Println("no token provided")
		os.Exit(1)
	}

	if accountId == "" {
		fmt.Println("no account id provided")
		os.Exit(1)
	}

	catoClient, _ := cato.New(url, token, nil)

	ctx := context.Background()

	input := cato_models.SiteBgpStatusInput{
		Site: &cato_models.SiteRefInput{
			By:    "ID",
			Input: "15691",
		},
	}

	queryResult, rawResult, err := catoClient.SiteBgpStatus(ctx, input, accountId)
	if err != nil {
		fmt.Println("policy query error: ", err)
		return
	}

	for _, rawElement := range queryResult.Site.SiteBgpStatus.RawStatus {
		result := &cato.SiteBgpStatusResult{}
		err2 := json.Unmarshal([]byte(rawElement), &result)
		if err2 != nil {
			fmt.Println("Error parsing JSON:", err2)
			return
		}

		queryResultResultJson, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(queryResultResultJson))
	}

	queryResultJson, err := json.Marshal(rawResult)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(queryResultJson))
}
