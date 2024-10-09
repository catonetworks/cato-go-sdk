package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	cato "github.com/cato-networks/cato-go-sdk"
	cato_models "github.com/cato-networks/cato-go-sdk/models"
	"github.com/cato-networks/cato-go-sdk/scalars"
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

	catoClient, _ := cato.New(url, token, *http.DefaultClient)

	ctx := context.Background()

	position := cato_models.PolicyRulePositionEnum("LAST_IN_POLICY")
	hostRefInput := []*cato_models.HostRefInput{}
	siteRefInput := []*cato_models.SiteRefInput{}
	iprange := []*cato_models.IPAddressRangeInput{}
	globalIpRange := []*cato_models.GlobalIPRangeRefInput{}
	networkInterfaceRefInput := []*cato_models.NetworkInterfaceRefInput{}
	siteNetworkSubnetRefInput := []*cato_models.SiteNetworkSubnetRefInput{}
	floatingSubnetRefInput := []*cato_models.FloatingSubnetRefInput{}
	userRefInput := []*cato_models.UserRefInput{}
	usersGroupRefInput := []*cato_models.UsersGroupRefInput{}
	groupRefInput := []*cato_models.GroupRefInput{}
	systemGroupRefInput := []*cato_models.SystemGroupRefInput{}

	connectionOrigin := cato_models.ConnectionOriginEnum("ANY")
	actionEnum := cato_models.InternetFirewallActionEnum("ALLOW")

	domainList := []string{"slashdot.org"}
	fqdnList := []string{"www.slashdot.org"}

	remoteAsnList := []scalars.Asn16{}

	inputRule := cato_models.InternetFirewallAddRuleInput{
		At: &cato_models.PolicyRulePositionInput{
			Position: &position,
		},
		Rule: &cato_models.InternetFirewallAddRuleDataInput{
			Enabled:     false,
			Name:        "TestScalarRule01",
			Description: "TestScalarRule01",
			Source: &cato_models.InternetFirewallSourceInput{
				IP:                []string{},
				Host:              hostRefInput,
				Site:              siteRefInput,
				Subnet:            []string{},
				IPRange:           iprange,
				GlobalIPRange:     globalIpRange,
				NetworkInterface:  networkInterfaceRefInput,
				SiteNetworkSubnet: siteNetworkSubnetRefInput,
				FloatingSubnet:    floatingSubnetRefInput,
				User:              userRefInput,
				UsersGroup:        usersGroupRefInput,
				Group:             groupRefInput,
				SystemGroup:       systemGroupRefInput,
			},
			ConnectionOrigin: connectionOrigin,
			Country:          []*cato_models.CountryRefInput{},
			Device:           []*cato_models.DeviceProfileRefInput{},
			DeviceOs:         []cato_models.OperatingSystem{},
			Destination: &cato_models.InternetFirewallDestinationInput{
				Application:            []*cato_models.ApplicationRefInput{},
				CustomApp:              []*cato_models.CustomApplicationRefInput{},
				AppCategory:            []*cato_models.ApplicationCategoryRefInput{},
				CustomCategory:         []*cato_models.CustomCategoryRefInput{},
				SanctionedAppsCategory: []*cato_models.SanctionedAppsCategoryRefInput{},
				Country:                []*cato_models.CountryRefInput{},
				Domain:                 domainList,
				Fqdn:                   fqdnList,
				IP:                     []string{},
				Subnet:                 []string{},
				IPRange:                iprange,
				GlobalIPRange:          globalIpRange,
				RemoteAsn:              remoteAsnList,
			},
			Service: &cato_models.InternetFirewallServiceTypeInput{},
			Action:  actionEnum,
			Schedule: &cato_models.PolicyScheduleInput{
				ActiveOn: "CUSTOM_RECURRING",
				CustomRecurring: &cato_models.PolicyCustomRecurringInput{
					// From: cato_scalars.Time("08:00"),
					From: "10:00:00",
					To:   "11:00:00",
					Days: []cato_models.DayOfWeek{
						"MONDAY",
					},
				},
			},
			Tracking: &cato_models.PolicyTrackingInput{
				Event: &cato_models.PolicyRuleTrackingEventInput{
					Enabled: true,
				},
				Alert: &cato_models.PolicyRuleTrackingAlertInput{
					Enabled:           false,
					Frequency:         "DAILY",
					SubscriptionGroup: []*cato_models.SubscriptionGroupRefInput{},
					MailingList:       []*cato_models.SubscriptionMailingListRefInput{},
					Webhook:           []*cato_models.SubscriptionWebhookRefInput{},
				},
			},
		},
	}

	policyChange, err := catoClient.PolicyInternetFirewallAddRule(ctx, inputRule, accountId)

	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	policyChangeJson, _ := json.Marshal(policyChange)
	fmt.Println(string(policyChangeJson))

	fmt.Println("policyChange: ", policyChange)

}
