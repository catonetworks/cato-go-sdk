package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	cato "github.com/catonetworks/cato-go-sdk"
)

func main() {
	token := os.Getenv("CATO_API_KEY")
	accountId := os.Getenv("CATO_ACCOUNT_ID")
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
		fmt.Println("no api url provided")
		os.Exit(1)
	}

	catoClient, _ := cato.New(url, token, accountId, nil, nil)

	ctx := context.Background()

	marker := ""
	fetchedCount := int64(0)
	now := time.Now()

	for {

		// "last.P2M"
		fmt.Println("CURRENT_MARKER: ", marker)
		eventsData, err := catoClient.EventsFeedIndex(ctx, []string{accountId}, &marker)
		if err != nil {
			fmt.Println("error in EventsFeedIndex: ", err)
			return
		}

		fetchedCount += eventsData.EventsFeed.FetchedCount

		// queryResultJson, err := json.Marshal(eventsData)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		// fmt.Println(string(queryResultJson))

		fmt.Println("time elapse:", time.Since(now))
		fmt.Println("fetchedCount: ", fetchedCount)

		if !strings.EqualFold(marker, *eventsData.EventsFeed.Marker) {

			// fmt.Println("updating marker to (*auditData.AuditFeed.Marker): ", *auditData.AuditFeed.Marker)
			marker = *eventsData.EventsFeed.Marker

		} else {
			// fmt.Println("end of auditfeed loop...exiting...")
			return
		}
	}

}
