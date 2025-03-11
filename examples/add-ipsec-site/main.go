package main

import (
	"context"
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

	siteDescription := "TestSite001 in Seattle"
	stateCode := "US-WA"
	address := "555 This Place"
	city := "Seattle"

	inputRule := cato_models.AddIpsecIkeV2SiteInput{
		Name:               "TestSite001",
		SiteType:           "BRANCH",
		Description:        &siteDescription,
		NativeNetworkRange: "10.99.0.0/16",
		SiteLocation: &cato_models.AddSiteLocationInput{
			CountryCode: "US",
			StateCode:   &stateCode,
			Timezone:    "America/Los_Angeles",
			Address:     &address,
			City:        &city,
		},
	}

	policyChange, err := catoClient.SiteAddIpsecIkeV2Site(ctx, inputRule, accountId)

	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	fmt.Println("SiteAddIpsecIkeV2Site: ", policyChange)

}
