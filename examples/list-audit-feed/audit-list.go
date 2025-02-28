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

	// "last.P2M"
	auditData, err := catoClient.AuditFeed(ctx, nil, []string{accountId}, nil, "last.P2M", nil, nil)
	if err != nil {
		fmt.Println("error in eventsfeed: ", err)
		return
	}

	queryResultJson, err := json.Marshal(auditData)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(queryResultJson))

}
