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
	accountId := "16359"
	// accountId := os.Getenv("CATO_ACCOUNT_ID")
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
		fmt.Println("no url provided")
		os.Exit(1)
	}

	catoClient, _ := cato.New(url, token, accountId, nil, nil)

	ctx := context.Background()

	// Update the account description
	newDescription := "Updated account description via SDK"
	updateAccountInput := cato_models.UpdateAccountInput{
		Description: &newDescription,
	}

	result, err := catoClient.AccountManagementUpdateAccount(ctx, updateAccountInput, accountId)
	if err != nil {
		fmt.Println("error updating account: ", err)
		os.Exit(1)
	}

	// Print the result
	resultJson, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println("Account updated successfully:")
	fmt.Println(string(resultJson))

	// Access specific fields
	if result.AccountManagement.UpdateAccount != nil {
		account := result.AccountManagement.UpdateAccount
		fmt.Printf("\nUpdated Account Details:\n")
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
