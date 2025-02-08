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

	catoClient, _ := cato.New(url, token, nil)

	ctx := context.Background()

	internetFirewallPolicyMutationInput := &cato_models.InternetFirewallPolicyMutationInput{}
	internetFirewallRemoveRuleInput := cato_models.InternetFirewallRemoveRuleInput{
		ID: "UUID_OF_RULE",
	}

	queryResult, err := catoClient.PolicyInternetFirewallRemoveRule(ctx, internetFirewallPolicyMutationInput, internetFirewallRemoveRuleInput, accountId)
	if err != nil {
		fmt.Println("policy query error: ", err)
		return
	}

	queryResultJson, err := json.Marshal(queryResult)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(queryResultJson))

	publishDataIfEnabled := &cato_models.PolicyPublishRevisionInput{}
	_, err = catoClient.PolicyInternetFirewallPublishPolicyRevision(ctx, &cato_models.InternetFirewallPolicyMutationInput{}, publishDataIfEnabled, accountId)
	if err != nil {
		fmt.Println("policy publish query error: ", err)
		return
	}
}
