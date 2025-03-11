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
	url := "https://api.catonetworks.com/api/v1/graphql2"

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

	inputRule := cato_models.PolicyAddSectionInput{
		At: &cato_models.PolicySectionPositionInput{
			Position: "LAST_IN_POLICY",
		},
		Section: &cato_models.PolicyAddSectionInfoInput{
			Name: "IFW Demo Section",
		},
	}

	// PolicyInternetFirewallAddSection(ctx, &cato_models.InternetFirewallPolicyMutationInput{}, input, r.client.AccountId)
	policyChange, err := catoClient.PolicyInternetFirewallAddSection(ctx, &cato_models.InternetFirewallPolicyMutationInput{}, inputRule, accountId)

	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	policyChangeJson, _ := json.Marshal(policyChange)
	fmt.Println(string(policyChangeJson))

	fmt.Println("policyChange: ", policyChange)

}
