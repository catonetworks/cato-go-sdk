package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

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

	marker := ""

	for {

		// "last.P2M"
		fmt.Println("CURRENT_MARKER: ", marker)
		auditData, err := catoClient.AuditFeed(ctx, nil, []string{accountId}, nil, "last.P2M", nil, &marker)
		if err != nil {
			fmt.Println("error in auditfeed: ", err)
			return
		}

		if auditData.GetAuditFeed() == nil {
			fmt.Println("GetAuditFeed is empty...maybe an error....")
			return
		}

		if len(auditData.AuditFeed.Accounts[0].Records) == 0 {
			fmt.Println("nothing to process...skipping....")
			return
		}

		queryResultJson, err := json.Marshal(auditData)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(queryResultJson))

		if !strings.EqualFold(marker, *auditData.AuditFeed.Marker) {

			// fmt.Println("updating marker to (*auditData.AuditFeed.Marker): ", *auditData.AuditFeed.Marker)
			marker = *auditData.AuditFeed.Marker

		} else {
			// fmt.Println("end of auditfeed loop...exiting...")
			return
		}
	}

}
