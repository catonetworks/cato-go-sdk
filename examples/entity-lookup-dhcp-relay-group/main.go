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
	token := os.Getenv("TF_VAR_cato_token_dev3")
	accountId := os.Getenv("TF_VAR_account_id_dev3")
	url := "https://api.catonetworks.com/api/v1/graphql2"

	if token == "" {
		fmt.Println("no token provided")
		os.Exit(1)
	}

	if accountId == "" {
		fmt.Println("no account id provided")
		os.Exit(1)
	}

	//  EntityLookup(ctx context.Context, accountID string, typeArg cato_models.EntityType, limit *int64, from *int64, parent *cato_models.EntityInput, search *string, entityIDs []string, sort []*cato_models.SortInput, filters []*cato_models.LookupFilterInput, helperFields []string, interceptors ...clientv2.RequestInterceptor) (*EntityLookup, error)
	catoClient, _ := cato.New(url, token, accountId, nil, nil)

	ctx := context.Background()

	// Query for DHCP relay groups using EntityLookupMinimal
	queryResult, err := catoClient.EntityLookupMinimal(ctx, accountId, cato_models.EntityTypeDhcpRelayGroup, nil, nil, nil, nil, nil)
	if err != nil {
		fmt.Println("EntityLookup query error: ", err)
		os.Exit(1)
	}

	queryResultJson, err := json.MarshalIndent(queryResult, "", "  ")
	if err != nil {
		fmt.Println("JSON marshaling error: ", err)
		os.Exit(1)
	}

	fmt.Println("DHCP Relay Group Entity Lookup Result:")
	fmt.Println(string(queryResultJson))

	// Print summary information
	if queryResult != nil {
		total := "unknown"
		if queryResult.EntityLookup.Total != nil {
			total = fmt.Sprintf("%d", *queryResult.EntityLookup.Total)
		}
		fmt.Printf("\nFound %s DHCP Relay Group entities\n", total)
		for i, item := range queryResult.EntityLookup.Items {
			if item != nil {
				entity := item.GetEntity()
				if entity != nil {
					name := "N/A"
					if entity.Name != nil {
						name = *entity.Name
					}
					fmt.Printf("  %d. ID: %s, Name: %s, Type: %s\n", i+1, entity.ID, name, entity.Type)
					if item.Description != "" {
						fmt.Printf("     Description: %s\n", item.Description)
					}
				}
			}
		}
	}
}
