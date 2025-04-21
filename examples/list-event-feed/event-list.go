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

	catoClient, _ := cato.New(url, token, accountId, nil, nil)

	ctx := context.Background()

	marker := ""

	eventsMarkerData, err := catoClient.EventsFeedIndex(ctx, []string{accountId}, &marker)
	if err != nil {
		fmt.Println("error in auditfeed: ", err)
		return
	}

	marker = *eventsMarkerData.EventsFeed.Marker

	queryResultJson, err := json.Marshal(eventsMarkerData)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(queryResultJson))

	for {

		// "last.P2M"
		fmt.Println("CURRENT_MARKER: ", marker)
		eventsData, err := catoClient.EventsFeed(ctx, nil, []string{accountId}, nil, &marker)
		if err != nil {
			fmt.Println("error in auditfeed: ", err)
			return
		}

		queryResultRespJson, err := json.Marshal(eventsData)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(queryResultRespJson))

		if eventsData.GetEventsFeed() == nil {
			fmt.Println("GetAuditFeed is empty...maybe an error....")
			return
		}

		if len(eventsData.EventsFeed.Accounts[0].Records) == 0 {
			fmt.Println("nothing to process...skipping....")
			return
		}

		queryResultJson, err := json.Marshal(eventsData)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(queryResultJson))

		if !strings.EqualFold(marker, *eventsData.EventsFeed.Marker) {

			// fmt.Println("updating marker to (*auditData.AuditFeed.Marker): ", *auditData.AuditFeed.Marker)
			marker = *eventsData.EventsFeed.Marker

		} else {
			// fmt.Println("end of auditfeed loop...exiting...")
			return
		}
	}

}
