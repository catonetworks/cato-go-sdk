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
	token_msp := os.Getenv("CATO_API_KEY_MSP")
	accountId_msp := os.Getenv("CATO_ACCOUNT_ID_MSP")
	url := os.Getenv("CATO_API_URL")

	if token_msp == "" {
		fmt.Println("no token provided")
		os.Exit(1)
	}

	if accountId_msp == "" {
		fmt.Println("no account id provided")
		os.Exit(1)
	}

	if url == "" {
		fmt.Println("no url provided")
		os.Exit(1)
	}

	catoClient_msp, _ := cato.New(url, token_msp, accountId_msp, nil, nil)

	ctx := context.Background()

	/////////////////
	// Add account //
	/////////////////
	description := "Example customer account created via SDK"
	addAccountInput := cato_models.AddAccountInput{
		Name:        "Example_Customer_Account2",
		Description: &description,
		Type:        cato_models.AccountProfileTypeCustomer,
		Tenancy:     cato_models.AccountTenancySingleTenant,
		Timezone:    "UTC",
	}

	addAccountResult, err := catoClient_msp.AccountManagementAddAccount(ctx, addAccountInput, accountId_msp)
	if err != nil {
		fmt.Println("error adding account: ", err)
		os.Exit(1)
	}

	// Print the result
	resultJson, _ := json.MarshalIndent(addAccountResult, "", "  ")
	fmt.Println("Account created successfully:")
	fmt.Println(string(resultJson))

	////////////////////
	// Update account //
	////////////////////
	if addAccountResult.AccountManagement.AddAccount != nil {
		account := addAccountResult.AccountManagement.AddAccount
		fmt.Printf("\nAccount Details:\n")
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

		// Update the account description
		newDescription := "Updated account description via SDK"
		updateAccountInput := cato_models.UpdateAccountInput{
			Description: &newDescription,
		}

		catoClient, _ := cato.New(url, token_msp, account.ID, nil, nil)

		result, err := catoClient.AccountManagementUpdateAccount(ctx, updateAccountInput, account.ID)
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

		//////////////////
		// Read account //
		//////////////////
		// Read the updated account using AccountManagement query
		readAccountResponse, err := catoClient.AccountManagement(ctx, account.ID)
		if err != nil {
			fmt.Println("error reading account: ", err)
			os.Exit(1)
		}

		// Print the account details
		readAccountJson, _ := json.MarshalIndent(readAccountResponse, "", "  ")
		fmt.Println("\nAccount read successfully:")
		fmt.Println(string(readAccountJson))

		// Access specific fields from read operation
		if readAccountResponse.AccountManagement != nil && readAccountResponse.AccountManagement.Account != nil {
			readAccountData := readAccountResponse.AccountManagement.Account
			fmt.Printf("\nRead Account Details:\n")
			fmt.Printf("ID: %s\n", readAccountData.ID)
			fmt.Printf("Name: %s\n", readAccountData.Name)
			fmt.Printf("Type: %s\n", readAccountData.Type)
			fmt.Printf("Tenancy: %s\n", readAccountData.Tenancy)
			fmt.Printf("TimeZone: %s\n", readAccountData.TimeZone)
			if readAccountData.Description != nil {
				fmt.Printf("Description: %s\n", *readAccountData.Description)
			}
			fmt.Printf("Created by: %s\n", readAccountData.Audit.CreatedBy)
			fmt.Printf("Created time: %s\n", readAccountData.Audit.CreatedTime)
		}
		////////////////////
		// Remove account //
		////////////////////
		removedAccountResult, err := catoClient.AccountManagementRemoveAccount(ctx, account.ID, accountId_msp)
		if err != nil {
			fmt.Println("error removing account: ", err)
			os.Exit(1)
		}

		// Print the result
		resultJson, _ = json.MarshalIndent(removedAccountResult, "", "  ")
		fmt.Println("\nAccount removed successfully:")
		fmt.Println(string(resultJson))

		// Access specific fields
		if removedAccountResult.AccountManagement.RemoveAccount != nil {
			account := removedAccountResult.AccountManagement.RemoveAccount.AccountInfo
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
}
