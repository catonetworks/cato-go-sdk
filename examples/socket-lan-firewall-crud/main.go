package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	cato "github.com/catonetworks/cato-go-sdk"
	cato_models "github.com/catonetworks/cato-go-sdk/models"
	"github.com/catonetworks/cato-go-sdk/scalars"
)

// Debug flag - set to true to print all requests and responses
var debug = os.Getenv("DEBUG") != ""

func debugPrint(label string, v any) {
	if !debug {
		return
	}
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("[DEBUG] %s: (failed to marshal: %v)\n", label, err)
		return
	}
	fmt.Printf("[DEBUG] %s:\n%s\n\n", label, string(jsonData))
}

// readAndDisplayPolicy queries the Socket LAN policy and displays the specified rule
func readAndDisplayPolicy(catoClient cato.CatoClient, ctx context.Context, accountId string, parentRuleId string) {
	policyResult, err := catoClient.PolicySocketLanPolicy(ctx, accountId, nil)
	if err != nil {
		fmt.Printf("  Warning: Could not read policy: %v\n", err)
		return
	}

	debugPrint("PolicySocketLanPolicy Response", policyResult)

	// Find the parent rule by ID
	for _, ruleWrapper := range policyResult.Policy.SocketLan.Policy.Rules {
		if ruleWrapper.Rule.ID == parentRuleId {
			// Display parent rule details
			fmt.Printf("  Parent Rule:\n")
			fmt.Printf("    ID: %s\n", ruleWrapper.Rule.ID)
			fmt.Printf("    Name: %s\n", ruleWrapper.Rule.Name)
			fmt.Printf("    Description: %s\n", ruleWrapper.Rule.Description)
			fmt.Printf("    Direction: %s\n", ruleWrapper.Rule.Direction)
			fmt.Printf("    Transport: %s\n", ruleWrapper.Rule.Transport)
			fmt.Printf("    NAT Enabled: %v\n", ruleWrapper.Rule.Nat.Enabled)

			// Display child firewall rules
			if len(ruleWrapper.Rule.Firewall) > 0 {
				fmt.Printf("  Child Firewall Rules (%d):\n", len(ruleWrapper.Rule.Firewall))
				for _, fwRule := range ruleWrapper.Rule.Firewall {
					fmt.Printf("    - ID: %s\n", fwRule.Rule.ID)
					fmt.Printf("      Name: %s\n", fwRule.Rule.Name)
					fmt.Printf("      Action: %s\n", fwRule.Rule.Action)
					fmt.Printf("      Direction: %s\n", fwRule.Rule.Direction)
				}
			} else {
				fmt.Printf("  Child Firewall Rules: (none)\n")
			}
			return
		}
	}
	fmt.Printf("  Warning: Parent rule %s not found in policy\n", parentRuleId)
}

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

	if debug {
		fmt.Println("[DEBUG] Debug mode enabled - printing all requests and responses")
		fmt.Printf("[DEBUG] API URL: %s\n", url)
		fmt.Printf("[DEBUG] Account ID: %s\n\n", accountId)
	}

	catoClient, _ := cato.New(url, token, accountId, nil, nil)

	ctx := context.Background()

	/////////////////////////////////////////////
	// Step 1: Create a new Section
	/////////////////////////////////////////////

	fmt.Printf("======================================\n")
	fmt.Printf("Step 1: Creating Socket LAN Policy Section\n")
	fmt.Printf("======================================\n")

	sectionPosition := cato_models.PolicySectionPositionEnumLastInPolicy
	addSectionInput := cato_models.PolicyAddSectionInput{
		At: &cato_models.PolicySectionPositionInput{
			Position: sectionPosition,
		},
		Section: &cato_models.PolicyAddSectionInfoInput{
			Name: "SDK Example Section",
		},
	}

	debugPrint("AddSection Request", addSectionInput)

	sectionResult, err := catoClient.PolicySocketLanAddSection(ctx, addSectionInput, accountId)
	if err != nil {
		fmt.Println("error adding section: ", err)
		os.Exit(1)
	}

	debugPrint("AddSection Response", sectionResult)

	sectionId := sectionResult.Policy.SocketLan.AddSection.Section.Section.ID
	fmt.Printf("Section created successfully:\n")
	fmt.Printf("  ID: %s\n", sectionId)
	fmt.Printf("  Name: %s\n", sectionResult.Policy.SocketLan.AddSection.Section.Section.Name)

	/////////////////////////////////////////////
	// Step 2: Create a MINIMAL parent Socket LAN rule in the new section
	/////////////////////////////////////////////

	fmt.Printf("\n======================================\n")
	fmt.Printf("Step 2: Creating Minimal Parent Socket LAN Rule\n")
	fmt.Printf("======================================\n")

	position := cato_models.PolicyRulePositionEnumLastInSection
	sectionRef := sectionId

	// Minimal parent rule - only required fields
	socketLanAddRuleInput := cato_models.SocketLanAddRuleInput{
		At: &cato_models.PolicyRulePositionInput{
			Position: &position,
			Ref:      &sectionRef,
		},
		Rule: &cato_models.SocketLanAddRuleDataInput{
			Name:        "Socket LAN Rule",
			Description: "Minimal rule - will be updated with kitchen sink params",
			Enabled:     true,
			Direction:   cato_models.SocketLanDirectionTo,
			Transport:   cato_models.SocketLanTransportTypeLan,
			Source: &cato_models.SocketLanSourceInput{
				Subnet:            []string{"192.168.1.0/24"},
				Vlan:              []scalars.Vlan{},
				IPRange:           []*cato_models.IPAddressRangeInput{},
				Group:             []*cato_models.GroupRefInput{},
				NetworkInterface:  []*cato_models.NetworkInterfaceRefInput{},
				SystemGroup:       []*cato_models.SystemGroupRefInput{},
				Host:              []*cato_models.HostRefInput{},
				IP:                []string{},
				GlobalIPRange:     []*cato_models.GlobalIPRangeRefInput{},
				FloatingSubnet:    []*cato_models.FloatingSubnetRefInput{},
				SiteNetworkSubnet: []*cato_models.SiteNetworkSubnetRefInput{},
			},
			Destination: &cato_models.SocketLanDestinationInput{
				Subnet:            []string{"10.0.0.0/24"},
				Vlan:              []scalars.Vlan{},
				IPRange:           []*cato_models.IPAddressRangeInput{},
				Group:             []*cato_models.GroupRefInput{},
				NetworkInterface:  []*cato_models.NetworkInterfaceRefInput{},
				SystemGroup:       []*cato_models.SystemGroupRefInput{},
				Host:              []*cato_models.HostRefInput{},
				IP:                []string{},
				GlobalIPRange:     []*cato_models.GlobalIPRangeRefInput{},
				FloatingSubnet:    []*cato_models.FloatingSubnetRefInput{},
				SiteNetworkSubnet: []*cato_models.SiteNetworkSubnetRefInput{},
			},
			Site: &cato_models.SocketLanSiteInput{
				Site: []*cato_models.SiteRefInput{
					{By: "ID", Input: "180279"},
				},
				Group: []*cato_models.GroupRefInput{},
			},
			Service: &cato_models.SocketLanServiceInput{
				Simple: []*cato_models.SimpleServiceInput{},
				Custom: []*cato_models.CustomServiceInput{},
			},
			Nat: &cato_models.SocketLanNatSettingsInput{
				Enabled: false,
				NatType: cato_models.SocketLanNatTypeDynamicPat,
			},
		},
	}

	debugPrint("AddRule Request - socketLanAddRuleInput (Minimal Parent)", socketLanAddRuleInput)

	parentResult, err := catoClient.PolicySocketLanAddRule(ctx, socketLanAddRuleInput, accountId)
	if err != nil {
		fmt.Println("error adding parent Socket LAN rule: ", err)
		os.Exit(1)
	}

	debugPrint("AddRule Response (Parent)", parentResult)

	if parentResult.Policy.SocketLan.AddRule.Rule == nil {
		fmt.Println("error: no rule returned from parent rule creation")
		os.Exit(1)
	}

	parentRule := parentResult.Policy.SocketLan.AddRule.Rule
	parentRuleId := parentRule.Rule.ID

	fmt.Printf("Parent Socket LAN rule created successfully:\n")
	fmt.Printf("  ID: %s\n", parentRuleId)
	fmt.Printf("  Name: %s\n", parentRule.Rule.Name)

	/////////////////////////////////////////////
	// Step 3: Create a MINIMAL child Firewall rule
	/////////////////////////////////////////////

	fmt.Printf("\n======================================\n")
	fmt.Printf("Step 3: Creating Minimal Child Firewall Rule\n")
	fmt.Printf("======================================\n")

	// Minimal child firewall rule - only required fields
	firewallAddRuleInput := cato_models.SocketLanFirewallAddRuleInput{
		At: &cato_models.PolicySubRulePositionInput{
			Position: cato_models.PolicySubRulePositionEnumFirstInRule,
			Ref:      parentRuleId,
		},
		Rule: &cato_models.SocketLanFirewallAddRuleDataInput{
			Name:        "Firewall Rule",
			Description: "Minimal rule - will be updated with kitchen sink params",
			Enabled:     true,
			Direction:   cato_models.SocketLanFirewallDirectionTo,
			Action:      cato_models.SocketLanFirewallActionAllow,
			Source: &cato_models.SocketLanFirewallSourceInput{
				Subnet:            []string{"192.168.1.0/24"},
				Vlan:              []scalars.Vlan{},
				Mac:               []string{},
				IPRange:           []*cato_models.IPAddressRangeInput{},
				Group:             []*cato_models.GroupRefInput{},
				Site:              []*cato_models.SiteRefInput{},
				NetworkInterface:  []*cato_models.NetworkInterfaceRefInput{},
				SystemGroup:       []*cato_models.SystemGroupRefInput{},
				Host:              []*cato_models.HostRefInput{},
				IP:                []string{},
				GlobalIPRange:     []*cato_models.GlobalIPRangeRefInput{},
				FloatingSubnet:    []*cato_models.FloatingSubnetRefInput{},
				SiteNetworkSubnet: []*cato_models.SiteNetworkSubnetRefInput{},
			},
			Destination: &cato_models.SocketLanFirewallDestinationInput{
				Subnet:            []string{"10.0.0.0/24"},
				Vlan:              []scalars.Vlan{},
				IPRange:           []*cato_models.IPAddressRangeInput{},
				Group:             []*cato_models.GroupRefInput{},
				Site:              []*cato_models.SiteRefInput{},
				NetworkInterface:  []*cato_models.NetworkInterfaceRefInput{},
				SystemGroup:       []*cato_models.SystemGroupRefInput{},
				Host:              []*cato_models.HostRefInput{},
				IP:                []string{},
				GlobalIPRange:     []*cato_models.GlobalIPRangeRefInput{},
				FloatingSubnet:    []*cato_models.FloatingSubnetRefInput{},
				SiteNetworkSubnet: []*cato_models.SiteNetworkSubnetRefInput{},
			},
			Application: &cato_models.SocketLanFirewallApplicationInput{
				Application:   []*cato_models.ApplicationRefInput{},
				CustomApp:     []*cato_models.CustomApplicationRefInput{},
				Domain:        []string{},
				Fqdn:          []string{},
				IP:            []string{},
				Subnet:        []string{},
				IPRange:       []*cato_models.IPAddressRangeInput{},
				GlobalIPRange: []*cato_models.GlobalIPRangeRefInput{},
			},
			Service: &cato_models.SocketLanFirewallServiceTypeInput{
				Simple:   []*cato_models.SimpleServiceInput{},
				Standard: []*cato_models.ServiceRefInput{},
				Custom:   []*cato_models.CustomServiceInput{},
			},
			Tracking: &cato_models.PolicyTrackingInput{
				Event: &cato_models.PolicyRuleTrackingEventInput{
					Enabled: false,
				},
				Alert: &cato_models.PolicyRuleTrackingAlertInput{
					Enabled:           false,
					Frequency:         cato_models.PolicyRuleTrackingFrequencyEnumDaily,
					MailingList:       []*cato_models.SubscriptionMailingListRefInput{},
					SubscriptionGroup: []*cato_models.SubscriptionGroupRefInput{},
					Webhook:           []*cato_models.SubscriptionWebhookRefInput{},
				},
			},
		},
	}

	debugPrint("AddRule Request - firewallAddRuleInput (Minimal Child)", firewallAddRuleInput)

	childResult, err := catoClient.PolicySocketLanFirewallAddRule(ctx, accountId, nil, firewallAddRuleInput)
	if err != nil {
		fmt.Println("error adding child Firewall rule: ", err)
		os.Exit(1)
	}

	debugPrint("AddRule Response (Child)", childResult)

	if childResult.Policy.SocketLan.Firewall.AddRule.Rule == nil {
		fmt.Println("error: no rule returned from child firewall rule creation")
		os.Exit(1)
	}

	childRule := childResult.Policy.SocketLan.Firewall.AddRule.Rule
	childRuleId := childRule.Rule.ID

	fmt.Printf("Child Firewall rule created successfully:\n")
	fmt.Printf("  ID: %s\n", childRuleId)
	fmt.Printf("  Name: %s\n", childRule.Rule.Name)

	/////////////////////////////////////////////
	// Step 4: Update parent Socket LAN rule with KITCHEN SINK params
	/////////////////////////////////////////////

	fmt.Printf("\n======================================\n")
	fmt.Printf("Step 4: Updating Parent Socket LAN Rule (Kitchen Sink)\n")
	fmt.Printf("======================================\n")

	// Kitchen sink values from JSON example
	updatedParentName := "Socket LAN Network Rule Kitchen Sink"
	updatedParentDescription := "Updated with all possible parameters"
	enabledFlag := true
	direction := cato_models.SocketLanDirectionTo
	transport := cato_models.SocketLanTransportTypeLan

	socketLanUpdateRuleInput := cato_models.SocketLanUpdateRuleInput{
		ID: parentRuleId,
		Rule: &cato_models.SocketLanUpdateRuleDataInput{
			Name:        &updatedParentName,
			Description: &updatedParentDescription,
			Enabled:     &enabledFlag,
			Direction:   &direction,
			Transport:   &transport,
			// Site - both site and group references
			Site: &cato_models.SocketLanSiteUpdateInput{
				Site: []*cato_models.SiteRefInput{
					{By: "ID", Input: "180279"},
				},
				Group: []*cato_models.GroupRefInput{
					{By: "ID", Input: "636881"},
				},
			},
			// Source - all possible fields from kitchen sink
			Source: &cato_models.SocketLanSourceUpdateInput{
				Vlan: []scalars.Vlan{1},
				IPRange: []*cato_models.IPAddressRangeInput{
					{From: "1.2.3.1", To: "1.2.3.3"},
				},
				Group: []*cato_models.GroupRefInput{
					{By: "ID", Input: "636881"},
				},
				Subnet: []string{"1.2.3.0/24"},
				NetworkInterface: []*cato_models.NetworkInterfaceRefInput{
					{By: "ID", Input: "216208"},
				},
				SystemGroup: []*cato_models.SystemGroupRefInput{
					{By: "ID", Input: "7S"},
				},
				Host: []*cato_models.HostRefInput{
					{By: "ID", Input: "2762912"},
				},
				IP: []string{"1.2.3.4"},
				GlobalIPRange: []*cato_models.GlobalIPRangeRefInput{
					{By: "ID", Input: "2528637"},
				},
				FloatingSubnet: []*cato_models.FloatingSubnetRefInput{
					{By: "ID", Input: "2528636"},
				},
				SiteNetworkSubnet: []*cato_models.SiteNetworkSubnetRefInput{
					{By: "ID", Input: "TjI3NjI5MTA="},
				},
			},
			// Destination - all possible fields from kitchen sink
			Destination: &cato_models.SocketLanDestinationUpdateInput{
				Vlan: []scalars.Vlan{1},
				IPRange: []*cato_models.IPAddressRangeInput{
					{From: "1.2.3.4", To: "1.2.3.5"},
				},
				Group: []*cato_models.GroupRefInput{
					{By: "ID", Input: "636881"},
				},
				Subnet: []string{"1.2.3.0/24"},
				NetworkInterface: []*cato_models.NetworkInterfaceRefInput{
					{By: "ID", Input: "216208"},
				},
				SystemGroup: []*cato_models.SystemGroupRefInput{
					{By: "ID", Input: "7S"},
				},
				Host: []*cato_models.HostRefInput{
					{By: "ID", Input: "2762912"},
				},
				IP: []string{"1.2.3.4"},
				GlobalIPRange: []*cato_models.GlobalIPRangeRefInput{
					{By: "ID", Input: "2528637"},
				},
				FloatingSubnet: []*cato_models.FloatingSubnetRefInput{
					{By: "ID", Input: "2528636"},
				},
				SiteNetworkSubnet: []*cato_models.SiteNetworkSubnetRefInput{
					{By: "ID", Input: "TjI3NjI5MTA="},
				},
			},
			// Service - simple and custom services
			Service: &cato_models.SocketLanServiceUpdateInput{
				Simple: []*cato_models.SimpleServiceInput{
					{Name: "HTTP"},
				},
				Custom: []*cato_models.CustomServiceInput{
					{
						Port:     []scalars.Port{"80"},
						Protocol: cato_models.IPProtocolTCP,
					},
					{
						PortRange: &cato_models.PortRangeInput{From: "80", To: "81"},
						Protocol:  cato_models.IPProtocolUDP,
					},
				},
			},
			// NAT settings
			Nat: &cato_models.SocketLanNatSettingsUpdateInput{
				Enabled: &enabledFlag,
				NatType: func() *cato_models.SocketLanNatType { v := cato_models.SocketLanNatTypeDynamicPat; return &v }(),
			},
		},
	}

	debugPrint("UpdateRule Request - socketLanUpdateRuleInput (Kitchen Sink Parent)", socketLanUpdateRuleInput)

	updateParentResult, err := catoClient.PolicySocketLanUpdateRule(ctx, nil, socketLanUpdateRuleInput, accountId)
	if err != nil {
		fmt.Println("error updating parent Socket LAN rule: ", err)
		os.Exit(1)
	}

	debugPrint("UpdateRule Response (Parent)", updateParentResult)

	fmt.Printf("Parent Socket LAN rule updated with kitchen sink params:\n")
	if updateParentResult.Policy.SocketLan.UpdateRule.Rule != nil {
		updatedRule := updateParentResult.Policy.SocketLan.UpdateRule.Rule
		fmt.Printf("  ID: %s\n", updatedRule.Rule.ID)
		fmt.Printf("  Name: %s\n", updatedRule.Rule.Name)
		fmt.Printf("  Description: %s\n", updatedRule.Rule.Description)
		fmt.Printf("  Direction: %s\n", updatedRule.Rule.Direction)
		fmt.Printf("  Transport: %s\n", updatedRule.Rule.Transport)
	}

	// Query the policy to verify the update
	fmt.Printf("\nQuerying Socket LAN policy to verify parent rule update:\n")
	readAndDisplayPolicy(catoClient, ctx, accountId, parentRuleId)

	/////////////////////////////////////////////
	// Step 5: Update child Firewall rule with KITCHEN SINK params
	/////////////////////////////////////////////

	fmt.Printf("\n======================================\n")
	fmt.Printf("Step 5: Updating Child Firewall Rule (Kitchen Sink)\n")
	fmt.Printf("======================================\n")

	updatedChildName := "Socket LAN Firewall Sub-Rule Kitchen Sink"
	updatedChildDescription := "Updated with all possible parameters"
	childEnabledFlag := true
	childDirection := cato_models.SocketLanFirewallDirectionTo
	childAction := cato_models.SocketLanFirewallActionAllow

	firewallUpdateRuleInput := cato_models.SocketLanFirewallUpdateRuleInput{
		ID: childRuleId,
		Rule: &cato_models.SocketLanFirewallUpdateRuleDataInput{
			Name:        &updatedChildName,
			Description: &updatedChildDescription,
			Enabled:     &childEnabledFlag,
			Direction:   &childDirection,
			Action:      &childAction,
			// Source - all possible fields from kitchen sink (includes mac and site)
			Source: &cato_models.SocketLanFirewallSourceUpdateInput{
				Vlan: []scalars.Vlan{1},
				Mac:  []string{"a2:a9:97:1a:77:55"},
				IPRange: []*cato_models.IPAddressRangeInput{
					{From: "1.2.3.1", To: "1.2.3.4"},
				},
				Group: []*cato_models.GroupRefInput{
					{By: "ID", Input: "636881"},
				},
				Subnet: []string{"1.2.3.0/24"},
				Site: []*cato_models.SiteRefInput{
					{By: "ID", Input: "180279"},
				},
				NetworkInterface: []*cato_models.NetworkInterfaceRefInput{
					{By: "ID", Input: "216208"},
				},
				SystemGroup: []*cato_models.SystemGroupRefInput{
					{By: "ID", Input: "7S"},
				},
				Host: []*cato_models.HostRefInput{
					{By: "ID", Input: "2762912"},
				},
				IP: []string{"1.2.3.4"},
				GlobalIPRange: []*cato_models.GlobalIPRangeRefInput{
					{By: "ID", Input: "2528637"},
				},
				FloatingSubnet: []*cato_models.FloatingSubnetRefInput{
					{By: "ID", Input: "2528636"},
				},
				SiteNetworkSubnet: []*cato_models.SiteNetworkSubnetRefInput{
					{By: "ID", Input: "TjI3NjI5MTA="},
				},
			},
			// Destination - all possible fields from kitchen sink
			Destination: &cato_models.SocketLanFirewallDestinationUpdateInput{
				Vlan: []scalars.Vlan{1},
				IPRange: []*cato_models.IPAddressRangeInput{
					{From: "1.2.3.4", To: "1.2.3.5"},
				},
				Group: []*cato_models.GroupRefInput{
					{By: "ID", Input: "636881"},
				},
				Subnet: []string{"1.2.3.0/24"},
				Site: []*cato_models.SiteRefInput{
					{By: "ID", Input: "180279"},
				},
				NetworkInterface: []*cato_models.NetworkInterfaceRefInput{
					{By: "ID", Input: "216208"},
				},
				SystemGroup: []*cato_models.SystemGroupRefInput{
					{By: "ID", Input: "7S"},
				},
				Host: []*cato_models.HostRefInput{
					{By: "ID", Input: "2762912"},
				},
				IP: []string{"1.2.3.4"},
				GlobalIPRange: []*cato_models.GlobalIPRangeRefInput{
					{By: "ID", Input: "2528637"},
				},
				FloatingSubnet: []*cato_models.FloatingSubnetRefInput{
					{By: "ID", Input: "2528636"},
				},
				SiteNetworkSubnet: []*cato_models.SiteNetworkSubnetRefInput{
					{By: "ID", Input: "TjI3NjI5MTA="},
				},
			},
			// Application - all possible fields from kitchen sink
			Application: &cato_models.SocketLanFirewallApplicationUpdateInput{
				Application: []*cato_models.ApplicationRefInput{
					{By: "ID", Input: "adobe"},
				},
				CustomApp: []*cato_models.CustomApplicationRefInput{
					{By: "ID", Input: "CustomApp_11360_33973"},
				},
				Domain: []string{"something.com"},
				Fqdn:   []string{"www.something.com"},
				IP:     []string{"1.2.3.1"},
				Subnet: []string{"1.2.3.0/24"},
				IPRange: []*cato_models.IPAddressRangeInput{
					{From: "1.2.3.4", To: "1.2.3.5"},
				},
				GlobalIPRange: []*cato_models.GlobalIPRangeRefInput{
					{By: "ID", Input: "2528637"},
				},
			},
			// Service - simple, standard, and custom services
			Service: &cato_models.SocketLanFirewallServiceTypeUpdateInput{
				Simple: []*cato_models.SimpleServiceInput{
					{Name: "HTTP"},
				},
				Standard: []*cato_models.ServiceRefInput{
					{By: "ID", Input: "THREEPC"},
				},
				Custom: []*cato_models.CustomServiceInput{
					{
						Port:     []scalars.Port{"80"},
						Protocol: cato_models.IPProtocolTCP,
					},
					{
						PortRange: &cato_models.PortRangeInput{From: "80", To: "81"},
						Protocol:  cato_models.IPProtocolUDP,
					},
				},
			},
			// Tracking - event and alert with mailing list
			Tracking: &cato_models.PolicyTrackingUpdateInput{
				Event: &cato_models.PolicyRuleTrackingEventUpdateInput{
					Enabled: &childEnabledFlag,
				},
				Alert: &cato_models.PolicyRuleTrackingAlertUpdateInput{
					Enabled: &childEnabledFlag,
					Frequency: func() *cato_models.PolicyRuleTrackingFrequencyEnum {
						v := cato_models.PolicyRuleTrackingFrequencyEnumHourly
						return &v
					}(),
					MailingList: []*cato_models.SubscriptionMailingListRefInput{
						{By: "ID", Input: "13184"},
					},
				},
			},
		},
	}

	debugPrint("UpdateRule Request - firewallUpdateRuleInput (Kitchen Sink Child)", firewallUpdateRuleInput)

	updateChildResult, err := catoClient.PolicySocketLanFirewallUpdateRule(ctx, accountId, nil, firewallUpdateRuleInput)
	if err != nil {
		fmt.Println("error updating child Firewall rule: ", err)
		os.Exit(1)
	}

	debugPrint("UpdateRule Response (Child)", updateChildResult)

	fmt.Printf("Child Firewall rule updated with kitchen sink params:\n")
	if updateChildResult.Policy.SocketLan.Firewall.UpdateRule.Rule != nil {
		updatedRule := updateChildResult.Policy.SocketLan.Firewall.UpdateRule.Rule
		fmt.Printf("  ID: %s\n", updatedRule.Rule.ID)
		fmt.Printf("  Name: %s\n", updatedRule.Rule.Name)
		fmt.Printf("  Description: %s\n", updatedRule.Rule.Description)
		fmt.Printf("  Direction: %s\n", updatedRule.Rule.Direction)
		fmt.Printf("  Action: %s\n", updatedRule.Rule.Action)
	}

	// Query the policy to verify the update
	fmt.Printf("\nQuerying Socket LAN policy to verify child firewall rule update:\n")
	readAndDisplayPolicy(catoClient, ctx, accountId, parentRuleId)

	/////////////////////////////////////////////
	// Step 6: Update the Section
	/////////////////////////////////////////////

	fmt.Printf("\n======================================\n")
	fmt.Printf("Step 6: Updating Section Name\n")
	fmt.Printf("======================================\n")

	updatedSectionName := "SDK Example Section (Updated)"
	updateSectionInput := cato_models.PolicyUpdateSectionInput{
		ID: sectionId,
		Section: &cato_models.PolicyUpdateSectionInfoInput{
			Name: &updatedSectionName,
		},
	}

	debugPrint("UpdateSection Request", updateSectionInput)

	updateSectionResult, err := catoClient.PolicySocketLanUpdateSection(ctx, nil, updateSectionInput, accountId)
	if err != nil {
		fmt.Println("error updating section: ", err)
		os.Exit(1)
	}

	debugPrint("UpdateSection Response", updateSectionResult)

	fmt.Printf("Section updated successfully:\n")
	fmt.Printf("  ID: %s\n", updateSectionResult.Policy.SocketLan.UpdateSection.Section.Section.ID)
	fmt.Printf("  Name: %s\n", updateSectionResult.Policy.SocketLan.UpdateSection.Section.Section.Name)

	/////////////////////////////////////////////
	// Step 7: Publish the Socket LAN policy
	/////////////////////////////////////////////

	fmt.Printf("\n======================================\n")
	fmt.Printf("Step 7: Publishing Socket LAN Policy\n")
	fmt.Printf("======================================\n")

	publishResult, err := catoClient.PolicySocketLanPublishPolicyRevision(ctx, nil, nil, accountId)
	if err != nil {
		fmt.Println("error publishing Socket LAN policy revision: ", err)
		os.Exit(1)
	}

	debugPrint("PublishPolicyRevision Response", publishResult)

	fmt.Printf("Socket LAN policy revision published:\n")
	fmt.Printf("  Status: %s\n", publishResult.Policy.SocketLan.PublishPolicyRevision.Status)

	if len(publishResult.Policy.SocketLan.PublishPolicyRevision.Errors) > 0 {
		fmt.Printf("  Errors:\n")
		for _, err := range publishResult.Policy.SocketLan.PublishPolicyRevision.Errors {
			fmt.Printf("    - %s (Code: %s)\n", *err.ErrorMessage, *err.ErrorCode)
		}
	}

	/////////////////////////////////////////////
	// Step 8: Remove the child Firewall rule
	/////////////////////////////////////////////

	fmt.Printf("\n======================================\n")
	fmt.Printf("Step 8: Removing Child Firewall Rule\n")
	fmt.Printf("======================================\n")

	firewallRemoveRuleInput := cato_models.SocketLanFirewallRemoveRuleInput{
		ID: childRuleId,
	}

	debugPrint("RemoveRule Request - firewallRemoveRuleInput (Child)", firewallRemoveRuleInput)

	removeChildResult, err := catoClient.PolicySocketLanFirewallRemoveRule(ctx, accountId, nil, firewallRemoveRuleInput)
	if err != nil {
		fmt.Println("error removing child Firewall rule: ", err)
		os.Exit(1)
	}

	debugPrint("RemoveRule Response (Child)", removeChildResult)

	fmt.Printf("Child Firewall rule removed:\n")
	fmt.Printf("  Status: %s\n", removeChildResult.Policy.SocketLan.Firewall.RemoveRule.Status)

	/////////////////////////////////////////////
	// Step 9: Remove the parent Socket LAN rule
	/////////////////////////////////////////////

	fmt.Printf("\n======================================\n")
	fmt.Printf("Step 9: Removing Parent Socket LAN Rule\n")
	fmt.Printf("======================================\n")

	socketLanRemoveRuleInput := cato_models.SocketLanRemoveRuleInput{
		ID: parentRuleId,
	}

	debugPrint("RemoveRule Request - socketLanRemoveRuleInput (Parent)", socketLanRemoveRuleInput)

	removeParentResult, err := catoClient.PolicySocketLanRemoveRule(ctx, nil, socketLanRemoveRuleInput, accountId)
	if err != nil {
		fmt.Println("error removing parent Socket LAN rule: ", err)
		os.Exit(1)
	}

	debugPrint("RemoveRule Response (Parent)", removeParentResult)

	fmt.Printf("Parent Socket LAN rule removed:\n")
	fmt.Printf("  Status: %s\n", removeParentResult.Policy.SocketLan.RemoveRule.Status)

	/////////////////////////////////////////////
	// Step 10: Remove the Section
	/////////////////////////////////////////////

	fmt.Printf("\n======================================\n")
	fmt.Printf("Step 10: Removing Section\n")
	fmt.Printf("======================================\n")

	removeSectionInput := cato_models.PolicyRemoveSectionInput{
		ID: sectionId,
	}

	debugPrint("RemoveSection Request", removeSectionInput)

	removeSectionResult, err := catoClient.PolicySocketLanRemoveSection(ctx, nil, removeSectionInput, accountId)
	if err != nil {
		fmt.Println("error removing section: ", err)
		os.Exit(1)
	}

	debugPrint("RemoveSection Response", removeSectionResult)

	fmt.Printf("Section removed:\n")
	fmt.Printf("  Status: %s\n", removeSectionResult.Policy.SocketLan.RemoveSection.Status)

	/////////////////////////////////////////////
	// Step 11: Publish the final changes
	/////////////////////////////////////////////

	fmt.Printf("\n======================================\n")
	fmt.Printf("Step 11: Publishing Final Changes\n")
	fmt.Printf("======================================\n")

	finalPublishResult, err := catoClient.PolicySocketLanPublishPolicyRevision(ctx, nil, nil, accountId)
	if err != nil {
		fmt.Println("error publishing final Socket LAN policy revision: ", err)
		os.Exit(1)
	}

	debugPrint("Final PublishPolicyRevision Response", finalPublishResult)

	fmt.Printf("Final Socket LAN policy revision published:\n")
	fmt.Printf("  Status: %s\n", finalPublishResult.Policy.SocketLan.PublishPolicyRevision.Status)

	/////////////////////////////////////////////
	// Summary
	/////////////////////////////////////////////

	fmt.Printf("\n======================================\n")
	fmt.Printf("Complete Workflow Summary\n")
	fmt.Printf("======================================\n")
	fmt.Printf("1. Created Section (ID: %s)\n", sectionId)
	fmt.Printf("2. Created MINIMAL parent Socket LAN rule in section (ID: %s)\n", parentRuleId)
	fmt.Printf("3. Created MINIMAL child Firewall rule (ID: %s)\n", childRuleId)
	fmt.Printf("4. Updated parent rule with KITCHEN SINK params:\n")
	fmt.Printf("   - Site: site + group references\n")
	fmt.Printf("   - Source: vlan, ipRange, group, subnet, networkInterface, systemGroup, host, ip, globalIpRange, floatingSubnet, siteNetworkSubnet\n")
	fmt.Printf("   - Destination: vlan, ipRange, group, subnet, networkInterface, systemGroup, host, ip, globalIpRange, floatingSubnet, siteNetworkSubnet\n")
	fmt.Printf("   - Service: simple (HTTP), custom (port+protocol, portRange+protocol)\n")
	fmt.Printf("   - NAT: enabled with DYNAMIC_PAT\n")
	fmt.Printf("5. Updated child Firewall rule with KITCHEN SINK params:\n")
	fmt.Printf("   - Source: vlan, mac, ipRange, group, subnet, site, networkInterface, systemGroup, host, ip, globalIpRange, floatingSubnet, siteNetworkSubnet\n")
	fmt.Printf("   - Destination: vlan, ipRange, group, subnet, site, networkInterface, systemGroup, host, ip, globalIpRange, floatingSubnet, siteNetworkSubnet\n")
	fmt.Printf("   - Application: application (Adobe), customApp, domain, fqdn, ip, subnet, ipRange, globalIpRange\n")
	fmt.Printf("   - Service: simple (HTTP), standard (3PC), custom (port+protocol, portRange+protocol)\n")
	fmt.Printf("   - Tracking: event enabled, alert enabled with HOURLY frequency and mailing list\n")
	fmt.Printf("6. Updated Section name\n")
	fmt.Printf("7. Published policy revision\n")
	fmt.Printf("8. Removed child Firewall rule\n")
	fmt.Printf("9. Removed parent Socket LAN rule\n")
	fmt.Printf("10. Removed Section\n")
	fmt.Printf("11. Published final changes\n")
	fmt.Printf("\nAll operations completed successfully!\n")
}
