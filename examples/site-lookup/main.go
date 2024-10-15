package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	cato "github.com/catonetworks/cato-go-sdk"
	cato_models "github.com/catonetworks/cato-go-sdk/models"
)

func main() {
	token := os.Getenv("CATO_API_KEY")
	accountId := os.Getenv("CATO_ACCOUNT_ID")
	siteId := os.Getenv("CATO_SITE_ID")
	url := "https://api.catonetworks.com/api/v1/graphql2"

	if token == "" {
		fmt.Println("no token provided")
		os.Exit(1)
	}

	if accountId == "" {
		fmt.Println("no account id provided")
		os.Exit(1)
	}

	catoClient, _ := cato.New(url, token, *http.DefaultClient)

	ctx := context.Background()

	entityIds := []string{siteId}

	queryResult, err := catoClient.EntityLookup(ctx, accountId, cato_models.EntityType("site"), nil, nil, nil, nil, entityIds, nil, nil, nil)
	if err != nil {
		fmt.Println("error in EntityLookup: ", err)
		os.Exit(1)
	}

	for _, val := range queryResult.EntityLookup.GetItems() {
		fmt.Println("entity item: ", *val.Entity.Name)
	}

	queryResultJson, err := json.Marshal(queryResult)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(queryResultJson))
}
