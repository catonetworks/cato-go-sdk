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

	/////////////////////////////////////////////
	// Create a new remote port forwarding rule //
	/////////////////////////////////////////////
	position := cato_models.PolicyRulePositionEnumLastInPolicy
	remotePortFwdAddRuleInput := cato_models.RemotePortFwdAddRuleInput{
		Rule: &cato_models.RemotePortFwdAddRuleDataInput{
			Name:        "Remote Port Forwarding Rule",
			Description: "Example remote port forwarding rule created via SDK",
			Enabled:     true,
			ForwardICMP: true,
			ExternalIP: &cato_models.AllocatedIPRefInput{
				By:    cato_models.ObjectRefByID,
				Input: "REPLACE_WITH_ALLOCATED_IP_ID",
			},
			ExternalPortRange: &cato_models.PortRangeInput{
				From: "443",
				To:   "443",
			},
			InternalIP: "192.168.1.2",
			InternalPortRange: &cato_models.PortRangeInput{
				From: "443",
				To:   "443",
			},
			RemoteIPs: &cato_models.RemotePortFwdRemoteIpsInput{
				GlobalIPRange: []*cato_models.GlobalIPRangeRefInput{},
				IP:            []string{"192.168.1.3"},
				IPRange: []*cato_models.IPAddressRangeInput{
					{
						From: "192.168.1.4",
						To:   "192.168.1.5",
					},
				},
				Subnet: []string{"192.168.1.3/32"},
			},
			RestrictionType: cato_models.RemotePortFwdRestrictionTypeAllowList,
			Tracking: &cato_models.PolicyRuleTrackingAlertInput{
				Enabled:   true,
				Frequency: cato_models.PolicyRuleTrackingFrequencyEnumHourly,
				MailingList: []*cato_models.SubscriptionMailingListRefInput{
					{
						By:    cato_models.ObjectRefByID,
						Input: "-100",
					},
				},
				SubscriptionGroup: []*cato_models.SubscriptionGroupRefInput{},
				Webhook:           []*cato_models.SubscriptionWebhookRefInput{},
			},
		},
		At: &cato_models.PolicyRulePositionInput{
			Position: &position,
		},
	}

	result, err := catoClient.PolicyRemotePortFwdAddRule(ctx, remotePortFwdAddRuleInput, accountId, nil)
	if err != nil {
		fmt.Println("error adding remote port forwarding rule: ", err)
		os.Exit(1)
	}

	// Print the result
	resultJson, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println("Remote port forwarding rule added successfully:")
	fmt.Println(string(resultJson))

	// Access specific fields
	if result.Policy.RemotePortFwd.AddRule.Rule != nil {
		rule := result.Policy.RemotePortFwd.AddRule.Rule
		fmt.Printf("\nRule Details:\n")
		fmt.Printf("ID: %s\n", rule.Rule.ID)
		fmt.Printf("Name: %s\n", rule.Rule.Name)
		fmt.Printf("Description: %s\n", rule.Rule.Description)
		fmt.Printf("Enabled: %t\n", rule.Rule.Enabled)
		fmt.Printf("Internal IP: %s\n", rule.Rule.InternalIP)
		fmt.Printf("Restriction Type: %s\n", rule.Rule.RestrictionType)
		fmt.Printf("Updated by: %s\n", rule.Audit.UpdatedBy)
		fmt.Printf("Updated time: %s\n", rule.Audit.UpdatedTime)
	}

	//////////////////////////////////////////
	// Read the new remote port forwarding rule //
	//////////////////////////////////////////

	if result.Policy.RemotePortFwd.AddRule.Rule != nil {
		rule := result.Policy.RemotePortFwd.AddRule.Rule
		ruleId := rule.Rule.ID
		ruleName := rule.Rule.Name

		fmt.Printf("\n======================================\n")
		fmt.Printf("Reading Remote Port Forwarding Policy\n")
		fmt.Printf("======================================\n")
		// Query the remote port forwarding policy to get the current state of all rules
		policyResult, err := catoClient.RemotePortFwdPolicy(ctx, accountId, nil)
		if err != nil {
			fmt.Println("error reading remote port forwarding policy: ", err)
			os.Exit(1)
		}

		// Display the rule details that we have from the creation response
		fmt.Printf("Rule ID: %s\n", ruleId)
		fmt.Printf("Rule Name: %s\n", ruleName)

		// Display properties if available
		if len(policyResult.Policy.RemotePortFwd.Policy.Rules) > 0 {
			fmt.Printf("Total rules in policy: %d\n", len(policyResult.Policy.RemotePortFwd.Policy.Rules))
			// Look for our specific rule
			for _, rule := range policyResult.Policy.RemotePortFwd.Policy.Rules {
				if rule.Rule.ID == ruleId {
					fmt.Printf("Found our rule: %s\n", rule.Rule.Name)
					break
				}
			}
		} else {
			fmt.Printf("No rules found in policy\n")
		}

		/////////////////////////////////////////////
		// Update the remote port forwarding rule //
		/////////////////////////////////////////////

		fmt.Printf("\n======================================\n")
		fmt.Printf("Updating Remote Port Forwarding Rule\n")
		fmt.Printf("======================================\n")

		// Create update input with modified values
		updatedName := "Updated Remote Port Forwarding Rule"
		updatedDescription := "Updated remote port forwarding rule description via SDK"
		forwardICMP := false

		remotePortFwdUpdateRuleInput := cato_models.RemotePortFwdUpdateRuleInput{
			ID: ruleId,
			Rule: &cato_models.RemotePortFwdUpdateRuleDataInput{
				Name:        &updatedName,
				Description: &updatedDescription,
				ForwardICMP: &forwardICMP,
			},
		}

		// Perform the update
		updateResult, err := catoClient.PolicyRemotePortFwdUpdateRule(ctx, remotePortFwdUpdateRuleInput, accountId, nil)
		if err != nil {
			fmt.Println("error updating remote port forwarding rule: ", err)
			os.Exit(1)
		}

		// Print the update result
		updateResultJson, _ := json.MarshalIndent(updateResult, "", "  ")
		fmt.Println("Remote port forwarding rule updated successfully:")
		fmt.Println(string(updateResultJson))

		// Access specific fields from update result
		if updateResult.Policy.RemotePortFwd.UpdateRule.Rule != nil {
			updatedRule := updateResult.Policy.RemotePortFwd.UpdateRule.Rule
			fmt.Printf("\nUpdated Rule Details:\n")
			fmt.Printf("ID: %s\n", updatedRule.Rule.ID)
			fmt.Printf("Name: %s\n", updatedRule.Rule.Name)
			fmt.Printf("Description: %s\n", updatedRule.Rule.Description)
			fmt.Printf("Enabled: %t\n", updatedRule.Rule.Enabled)
			fmt.Printf("Forward ICMP: %t\n", updatedRule.Rule.ForwardICMP)
			fmt.Printf("Updated by: %s\n", updatedRule.Audit.UpdatedBy)
			fmt.Printf("Updated time: %s\n", updatedRule.Audit.UpdatedTime)
		}

		// Check for any update errors
		if len(updateResult.Policy.RemotePortFwd.UpdateRule.Errors) > 0 {
			fmt.Printf("\nUpdate Errors:\n")
			for _, err := range updateResult.Policy.RemotePortFwd.UpdateRule.Errors {
				fmt.Printf("- %s (Code: %s)\n", *err.ErrorMessage, *err.ErrorCode)
			}
		} else {
			fmt.Printf("\n✓ Remote port forwarding rule updated successfully!\n")
			fmt.Printf("  - Name changed to: %s\n", updatedName)
			fmt.Printf("  - Description updated\n")
			fmt.Printf("  - Forward ICMP disabled\n")
		}

		/////////////////////////////////////////////
		// Delete the remote port forwarding rule //
		/////////////////////////////////////////////

		fmt.Printf("\n======================================\n")
		fmt.Printf("Deleting Remote Port Forwarding Rule\n")
		fmt.Printf("======================================\n")

		// Create remove input
		remotePortFwdRemoveRuleInput := cato_models.RemotePortFwdRemoveRuleInput{
			ID: ruleId,
		}

		// Perform the delete operation
		deleteResult, err := catoClient.PolicyRemotePortFwdRemoveRule(ctx, remotePortFwdRemoveRuleInput, accountId, nil)
		if err != nil {
			fmt.Println("error deleting remote port forwarding rule: ", err)
			os.Exit(1)
		}

		// Print the delete result
		deleteResultJson, _ := json.MarshalIndent(deleteResult, "", "  ")
		fmt.Println("Remote port forwarding rule deletion initiated:")
		fmt.Println(string(deleteResultJson))

		// Access specific fields from delete result
		fmt.Printf("\nDeletion Status: %s\n", deleteResult.Policy.RemotePortFwd.RemoveRule.Status)

		// Check for any delete errors
		if len(deleteResult.Policy.RemotePortFwd.RemoveRule.Errors) > 0 {
			fmt.Printf("\nDelete Errors:\n")
			for _, err := range deleteResult.Policy.RemotePortFwd.RemoveRule.Errors {
				fmt.Printf("- %s (Code: %s)\n", *err.ErrorMessage, *err.ErrorCode)
			}
		} else {
			fmt.Printf("\n✓ Remote port forwarding rule deletion completed successfully!\n")
			fmt.Printf("  - Rule ID %s has been marked for removal\n", ruleId)
			fmt.Printf("  - The rule will be removed from the policy after publishing\n")
		}

		/////////////////////////////////////////////
		// Publish the remote port forwarding policy //
		/////////////////////////////////////////////

		publishResult, err := catoClient.PolicyRemotePortFwdPublishPolicyRevision(ctx, accountId, nil, nil)
		if err != nil {
			fmt.Println("error publishing remote port forwarding policy revision: ", err)
			os.Exit(1)
		}

		// Print the publish result
		publishResultJson, _ := json.MarshalIndent(publishResult, "", "  ")
		fmt.Println("\nRemote port forwarding policy revision published successfully:")
		fmt.Println(string(publishResultJson))

		// Access specific fields
		fmt.Printf("\nPublish Status: %s\n", publishResult.Policy.RemotePortFwd.PublishPolicyRevision.Status)

		// Check for any errors
		if len(publishResult.Policy.RemotePortFwd.PublishPolicyRevision.Errors) > 0 {
			fmt.Printf("\nPublish Errors:\n")
			errorsJson, _ := json.MarshalIndent(publishResult.Policy.RemotePortFwd.PublishPolicyRevision.Errors, "", "  ")
			fmt.Println(string(errorsJson))
		} else {
			fmt.Printf("\nThe remote port forwarding policy revision has been successfully published and is now live.\n")
			fmt.Printf("All changes made to the draft revision are now active in the production environment.\n")
			fmt.Printf("\n🎉 Complete CRUD workflow finished successfully!\n")
			fmt.Printf("  ✓ Created remote port forwarding rule: %s\n", "Remote Port Forwarding Rule")
			fmt.Printf("  ✓ Read rule from policy\n")
			fmt.Printf("  ✓ Updated rule name to: %s\n", updatedName)
			fmt.Printf("  ✓ Deleted rule with ID: %s\n", ruleId)
			fmt.Printf("  ✓ Published all changes to production\n")
		}

	}
}
