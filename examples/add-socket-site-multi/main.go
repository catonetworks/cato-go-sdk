package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"

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

	if url == "" {
		fmt.Println("no url provided")
		os.Exit(1)
	}

	numberOfSites := 5

	var wg sync.WaitGroup

	inputList := make([]cato_models.AddSocketSiteInput, numberOfSites)

	siteDescription := "TestSiteExample in Seattle"
	stateCode := "US-WA"
	address := "555 This Place"
	city := "Seattle"

	for x := 0; x < numberOfSites; x++ {
		inputList[x] = cato_models.AddSocketSiteInput{
			Name:               "TestSiteExample" + strconv.Itoa(x),
			ConnectionType:     "SOCKET_X1500",
			SiteType:           "BRANCH",
			Description:        &siteDescription,
			NativeNetworkRange: "10.99." + strconv.Itoa(x) + ".0/24",
			SiteLocation: &cato_models.AddSiteLocationInput{
				CountryCode: "US",
				StateCode:   &stateCode,
				Timezone:    "America/Los_Angeles",
				Address:     &address,
				City:        &city,
			},
		}
	}

	inputList = append(inputList, cato_models.AddSocketSiteInput{
		Name:               "TestSiteAWSExample",
		ConnectionType:     "SOCKET_AWS1500",
		SiteType:           "CLOUD_DC",
		Description:        &siteDescription,
		NativeNetworkRange: "10.98.0.0/24",
		SiteLocation: &cato_models.AddSiteLocationInput{
			CountryCode: "US",
			StateCode:   &stateCode,
			Timezone:    "America/Los_Angeles",
			Address:     &address,
			City:        &city,
		},
	})

	wg.Add(numberOfSites + 1)

	catoClient, _ := cato.New(url, token, accountId, nil, nil)
	ctx := context.Background()

	for _, v := range inputList {
		go doCatoCall(ctx, &wg, catoClient, v, accountId)
	}

	wg.Wait()

}

func doCatoCall(ctx context.Context, wg *sync.WaitGroup, c *cato.Client, v cato_models.AddSocketSiteInput, accountId string) {
	defer wg.Done()
	policyChange, err := c.SiteAddSocketSite(ctx, v, accountId)
	if err != nil {
		fmt.Println("error: ", err)
	}
	queryResultJson, err := json.Marshal(policyChange)
	if err != nil {
		fmt.Println("SiteAddSocketSite error: ", err)
	}

	fmt.Println(string(queryResultJson))
}
