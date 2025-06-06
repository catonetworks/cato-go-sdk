package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	cato "github.com/catonetworks/cato-go-sdk"
	cato_models "github.com/catonetworks/cato-go-sdk/models"
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

	catoClient, _ := cato.New(url, token, accountId, nil, nil)

	ctx := context.Background()

	marker := ""

	eventsMarkerData, err := catoClient.EventsFeedIndex(ctx, []string{accountId}, &marker)
	if err != nil {
		fmt.Println("error in auditfeed: ", err)
		return
	}

	marker = *eventsMarkerData.EventsFeed.Marker

	filterContents := `[{"fieldName":"event_type","operator":"is","values":["Security"]}]`

	fmt.Println("Len filterContents: ", len(filterContents))

	filter := make([]*cato_models.EventFeedFieldFilterInput, 0)

	errJ := json.Unmarshal([]byte(filterContents), &filter)
	if errJ != nil {
		fmt.Println(errJ)
		return
	}
	fmt.Println("filter_filled: ", filter)
	filterJson, _ := json.Marshal(filter)
	fmt.Println("Unmarshal: ", string(filterJson))
	//filter0 := &cato_models.EventFeedFieldFilterInput{}

	fmt.Println("Second Attempt")
	f2 := make([]*cato_models.EventFeedFieldFilterInput, 0)
	f2 = append(f2, &cato_models.EventFeedFieldFilterInput{
		FieldName: "event_type",
		Operator:  "is",
		Values:    []string{"Security"},
	})
	f2Json, _ := json.Marshal(f2)
	fmt.Println("Unmarshal: ", string(f2Json))

	return
	// filter = append(filter, filter0)

	queryResultJson, err := json.Marshal(eventsMarkerData)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Marker Run: ", string(queryResultJson))

	for {

		// "last.P2M"
		fmt.Println("CURRENT_MARKER: ", marker)
		eventsData, err := catoClient.EventsFeed(ctx, nil, []string{accountId}, filter, &marker)
		if err != nil {
			fmt.Println("error in auditfeed: ", err)
			return
		}

		queryResultRespJson, err := json.Marshal(eventsData)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Log Run: ", string(queryResultRespJson))

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
