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
	// Create a new Socket LAN rule (parent)  //
	/////////////////////////////////////////////

	fmt.Printf("======================================\n")
	fmt.Printf("Creating Socket LAN Rule (Parent)\n")
	fmt.Printf("======================================\n")

	position := cato_models.PolicyRulePositionEnumLastInPolicy

	socketLanAddRuleInput := cato_models.SocketLanAddRuleInput{
		At: &cato_models.PolicyRulePositionInput{
			Position: &position,
		},
		Rule: &cato_models.SocketLanAddRuleDataInput{
			Name:        "Example Socket LAN Rule",
			Description: "Socket LAN rule created via SDK",
			Enabled:     true,
			Direction:   cato_models.SocketLanDirectionTo,
			Transport:   cato_models.SocketLanTransportTypeLan,
			Source: &cato_models.SocketLanSourceInput{
				IP:                []string{},
				Subnet:            []string{"192.168.1.0/24"},
				Vlan:              []scalars.Vlan{},
				Host:              []*cato_models.HostRefInput{},
				Group:             []*cato_models.GroupRefInput{},
				SystemGroup:       []*cato_models.SystemGroupRefInput{},
				NetworkInterface:  []*cato_models.NetworkInterfaceRefInput{},
				GlobalIPRange:     []*cato_models.GlobalIPRangeRefInput{},
				FloatingSubnet:    []*cato_models.FloatingSubnetRefInput{},
				SiteNetworkSubnet: []*cato_models.SiteNetworkSubnetRefInput{},
				IPRange:           []*cato_models.IPAddressRangeInput{},
			},
			Destination: &cato_models.SocketLanDestinationInput{
				IP:                []string{},
				Subnet:            []string{"10.0.0.0/24"},
				Vlan:              []scalars.Vlan{},
				Host:              []*cato_models.HostRefInput{},
				Group:             []*cato_models.GroupRefInput{},
				SystemGroup:       []*cato_models.SystemGroupRefInput{},
				NetworkInterface:  []*cato_models.NetworkInterfaceRefInput{},
				GlobalIPRange:     []*cato_models.GlobalIPRangeRefInput{},
				FloatingSubnet:    []*cato_models.FloatingSubnetRefInput{},
				SiteNetworkSubnet: []*cato_models.SiteNetworkSubnetRefInput{},
				IPRange:           []*cato_models.IPAddressRangeInput{},
			},
			Site: &cato_models.SocketLanSiteInput{
				Site: []*cato_models.SiteRefInput{
					{
						By:    "ID",
						Input: "180279",
					},
				},
				Group: []*cato_models.GroupRefInput{},
			},
			Service: &cato_models.SocketLanServiceInput{
				Custom: []*cato_models.CustomServiceInput{
					{
						Port:     []scalars.Port{"443"},
						Protocol: cato_models.IPProtocolTCP,
					},
				},
				Simple: []*cato_models.SimpleServiceInput{},
			},
			Nat: &cato_models.SocketLanNatSettingsInput{
				Enabled: true,
				NatType: cato_models.SocketLanNatType("DYNAMIC_PAT"),
			},
		},
	}

	debugPrint("AddRule Request - socketLanAddRuleInput", socketLanAddRuleInput)

	result, err := catoClient.PolicySocketLanAddRule(ctx, socketLanAddRuleInput, accountId)
	if err != nil {
		fmt.Println("error adding Socket LAN rule: ", err)
		os.Exit(1)
	}

	debugPrint("AddRule Response", result)

	// Print the result
	resultJson, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println("Socket LAN rule added successfully:")
	fmt.Println(string(resultJson))

	// Access specific fields - rule.Rule contains the rule details directly
	if result.Policy.SocketLan.AddRule.Rule != nil {
		rule := result.Policy.SocketLan.AddRule.Rule
		fmt.Printf("\nRule Details:\n")
		fmt.Printf("ID: %s\n", rule.Rule.ID)
		fmt.Printf("Name: %s\n", rule.Rule.Name)
		fmt.Printf("Description: %s\n", rule.Rule.Description)
		fmt.Printf("Enabled: %t\n", rule.Rule.Enabled)
		fmt.Printf("Direction: %s\n", rule.Rule.Direction)
		fmt.Printf("Transport: %s\n", rule.Rule.Transport)
		fmt.Printf("Updated Time: %s\n", rule.Audit.UpdatedTime)

		ruleId := rule.Rule.ID

		/////////////////////////////////////////////
		// Update the Socket LAN rule             //
		/////////////////////////////////////////////

		fmt.Printf("\n======================================\n")
		fmt.Printf("Updating Socket LAN Rule\n")
		fmt.Printf("======================================\n")

		// Create update input with modified values
		updatedName := "Updated Socket LAN Rule"
		updatedDescription := "Updated Socket LAN rule description via SDK"
		enabledFlag := true
		direction := cato_models.SocketLanDirectionTo
		transport := cato_models.SocketLanTransportTypeWan

		socketLanUpdateRuleInput := cato_models.SocketLanUpdateRuleInput{
			ID: ruleId,
			Rule: &cato_models.SocketLanUpdateRuleDataInput{
				Name:        &updatedName,
				Description: &updatedDescription,
				Enabled:     &enabledFlag,
				Direction:   &direction,
				Transport:   &transport,
			},
		}

		debugPrint("UpdateRule Request - socketLanUpdateRuleInput", socketLanUpdateRuleInput)

		// Perform the update
		updateResult, err := catoClient.PolicySocketLanUpdateRule(ctx, nil, socketLanUpdateRuleInput, accountId)
		if err != nil {
			fmt.Println("error updating Socket LAN rule: ", err)
			os.Exit(1)
		}

		debugPrint("UpdateRule Response", updateResult)

		// Print the update result
		updateResultJson, _ := json.MarshalIndent(updateResult, "", "  ")
		fmt.Println("Socket LAN rule updated successfully:")
		fmt.Println(string(updateResultJson))

		// Access specific fields from update result
		if updateResult.Policy.SocketLan.UpdateRule.Rule != nil {
			updatedRule := updateResult.Policy.SocketLan.UpdateRule.Rule
			fmt.Printf("\nUpdated Rule Details:\n")
			fmt.Printf("ID: %s\n", updatedRule.Rule.ID)
			fmt.Printf("Name: %s\n", updatedRule.Rule.Name)
			fmt.Printf("Description: %s\n", updatedRule.Rule.Description)
			fmt.Printf("Enabled: %t\n", updatedRule.Rule.Enabled)
			fmt.Printf("Direction: %s\n", updatedRule.Rule.Direction)
			fmt.Printf("Transport: %s\n", updatedRule.Rule.Transport)
			fmt.Printf("Updated Time: %s\n", updatedRule.Audit.UpdatedTime)
		}

		// Check for any update errors
		if len(updateResult.Policy.SocketLan.UpdateRule.Errors) > 0 {
			fmt.Printf("\nUpdate Errors:\n")
			for _, err := range updateResult.Policy.SocketLan.UpdateRule.Errors {
				fmt.Printf("- %s (Code: %s)\n", *err.ErrorMessage, *err.ErrorCode)
			}
		} else {
			fmt.Printf("\nSocket LAN rule updated successfully!\n")
		}

		/////////////////////////////////////////////
		// Move the Socket LAN rule               //
		/////////////////////////////////////////////

		fmt.Printf("\n======================================\n")
		fmt.Printf("Moving Socket LAN Rule\n")
		fmt.Printf("======================================\n")

		// Move the rule to the first position
		movePosition := cato_models.PolicyRulePositionEnumFirstInPolicy
		policyMoveRuleInput := cato_models.PolicyMoveRuleInput{
			ID: ruleId,
			To: &cato_models.PolicyRulePositionInput{
				Position: &movePosition,
			},
		}

		debugPrint("MoveRule Request - policyMoveRuleInput", policyMoveRuleInput)

		moveResult, err := catoClient.PolicySocketLanMoveRule(ctx, policyMoveRuleInput, accountId)
		if err != nil {
			fmt.Println("error moving Socket LAN rule: ", err)
			os.Exit(1)
		}

		debugPrint("MoveRule Response", moveResult)

		// Print the move result
		moveResultJson, _ := json.MarshalIndent(moveResult, "", "  ")
		fmt.Println("Socket LAN rule moved successfully:")
		fmt.Println(string(moveResultJson))

		if moveResult.Policy.SocketLan.MoveRule.Rule != nil {
			movedRule := moveResult.Policy.SocketLan.MoveRule.Rule
			fmt.Printf("\nMoved Rule Details:\n")
			fmt.Printf("ID: %s\n", movedRule.Rule.ID)
			fmt.Printf("Name: %s\n", movedRule.Rule.Name)
			fmt.Printf("Index: %d\n", movedRule.Rule.Index)
		}

		/////////////////////////////////////////////
		// Delete the Socket LAN rule             //
		/////////////////////////////////////////////

		fmt.Printf("\n======================================\n")
		fmt.Printf("Deleting Socket LAN Rule\n")
		fmt.Printf("======================================\n")

		// Create remove input
		socketLanRemoveRuleInput := cato_models.SocketLanRemoveRuleInput{
			ID: ruleId,
		}

		debugPrint("RemoveRule Request - socketLanRemoveRuleInput", socketLanRemoveRuleInput)

		// Perform the delete operation
		deleteResult, err := catoClient.PolicySocketLanRemoveRule(ctx, nil, socketLanRemoveRuleInput, accountId)
		if err != nil {
			fmt.Println("error deleting Socket LAN rule: ", err)
			os.Exit(1)
		}

		debugPrint("RemoveRule Response", deleteResult)

		// Print the delete result
		deleteResultJson, _ := json.MarshalIndent(deleteResult, "", "  ")
		fmt.Println("Socket LAN rule deletion initiated:")
		fmt.Println(string(deleteResultJson))

		// Access specific fields from delete result
		fmt.Printf("\nDeletion Status: %s\n", deleteResult.Policy.SocketLan.RemoveRule.Status)

		// Check for any delete errors
		if len(deleteResult.Policy.SocketLan.RemoveRule.Errors) > 0 {
			fmt.Printf("\nDelete Errors:\n")
			for _, err := range deleteResult.Policy.SocketLan.RemoveRule.Errors {
				fmt.Printf("- %s (Code: %s)\n", *err.ErrorMessage, *err.ErrorCode)
			}
		} else {
			fmt.Printf("\nSocket LAN rule deletion completed successfully!\n")
			fmt.Printf("  - Rule ID %s has been marked for removal\n", ruleId)
			fmt.Printf("  - The rule will be removed from the policy after publishing\n")
		}

		/////////////////////////////////////////////
		// Publish the Socket LAN policy          //
		/////////////////////////////////////////////

		fmt.Printf("\n======================================\n")
		fmt.Printf("Publishing Socket LAN Policy\n")
		fmt.Printf("======================================\n")

		if debug {
			fmt.Println("[DEBUG] PublishPolicyRevision Request - no input parameters")
		}

		publishResult, err := catoClient.PolicySocketLanPublishPolicyRevision(ctx, nil, nil, accountId)
		if err != nil {
			fmt.Println("error publishing Socket LAN policy revision: ", err)
			os.Exit(1)
		}

		debugPrint("PublishPolicyRevision Response", publishResult)

		// Print the publish result
		publishResultJson, _ := json.MarshalIndent(publishResult, "", "  ")
		fmt.Println("Socket LAN policy revision published successfully:")
		fmt.Println(string(publishResultJson))

		// Access specific fields
		fmt.Printf("\nPublish Status: %s\n", publishResult.Policy.SocketLan.PublishPolicyRevision.Status)

		// Check for any errors
		if len(publishResult.Policy.SocketLan.PublishPolicyRevision.Errors) > 0 {
			fmt.Printf("\nPublish Errors:\n")
			errorsJson, _ := json.MarshalIndent(publishResult.Policy.SocketLan.PublishPolicyRevision.Errors, "", "  ")
			fmt.Println(string(errorsJson))
		} else {
			fmt.Printf("\nThe Socket LAN policy revision has been successfully published and is now live.\n")
			fmt.Printf("All changes made to the draft revision are now active in the production environment.\n")
			fmt.Printf("\nComplete CRUD workflow finished successfully!\n")
			fmt.Printf("  - Created Socket LAN rule: %s\n", "Example Socket LAN Rule")
			fmt.Printf("  - Updated rule name to: %s\n", updatedName)
			fmt.Printf("  - Moved rule to first position\n")
			fmt.Printf("  - Deleted rule with ID: %s\n", ruleId)
			fmt.Printf("  - Published all changes to production\n")
		}
	}
}
