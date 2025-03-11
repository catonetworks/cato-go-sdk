# Cato Networks Go Client SDK

This client is built to support the Cato Networks GraphQL API as outlined here:

https://www.catonetworks.com/platform/cato-api/

### Installation

```bash
go install github.com/catonetworks/cato-go-sdk
```

### Client Initialization

UPDATE: as of version 0.2.0, the client instantiation call has changed along with the format of additional HTTP headers passed to the client.

URL: The URL to the Cato Networks GraphQL endpoint (https://api.catonetworks.com/api/v1/graphql2)
Token: Your API access token
HTTP Client: Use either the default HTTP client or leverage more advanced configuration options by passing in a client.

```go
catoClient, _ := cato.New(url, token, accountId, nil, nil)
```

### Client Queries
In this example, we are looking up Entity information regarding Cato Site data from the API. We pass in the EntityID as a string slice ([]string{}) which can leverage multiple EntityIDs.

To read the returned data, we can either leverage a JSON Marshall call to read the JSON directly or perform GetItems()/GetItem() on the elements.

```go
	queryResult, err := catoClient.EntityLookup(ctx, accountId, cato_models.EntityType("site"), nil, nil, nil, nil, entityIds, nil, nil, nil)
	if err != nil {
		fmt.Println("error in EntityLookup: ", err)
		os.Exit(1)
	}

	for _, val := range queryResult.EntityLookup.GetItems() {
		fmt.Println("entity item: ", *val.Entity.Name)
	}

	queryResultJson, err := json.Marshal(queryResult)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(queryResultJson))

```

### Policy Queries
#### Note:  The Policy query has been updated as of 9/17 to support the new Cato Networks WAN query parameters in addition to the previous Internet Firewall parameters. Those who have existing code leveraging a single input parameter to the query will need to ensure that you are using version < 0.2.20

```go

accountId := "12345"
queryIfwPolicy := &cato_models.InternetFirewallPolicyInput{}
queryWanPolicy := &cato_models.WanFirewallPolicyInput{}
queryResult, err := catoClient.Policy(ctx, queryIfwPolicy, queryWanPolicy, accountId)

```

### Client Mutations
Mustations are used in GraphQL to perform a change. This can include a create/update/delete operation. In the example below, we are creating a mostly blank internet firewall rule which is named, "TestRule101". This is set to be inserted as the last rule in the policy to ALLOW any traffic to slashdot.org.

```go
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
	actionEnum := cato_models.InternetFirewallActionEnum("ALLOW")

	domainList := []string{"slashdot.org"}
	fqdnList := []string{}

	inputRule := cato_models.InternetFirewallAddRuleInput{
		At: &cato_models.PolicyRulePositionInput{
			Position: &position,
		},
		Rule: &cato_models.InternetFirewallAddRuleDataInput{
			Enabled:     false,
			Name:        "TestRule101",
			Description: "TestRule101",
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
				RemoteAsn:              []string{},
			},
			Service: &cato_models.InternetFirewallServiceTypeInput{},
			Action:  actionEnum,
			Schedule: &cato_models.PolicyScheduleInput{
				ActiveOn: "ALWAYS",
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
```

### Examples

Additional examples can be found in the [examples](examples/) folder.