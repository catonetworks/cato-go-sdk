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

	//////////////////////////////////////
	// Create a new WAN network rule //
	//////////////////////////////////////
	position := cato_models.PolicyRulePositionEnumLastInPolicy
	wanNetworkAddRuleInput := cato_models.WanNetworkAddRuleInput{
		Rule: &cato_models.WanNetworkAddRuleDataInput{
			Name:        "Example WAN Rule",
			Description: "Example WAN network rule created via SDK",
			RuleType:    cato_models.WanNetworkRuleTypeWan,
			RouteType:   cato_models.WanNetworkRuleRouteTypeOptimized,
			Enabled:     true,
			Source:      &cato_models.WanNetworkRuleSourceInput{},
			Destination: &cato_models.WanNetworkRuleDestinationInput{},
			Application: &cato_models.WanNetworkRuleApplicationInput{},
			BandwidthPriority: &cato_models.BandwidthManagementRefInput{
				By:    cato_models.ObjectRefBy("NAME"),
				Input: "10",
			},
			Configuration: &cato_models.WanNetworkRuleConfigurationInput{
				ActiveTCPAcceleration: true,
				PacketLossMitigation:  false,
				PreserveSourcePort:    true,
				PrimaryTransport: &cato_models.WanNetworkRuleTransportInput{
					TransportType:          cato_models.WanNetworkRuleTransportType("WAN"),
					PrimaryInterfaceRole:   "WAN1",
					SecondaryInterfaceRole: "AUTOMATIC",
				},
				SecondaryTransport: &cato_models.WanNetworkRuleTransportInput{
					TransportType:          cato_models.WanNetworkRuleTransportType("AUTOMATIC"),
					PrimaryInterfaceRole:   "AUTOMATIC",
					SecondaryInterfaceRole: "AUTOMATIC",
				},
			},
		},
		At: &cato_models.PolicyRulePositionInput{
			Position: &position,
		},
	}

	result, err := catoClient.PolicyWanNetworkAddRule(ctx, wanNetworkAddRuleInput, accountId)
	if err != nil {
		fmt.Println("error adding WAN network rule: ", err)
		os.Exit(1)
	}

	// Print the result
	resultJson, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println("WAN network rule added successfully:")
	fmt.Println(string(resultJson))

	// Access specific fields
	if result.Policy.WanNetwork.AddRule.Rule != nil {
		rule := result.Policy.WanNetwork.AddRule.Rule
		fmt.Printf("\nRule Details:\n")
		fmt.Printf("ID: %s\n", rule.Rule.ID)
		fmt.Printf("Name: %s\n", rule.Rule.Name)
		fmt.Printf("Description: %s\n", rule.Rule.Description)
		fmt.Printf("Enabled: %t\n", rule.Rule.Enabled)
		fmt.Printf("Rule Type: %s\n", rule.Rule.RuleType)
		fmt.Printf("Updated by: %s\n", rule.Audit.UpdatedBy)
		fmt.Printf("Updated time: %s\n", rule.Audit.UpdatedTime)
	}

	///////////////////////////////////
	// Read the new WAN network rule //
	///////////////////////////////////

	if result.Policy.WanNetwork.AddRule.Rule != nil {
		rule := result.Policy.WanNetwork.AddRule.Rule
		ruleId := rule.Rule.ID
		ruleName := rule.Rule.Name

		fmt.Printf("\n======================================\n")
		fmt.Printf("Reading WAN Network Section\n")
		fmt.Printf("======================================\n")
		// Query the WAN network policy to get the current state of all sections
		policyResult, err := catoClient.WanNetworkPolicy(ctx, accountId)
		if err != nil {
			fmt.Println("error reading WAN network policy: ", err)
			os.Exit(1)
		}

		// Display the section details that we have from the creation response
		fmt.Printf("Rule ID: %s\n", ruleId)
		fmt.Printf("Rule Name: %s\n", ruleName)

		// Display properties if available
		if len(policyResult.Policy.WanNetwork.Policy.Sections) > 0 {
			fmt.Printf("Total sections in policy: %d\n", len(policyResult.Policy.WanNetwork.Policy.Sections))
			// Look for our specific section
			for _, rule := range policyResult.Policy.WanNetwork.Policy.Rules {
				if rule.Rule.ID == ruleId {
					fmt.Printf("Found our rule: %s\n", rule.Rule.Name)
					break
				}
			}
		} else {
			fmt.Printf("No rules found in policy\n")
		}

		//////////////////////////////////////
		// Update the WAN network rule      //
		//////////////////////////////////////

		fmt.Printf("\n======================================\n")
		fmt.Printf("Updating WAN Network Rule\n")
		fmt.Printf("======================================\n")

		// Create update input with modified values
		updatedName := "Updated Example WAN Rule"
		updatedDescription := "Updated WAN network rule description via SDK"
		enabledFlag := true
		ruleType := cato_models.WanNetworkRuleTypeWan
		routeType := cato_models.WanNetworkRuleRouteTypeOptimized
		activeTCPAcceleration := false
		packetLossMitigation := true
		preserveSourcePort := true
		primaryInterfaceRole := cato_models.WanNetworkRuleInterfaceRole("WAN2")
		secondaryInterfaceRole := cato_models.WanNetworkRuleInterfaceRole("AUTOMATIC")
		automaticInterfaceRole := cato_models.WanNetworkRuleInterfaceRole("AUTOMATIC")
		transportTypeWan := cato_models.WanNetworkRuleTransportType("WAN")
		transportTypeAuto := cato_models.WanNetworkRuleTransportType("AUTOMATIC")

		wanNetworkUpdateRuleInput := cato_models.WanNetworkUpdateRuleInput{
			ID: ruleId,
			Rule: &cato_models.WanNetworkUpdateRuleDataInput{
				Name:        &updatedName,
				Description: &updatedDescription,
				Enabled:     &enabledFlag,
				RuleType:    &ruleType,
				RouteType:   &routeType,
				// Keep the same source, destination, and application settings
				Source:      &cato_models.WanNetworkRuleSourceUpdateInput{},
				Destination: &cato_models.WanNetworkRuleDestinationUpdateInput{},
				Application: &cato_models.WanNetworkRuleApplicationUpdateInput{},
				// Update bandwidth priority
				BandwidthPriority: &cato_models.BandwidthManagementRefInput{
					By:    cato_models.ObjectRefBy("NAME"),
					Input: "20", // Changed from "10" to "20"
				},
				// Update configuration
				Configuration: &cato_models.WanNetworkRuleConfigurationUpdateInput{
					ActiveTCPAcceleration: &activeTCPAcceleration, // Changed from true to false
					PacketLossMitigation:  &packetLossMitigation,  // Changed from false to true
					PreserveSourcePort:    &preserveSourcePort,
					PrimaryTransport: &cato_models.WanNetworkRuleTransportUpdateInput{
						TransportType:          &transportTypeWan,
						PrimaryInterfaceRole:   &primaryInterfaceRole, // Changed from "WAN1" to "WAN2"
						SecondaryInterfaceRole: &secondaryInterfaceRole,
					},
					SecondaryTransport: &cato_models.WanNetworkRuleTransportUpdateInput{
						TransportType:          &transportTypeAuto,
						PrimaryInterfaceRole:   &automaticInterfaceRole,
						SecondaryInterfaceRole: &automaticInterfaceRole,
					},
				},
			},
		}

		// Perform the update
		updateResult, err := catoClient.PolicyWanNetworkUpdateRule(ctx, wanNetworkUpdateRuleInput, accountId)
		if err != nil {
			fmt.Println("error updating WAN network rule: ", err)
			os.Exit(1)
		}

		// Print the update result
		updateResultJson, _ := json.MarshalIndent(updateResult, "", "  ")
		fmt.Println("WAN network rule updated successfully:")
		fmt.Println(string(updateResultJson))

		// Access specific fields from update result
		if updateResult.Policy.WanNetwork.UpdateRule.Rule != nil {
			updatedRule := updateResult.Policy.WanNetwork.UpdateRule.Rule
			fmt.Printf("\nUpdated Rule Details:\n")
			fmt.Printf("ID: %s\n", updatedRule.Rule.ID)
			fmt.Printf("Name: %s\n", updatedRule.Rule.Name)
			fmt.Printf("Description: %s\n", updatedRule.Rule.Description)
			fmt.Printf("Enabled: %t\n", updatedRule.Rule.Enabled)
			fmt.Printf("Rule Type: %s\n", updatedRule.Rule.RuleType)
			fmt.Printf("Route Type: %s\n", updatedRule.Rule.RouteType)
			fmt.Printf("Updated by: %s\n", updatedRule.Audit.UpdatedBy)
			fmt.Printf("Updated time: %s\n", updatedRule.Audit.UpdatedTime)
		}

		// Check for any update errors
		if len(updateResult.Policy.WanNetwork.UpdateRule.Errors) > 0 {
			fmt.Printf("\nUpdate Errors:\n")
			for _, err := range updateResult.Policy.WanNetwork.UpdateRule.Errors {
				fmt.Printf("- %s (Code: %s)\n", *err.ErrorMessage, *err.ErrorCode)
			}
		} else {
			fmt.Printf("\n✓ WAN network rule updated successfully!\n")
			fmt.Printf("  - Name changed to: %s\n", updatedName)
			fmt.Printf("  - Description updated\n")
			fmt.Printf("  - Bandwidth priority changed from '10' to '20'\n")
			fmt.Printf("  - TCP acceleration disabled\n")
			fmt.Printf("  - Packet loss mitigation enabled\n")
			fmt.Printf("  - Primary interface changed from WAN1 to WAN2\n")
		}

		//////////////////////////////////////
		// Delete the WAN network rule      //
		//////////////////////////////////////

		fmt.Printf("\n======================================\n")
		fmt.Printf("Deleting WAN Network Rule\n")
		fmt.Printf("======================================\n")

		// Create remove input
		wanNetworkRemoveRuleInput := cato_models.WanNetworkRemoveRuleInput{
			ID: ruleId,
		}

		// Perform the delete operation
		deleteResult, err := catoClient.PolicyWanNetworkRemoveRule(ctx, wanNetworkRemoveRuleInput, accountId)
		if err != nil {
			fmt.Println("error deleting WAN network rule: ", err)
			os.Exit(1)
		}

		// Print the delete result
		deleteResultJson, _ := json.MarshalIndent(deleteResult, "", "  ")
		fmt.Println("WAN network rule deletion initiated:")
		fmt.Println(string(deleteResultJson))

		// Access specific fields from delete result
		fmt.Printf("\nDeletion Status: %s\n", deleteResult.Policy.WanNetwork.RemoveRule.Status)

		// Check for any delete errors
		if len(deleteResult.Policy.WanNetwork.RemoveRule.Errors) > 0 {
			fmt.Printf("\nDelete Errors:\n")
			for _, err := range deleteResult.Policy.WanNetwork.RemoveRule.Errors {
				fmt.Printf("- %s (Code: %s)\n", *err.ErrorMessage, *err.ErrorCode)
			}
		} else {
			fmt.Printf("\n✓ WAN network rule deletion completed successfully!\n")
			fmt.Printf("  - Rule ID %s has been marked for removal\n", ruleId)
			fmt.Printf("  - The rule will be removed from the policy after publishing\n")
		}

		//////////////////////////////////////
		// Publish the WAN network policy   //
		//////////////////////////////////////

		publishResult, err := catoClient.PolicyWanNetworkPublishPolicyRevision(ctx, accountId)
		if err != nil {
			fmt.Println("error publishing WAN network policy revision: ", err)
			os.Exit(1)
		}

		// Print the publish result
		publishResultJson, _ := json.MarshalIndent(publishResult, "", "  ")
		fmt.Println("\nWAN network policy revision published successfully:")
		fmt.Println(string(publishResultJson))

		// Access specific fields
		fmt.Printf("\nPublish Status: %s\n", publishResult.Policy.WanNetwork.PublishPolicyRevision.Status)

		// Check for any errors
		if len(publishResult.Policy.WanNetwork.PublishPolicyRevision.Errors) > 0 {
			fmt.Printf("\nPublish Errors:\n")
			errorsJson, _ := json.MarshalIndent(publishResult.Policy.WanNetwork.PublishPolicyRevision.Errors, "", "  ")
			fmt.Println(string(errorsJson))
		} else {
			fmt.Printf("\nThe WAN network policy revision has been successfully published and is now live.\n")
			fmt.Printf("All changes made to the draft revision are now active in the production environment.\n")
			fmt.Printf("\n🎉 Complete CRUD workflow finished successfully!\n")
			fmt.Printf("  ✓ Created WAN network rule: %s\n", "Example WAN Rule")
			fmt.Printf("  ✓ Read rule from policy\n")
			fmt.Printf("  ✓ Updated rule name to: %s\n", updatedName)
			fmt.Printf("  ✓ Deleted rule with ID: %s\n", ruleId)
			fmt.Printf("  ✓ Published all changes to production\n")
		}

	}
}
