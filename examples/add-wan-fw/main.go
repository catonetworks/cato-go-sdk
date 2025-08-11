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

	catoClient, _ := cato.New(url, token, accountId, nil, nil)

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

	inputRule := cato_models.WanFirewallAddRuleInput{
		At: &cato_models.PolicyRulePositionInput{
			Position: &position,
		},
		Rule: &cato_models.WanFirewallAddRuleDataInput{
			Enabled:      false,
			ActivePeriod: &cato_models.PolicyRuleActivePeriodInput{},
			Name:         "TestScalarRule07",
			Description:  "TestScalarRule07",
			Source: &cato_models.WanFirewallSourceInput{
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
			Destination: &cato_models.WanFirewallDestinationInput{
				IP:            []string{},
				Subnet:        []string{},
				IPRange:       iprange,
				GlobalIPRange: globalIpRange,
			},
			Service: &cato_models.WanFirewallServiceTypeInput{
				Standard: []*cato_models.ServiceRefInput{
					&cato_models.ServiceRefInput{
						By:    "NAME",
						Input: "Agora",
					},
				},
			},
			// Action: actionEnum,
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

	policyChange, err := catoClient.PolicyWanFirewallAddRule(ctx, inputRule, accountId)

	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	publishDataIfEnabled := &cato_models.PolicyPublishRevisionInput{}
	_, err = catoClient.PolicyInternetFirewallPublishPolicyRevision(ctx, &cato_models.InternetFirewallPolicyMutationInput{}, publishDataIfEnabled, accountId)
	if err != nil {
		fmt.Println("policy publish query error: ", err)
		return
	}

	policyChangeJson, _ := json.Marshal(policyChange)
	fmt.Println(string(policyChangeJson))

	fmt.Println("policyChange: ", policyChange)

}
