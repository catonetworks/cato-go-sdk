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
	// Create a new TLS inspection rule //
	/////////////////////////////////////////////
	position := cato_models.PolicyRulePositionEnumLastInPolicy
	tlsInspectAddRuleInput := cato_models.TLSInspectAddRuleInput{
		Rule: &cato_models.TLSInspectAddRuleDataInput{
			Name:                       "Example TLS Inspection Rule",
			Description:                "Example TLS inspection rule created via SDK",
			Enabled:                    true,
		Action:                     cato_models.TLSInspectActionInspect,
		UntrustedCertificateAction: cato_models.TLSInspectUntrustedCertificateActionBlock,
		ConnectionOrigin:           cato_models.ConnectionOriginEnumAny,
		Source: &cato_models.TLSInspectSourceInput{
			FloatingSubnet:    []*cato_models.FloatingSubnetRefInput{},
			GlobalIPRange:     []*cato_models.GlobalIPRangeRefInput{},
			Group:             []*cato_models.GroupRefInput{},
			Host:              []*cato_models.HostRefInput{},
			IP:                []string{},
			IPRange:           []*cato_models.IPAddressRangeInput{},
			NetworkInterface:  []*cato_models.NetworkInterfaceRefInput{},
			Site:              []*cato_models.SiteRefInput{},
			SiteNetworkSubnet: []*cato_models.SiteNetworkSubnetRefInput{},
			Subnet:            []string{},
			SystemGroup:       []*cato_models.SystemGroupRefInput{},
			User:              []*cato_models.UserRefInput{},
			UsersGroup:        []*cato_models.UsersGroupRefInput{},
		},
		Application: &cato_models.TLSInspectApplicationInput{
			AppCategory:        []*cato_models.ApplicationCategoryRefInput{},
			Application:        []*cato_models.ApplicationRefInput{},
			Country:            []*cato_models.CountryRefInput{},
			CustomApp:          []*cato_models.CustomApplicationRefInput{},
			CustomCategory:     []*cato_models.CustomCategoryRefInput{},
			CustomService:      []*cato_models.CustomServiceInput{},
			CustomServiceIP:    []*cato_models.CustomServiceIPInput{},
			Domain:             []string{},
			Fqdn:               []string{},
			GlobalIPRange:      []*cato_models.GlobalIPRangeRefInput{},
			IP:                 []string{},
			IPRange:            []*cato_models.IPAddressRangeInput{},
			RemoteAsn:          []scalars.Asn32{},
			Service:            []*cato_models.ServiceRefInput{},
			Subnet:             []string{},
			TLSInspectCategory: []cato_models.TLSInspectSystemCategory{},
		},
		Country:              []*cato_models.CountryRefInput{},
		DevicePostureProfile: []*cato_models.DeviceProfileRefInput{},
		Platform:             []cato_models.OperatingSystem{},
		},
		At: &cato_models.PolicyRulePositionInput{
			Position: &position,
		},
	}

	result, err := catoClient.PolicyTLSInspectAddRule(ctx, tlsInspectAddRuleInput, accountId)
	if err != nil {
		fmt.Println("error adding TLS inspection rule: ", err)
		os.Exit(1)
	}

	// Print the result
	resultJson, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println("TLS inspection rule added successfully:")
	fmt.Println(string(resultJson))

	// Access specific fields
	if result.Policy.TLSInspect.AddRule.Rule != nil {
		rule := result.Policy.TLSInspect.AddRule.Rule
		fmt.Printf("\nRule Details:\n")
		fmt.Printf("ID: %s\n", rule.Rule.ID)
		fmt.Printf("Name: %s\n", rule.Rule.Name)
		fmt.Printf("Description: %s\n", rule.Rule.Description)
		fmt.Printf("Enabled: %t\n", rule.Rule.Enabled)
		fmt.Printf("Action: %s\n", rule.Rule.Action)
		fmt.Printf("Untrusted Certificate Action: %s\n", rule.Rule.UntrustedCertificateAction)
		fmt.Printf("Updated by: %s\n", rule.Audit.UpdatedBy)
		fmt.Printf("Updated time: %s\n", rule.Audit.UpdatedTime)
	}

	//////////////////////////////////////////
	// Read the new TLS inspection rule //
	//////////////////////////////////////////

	if result.Policy.TLSInspect.AddRule.Rule != nil {
		rule := result.Policy.TLSInspect.AddRule.Rule
		ruleId := rule.Rule.ID
		ruleName := rule.Rule.Name

		fmt.Printf("\n======================================\n")
		fmt.Printf("Reading TLS Inspection Policy\n")
		fmt.Printf("======================================\n")
		// Query the TLS inspection policy to get the current state of all rules
		policyResult, err := catoClient.Tlsinspectpolicy(ctx, accountId)
		if err != nil {
			fmt.Println("error reading TLS inspection policy: ", err)
			os.Exit(1)
		}

		// Display the rule details that we have from the creation response
		fmt.Printf("Rule ID: %s\n", ruleId)
		fmt.Printf("Rule Name: %s\n", ruleName)

		// Display properties if available
		if len(policyResult.Policy.TLSInspect.Policy.Sections) > 0 {
			fmt.Printf("Total sections in policy: %d\n", len(policyResult.Policy.TLSInspect.Policy.Sections))
			// Look for our specific rule
			for _, rule := range policyResult.Policy.TLSInspect.Policy.Rules {
				if rule.Rule.ID == ruleId {
					fmt.Printf("Found our rule: %s\n", rule.Rule.Name)
					break
				}
			}
		} else {
			fmt.Printf("No rules found in policy\n")
		}

		/////////////////////////////////////////////
		// Update the TLS inspection rule      //
		/////////////////////////////////////////////

		fmt.Printf("\n======================================\n")
		fmt.Printf("Updating TLS Inspection Rule\n")
		fmt.Printf("======================================\n")

		// Create update input with modified values
		updatedName := "Updated Example TLS Inspection Rule"
		updatedDescription := "Updated TLS inspection rule description via SDK"
		enabledFlag := true
		action := cato_models.TLSInspectActionBypass
		untrustedCertAction := cato_models.TLSInspectUntrustedCertificateActionAllow
		connectionOrigin := cato_models.ConnectionOriginEnumRemote

		tlsInspectUpdateRuleInput := cato_models.TLSInspectUpdateRuleInput{
			ID: ruleId,
			Rule: &cato_models.TLSInspectUpdateRuleDataInput{
				Name:                       &updatedName,
				Description:                &updatedDescription,
				Enabled:                    &enabledFlag,
				Action:                     &action,
				UntrustedCertificateAction: &untrustedCertAction,
				ConnectionOrigin:           &connectionOrigin,
				// Keep the same source and application settings
				Source:               &cato_models.TLSInspectSourceUpdateInput{},
				Application:          &cato_models.TLSInspectApplicationUpdateInput{},
				Country:              []*cato_models.CountryRefInput{},
				DevicePostureProfile: []*cato_models.DeviceProfileRefInput{},
				Platform:             []cato_models.OperatingSystem{},
			},
		}

		// Perform the update
		updateResult, err := catoClient.PolicyTLSInspectUpdateRule(ctx, tlsInspectUpdateRuleInput, accountId)
		if err != nil {
			fmt.Println("error updating TLS inspection rule: ", err)
			os.Exit(1)
		}

		// Print the update result
		updateResultJson, _ := json.MarshalIndent(updateResult, "", "  ")
		fmt.Println("TLS inspection rule updated successfully:")
		fmt.Println(string(updateResultJson))

		// Access specific fields from update result
		if updateResult.Policy.TLSInspect.UpdateRule.Rule != nil {
			updatedRule := updateResult.Policy.TLSInspect.UpdateRule.Rule
			fmt.Printf("\nUpdated Rule Details:\n")
			fmt.Printf("ID: %s\n", updatedRule.Rule.ID)
			fmt.Printf("Name: %s\n", updatedRule.Rule.Name)
			fmt.Printf("Description: %s\n", updatedRule.Rule.Description)
			fmt.Printf("Enabled: %t\n", updatedRule.Rule.Enabled)
			fmt.Printf("Action: %s\n", updatedRule.Rule.Action)
			fmt.Printf("Untrusted Certificate Action: %s\n", updatedRule.Rule.UntrustedCertificateAction)
			fmt.Printf("Connection Origin: %s\n", updatedRule.Rule.ConnectionOrigin)
			fmt.Printf("Updated by: %s\n", updatedRule.Audit.UpdatedBy)
			fmt.Printf("Updated time: %s\n", updatedRule.Audit.UpdatedTime)
		}

		// Check for any update errors
		if len(updateResult.Policy.TLSInspect.UpdateRule.Errors) > 0 {
			fmt.Printf("\nUpdate Errors:\n")
			for _, err := range updateResult.Policy.TLSInspect.UpdateRule.Errors {
				fmt.Printf("- %s (Code: %s)\n", *err.ErrorMessage, *err.ErrorCode)
			}
		} else {
			fmt.Printf("\n✓ TLS inspection rule updated successfully!\n")
			fmt.Printf("  - Name changed to: %s\n", updatedName)
			fmt.Printf("  - Description updated\n")
			fmt.Printf("  - Action changed from INSPECT to BYPASS\n")
			fmt.Printf("  - Untrusted certificate action changed from BLOCK to ALLOW\n")
			fmt.Printf("  - Connection origin changed from ANY to REMOTE\n")
		}

		/////////////////////////////////////////////
		// Delete the TLS inspection rule      //
		/////////////////////////////////////////////

		fmt.Printf("\n======================================\n")
		fmt.Printf("Deleting TLS Inspection Rule\n")
		fmt.Printf("======================================\n")

		// Create remove input
		tlsInspectRemoveRuleInput := cato_models.TLSInspectRemoveRuleInput{
			ID: ruleId,
		}

		// Perform the delete operation
		deleteResult, err := catoClient.PolicyTLSInspectRemoveRule(ctx, tlsInspectRemoveRuleInput, accountId)
		if err != nil {
			fmt.Println("error deleting TLS inspection rule: ", err)
			os.Exit(1)
		}

		// Print the delete result
		deleteResultJson, _ := json.MarshalIndent(deleteResult, "", "  ")
		fmt.Println("TLS inspection rule deletion initiated:")
		fmt.Println(string(deleteResultJson))

		// Access specific fields from delete result
		fmt.Printf("\nDeletion Status: %s\n", deleteResult.Policy.TLSInspect.RemoveRule.Status)

		// Check for any delete errors
		if len(deleteResult.Policy.TLSInspect.RemoveRule.Errors) > 0 {
			fmt.Printf("\nDelete Errors:\n")
			for _, err := range deleteResult.Policy.TLSInspect.RemoveRule.Errors {
				fmt.Printf("- %s (Code: %s)\n", *err.ErrorMessage, *err.ErrorCode)
			}
		} else {
			fmt.Printf("\n✓ TLS inspection rule deletion completed successfully!\n")
			fmt.Printf("  - Rule ID %s has been marked for removal\n", ruleId)
			fmt.Printf("  - The rule will be removed from the policy after publishing\n")
		}

		/////////////////////////////////////////////
		// Publish the TLS inspection policy   //
		/////////////////////////////////////////////

	publishResult, err := catoClient.PolicyTLSInspectPublishPolicyRevision(ctx, accountId)
		if err != nil {
			fmt.Println("error publishing TLS inspection policy revision: ", err)
			os.Exit(1)
		}

		// Print the publish result
		publishResultJson, _ := json.MarshalIndent(publishResult, "", "  ")
		fmt.Println("\nTLS inspection policy revision published successfully:")
		fmt.Println(string(publishResultJson))

		// Access specific fields
		fmt.Printf("\nPublish Status: %s\n", publishResult.Policy.TLSInspect.PublishPolicyRevision.Status)

		// Check for any errors
		if len(publishResult.Policy.TLSInspect.PublishPolicyRevision.Errors) > 0 {
			fmt.Printf("\nPublish Errors:\n")
			errorsJson, _ := json.MarshalIndent(publishResult.Policy.TLSInspect.PublishPolicyRevision.Errors, "", "  ")
			fmt.Println(string(errorsJson))
		} else {
			fmt.Printf("\nThe TLS inspection policy revision has been successfully published and is now live.\n")
			fmt.Printf("All changes made to the draft revision are now active in the production environment.\n")
			fmt.Printf("\n🎉 Complete CRUD workflow finished successfully!\n")
			fmt.Printf("  ✓ Created TLS inspection rule: %s\n", "Example TLS Inspection Rule")
			fmt.Printf("  ✓ Read rule from policy\n")
			fmt.Printf("  ✓ Updated rule name to: %s\n", updatedName)
			fmt.Printf("  ✓ Deleted rule with ID: %s\n", ruleId)
			fmt.Printf("  ✓ Published all changes to production\n")
		}

	}
}
