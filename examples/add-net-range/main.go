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
	url := os.Getenv("CATO_API_URL")

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

	inputRange := cato_models.AddNetworkRangeInput{}

	policyChange, err := catoClient.SiteAddNetworkRange(ctx, "", inputRange, accountId)

	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	policyChangeJson, _ := json.Marshal(policyChange)
	fmt.Println(string(policyChangeJson))

	fmt.Println("policyChange: ", policyChange)

}
