package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

	//  EntityLookup(ctx context.Context, accountID string, typeArg cato_models.EntityType, limit *int64, from *int64, parent *cato_models.EntityInput, search *string, entityIDs []string, sort []*cato_models.SortInput, filters []*cato_models.LookupFilterInput, helperFields []string, interceptors ...clientv2.RequestInterceptor) (*EntityLookup, error)
	catoClient, _ := cato.New(url, token, *http.DefaultClient)

	ctx := context.Background()

	queryResult, err := catoClient.EntityLookup(ctx, accountId, cato_models.EntityType("site"), nil, nil, nil, nil, nil, nil, nil, nil)
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
}
