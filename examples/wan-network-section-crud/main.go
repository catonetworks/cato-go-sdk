package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

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
	// Create a new WAN network section //
	//////////////////////////////////////

	// Create unique section name with timestamp to avoid conflicts
	timestamp := time.Now().Format("20060102-150405")
	uniqueSectionName := fmt.Sprintf("Example WAN Section %s", timestamp)

	policyAddSectionInput := cato_models.PolicyAddSectionInput{
		At: &cato_models.PolicySectionPositionInput{
			Position: "LAST_IN_POLICY",
		},
		Section: &cato_models.PolicyAddSectionInfoInput{
			Name: uniqueSectionName,
		},
		// Add at parameter if you want to specify position
	}

	result, err := catoClient.PolicyWanNetworkAddSection(ctx, policyAddSectionInput, accountId, nil)
	if err != nil {
		fmt.Println("error adding WAN network section: ", err)
		os.Exit(1)
	}

	// Print the result
	resultJson, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println("WAN network section added successfully:")
	fmt.Println(string(resultJson))

	// Access specific fields
	if result.Policy.WanNetwork.AddSection.Section != nil {
		section := result.Policy.WanNetwork.AddSection.Section
		fmt.Printf("\nSection Details:\n")
		fmt.Printf("ID: %s\n", section.Section.ID)
		fmt.Printf("Name: %s\n", section.Section.Name)
		fmt.Printf("Updated by: %s\n", section.Audit.UpdatedBy)
		fmt.Printf("Updated time: %s\n", section.Audit.UpdatedTime)
	}

	// Check for any errors
	if len(result.Policy.WanNetwork.AddSection.Errors) > 0 {
		fmt.Printf("\nErrors:\n")
		for _, err := range result.Policy.WanNetwork.AddSection.Errors {
			fmt.Printf("- %s (Code: %s)\n", *err.ErrorMessage, *err.ErrorCode)
		}
	}

	//////////////////////////////////////
	// Read the WAN network section     //
	//////////////////////////////////////

	// Read and verify the section we just created
	if result.Policy.WanNetwork.AddSection.Section != nil {
		sectionId := result.Policy.WanNetwork.AddSection.Section.Section.ID
		sectionName := result.Policy.WanNetwork.AddSection.Section.Section.Name

		fmt.Printf("\n======================================\n")
		fmt.Printf("Reading WAN Network Section\n")
		fmt.Printf("======================================\n")
		// Query the WAN network policy to get the current state of all sections
		policyResult, err := catoClient.WanNetworkPolicy(ctx, accountId, nil)
		if err != nil {
			fmt.Println("error reading WAN network policy: ", err)
			os.Exit(1)
		}

		// Display the section details that we have from the creation response
		fmt.Printf("Section ID: %s\n", sectionId)
		fmt.Printf("Section Name: %s\n", sectionName)

		// Display properties if available
		if len(policyResult.Policy.WanNetwork.Policy.Sections) > 0 {
			fmt.Printf("Total sections in policy: %d\n", len(policyResult.Policy.WanNetwork.Policy.Sections))
			// Look for our specific section
			for _, section := range policyResult.Policy.WanNetwork.Policy.Sections {
				if section.Section.ID == sectionId {
					fmt.Printf("Found our section: %s\n", section.Section.Name)
					break
				}
			}
		} else {
			fmt.Printf("No sections found in policy\n")
		}

		// Note: This shows the section data from the create response.
		// Later in this example, we'll demonstrate how to query the full WAN network policy
		// to get the current state of all sections after publishing.

		fmt.Printf("\nSection read operation completed successfully!\n")
		fmt.Printf("Note: This example shows the section data from the create response.\n")
		fmt.Printf("After publishing, we'll query the full WAN network policy to verify\n")
		fmt.Printf("the section in the live environment.\n")

		//////////////////////////////////////
		// Update a new WAN network section //
		//////////////////////////////////////

		// Create update input
		updatedName := "Updated WAN Section Name"
		policyUpdateSectionInput := cato_models.PolicyUpdateSectionInput{
			ID: sectionId,
			Section: &cato_models.PolicyUpdateSectionInfoInput{
				Name: &updatedName,
			},
		}

		// Update the section
		updateResult, err := catoClient.PolicyWanNetworkUpdateSection(ctx, policyUpdateSectionInput, accountId, nil)
		if err != nil {
			fmt.Println("error updating WAN network section: ", err)
			os.Exit(1)
		}

		if updateResult.Policy.WanNetwork.UpdateSection.Section != nil {
			section := updateResult.Policy.WanNetwork.UpdateSection.Section
			fmt.Printf("Section ID: %s\n", section.Section.ID)
			fmt.Printf("Section Name: %s\n", section.Section.Name)
			if section.Section.ID == sectionId {
				fmt.Printf("Found our section: %s\n", section.Section.Name)
			}
		} else {
			fmt.Printf("No section found in policy\n")
		}

		// Check for any errors
		if len(updateResult.Policy.WanNetwork.UpdateSection.Errors) > 0 {
			fmt.Printf("\nUpdate Errors:\n")
			for _, err := range updateResult.Policy.WanNetwork.UpdateSection.Errors {
				fmt.Printf("- %s (Code: %s)\n", *err.ErrorMessage, *err.ErrorCode)
			}
		}

		//////////////////////////////////////
		// Delete the WAN network section   //
		//////////////////////////////////////

		// Create remove input
		policyRemoveSectionInput := cato_models.PolicyRemoveSectionInput{
			ID: sectionId,
		}

		// Remove the section
		removeResult, err := catoClient.PolicyWanNetworkRemoveSection(ctx, policyRemoveSectionInput, accountId, nil)
		if err != nil {
			fmt.Println("error removing WAN network section: ", err)
			os.Exit(1)
		}

		// Print the remove result
		removeResultJson, _ := json.MarshalIndent(removeResult, "", "  ")
		fmt.Println("\nWAN network section removed successfully:")
		fmt.Println(string(removeResultJson))

		// Check for any errors
		if len(removeResult.Policy.WanNetwork.RemoveSection.Errors) > 0 {
			fmt.Printf("\nRemove Errors:\n")
			for _, err := range removeResult.Policy.WanNetwork.RemoveSection.Errors {
				fmt.Printf("- %s (Code: %s)\n", *err.ErrorMessage, *err.ErrorCode)
			}
		}

		// Show confirmation of deletion
		if removeResult.Policy.WanNetwork.RemoveSection.Section != nil {
			section := removeResult.Policy.WanNetwork.RemoveSection.Section
			fmt.Printf("\nDeleted Section Details:\n")
			fmt.Printf("ID: %s\n", section.Section.ID)
			fmt.Printf("Name: %s\n", section.Section.Name)
			fmt.Printf("Deleted by: %s\n", section.Audit.UpdatedBy)
			fmt.Printf("Deleted time: %s\n", section.Audit.UpdatedTime)
		}

		fmt.Println("\nSection deletion completed successfully!")

		//////////////////////////////////////
		// Publish the WAN network policy   //
		//////////////////////////////////////

		// Publish the policy revision to make changes live
		publishResult, err := catoClient.PolicyWanNetworkPublishPolicyRevision(ctx, accountId, nil, nil)
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
		}

		return // Exit after successful deletion and publishing

	}

}
