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

	if token == "" {
		fmt.Println("no token provided")
		os.Exit(1)
	}

	if accountId == "" {
		fmt.Println("no account id provided")
		os.Exit(1)
	}

	if url == "" {
		fmt.Println("no api url provided")
		os.Exit(1)
	}

	catoClient, _ := cato.New(url, token, accountId, nil, nil)

	ctx := context.Background()

	licData, err := catoClient.Licensing(ctx, accountId)
	if err != nil {
		fmt.Println("error in auditfeed: ", err)
		return
	}

	queryResultJson, err := json.Marshal(licData)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(queryResultJson))

}
