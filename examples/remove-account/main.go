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

	// The account ID to remove (this should be provided as input)
	accountToRemove := "16359"

	if token == "" {
		fmt.Println("no token provided")
		os.Exit(1)
	}

	if accountId == "" {
		fmt.Println("no account id provided")
		os.Exit(1)
	}

	if url == "" {
		fmt.Println("no url provided")
		os.Exit(1)
	}

	if accountToRemove == "" {
		fmt.Println("no account to remove provided (set CATO_ACCOUNT_TO_REMOVE)")
		os.Exit(1)
	}

	catoClient, _ := cato.New(url, token, accountId, nil, nil)

	ctx := context.Background()

	// Remove the account
	result, err := catoClient.AccountManagementRemoveAccount(ctx, accountToRemove, accountId)
	if err != nil {
		fmt.Println("error removing account: ", err)
		os.Exit(1)
	}

	// Print the result
	resultJson, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println("Account removed successfully:")
	fmt.Println(string(resultJson))

	// Access specific fields
	if result.AccountManagement.RemoveAccount != nil {
		account := result.AccountManagement.RemoveAccount.AccountInfo
		fmt.Printf("\nRemoved Account Details:\n")
		fmt.Printf("ID: %s\n", account.ID)
		fmt.Printf("Name: %s\n", account.Name)
		fmt.Printf("Type: %s\n", account.Type)
		fmt.Printf("Tenancy: %s\n", account.Tenancy)
		fmt.Printf("TimeZone: %s\n", account.TimeZone)
		if account.Description != nil {
			fmt.Printf("Description: %s\n", *account.Description)
		}
		fmt.Printf("Created by: %s\n", account.Audit.CreatedBy)
		fmt.Printf("Created time: %s\n", account.Audit.CreatedTime)
	}
}
