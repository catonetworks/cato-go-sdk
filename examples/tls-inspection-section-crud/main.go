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
	// Create a new TLS inspection section //
	/////////////////////////////////////////////
	fmt.Printf("\n======================================\n")
	fmt.Printf("Creating TLS Inspection Section\n")
	fmt.Printf("======================================\n")

	position := cato_models.PolicyRulePositionEnumLastInPolicy
	policyAddSectionInput := cato_models.PolicyAddSectionInput{
		Section: &cato_models.PolicyAddSectionDataInput{
			Name: "Example TLS Inspection Section",
		},
		At: &cato_models.PolicyRulePositionInput{
			Position: &position,
		},
	}

	addResult, err := catoClient.PolicyTLSInspectAddSection(ctx, policyAddSectionInput, accountId)
	if err != nil {
		fmt.Println("error adding TLS inspection section: ", err)
		os.Exit(1)
	}

	// Print the result
	addResultJson, _ := json.MarshalIndent(addResult, "", "  ")
	fmt.Println("TLS inspection section added successfully:")
	fmt.Println(string(addResultJson))

	// Access specific fields
	if addResult.Policy.TlsInspect.AddSection.Section != nil {
		section := addResult.Policy.TlsInspect.AddSection.Section
		fmt.Printf("\nSection Details:\n")
		fmt.Printf("ID: %s\n", section.Section.ID)
		fmt.Printf("Name: %s\n", section.Section.Name)
		fmt.Printf("Updated by: %s\n", section.Audit.UpdatedBy)
		fmt.Printf("Updated time: %s\n", section.Audit.UpdatedTime)
	}

	// Check for any errors
	if len(addResult.Policy.TlsInspect.AddSection.Errors) > 0 {
		fmt.Printf("\nAdd Section Errors:\n")
		for _, err := range addResult.Policy.TlsInspect.AddSection.Errors {
			fmt.Printf("- %s (Code: %s)\n", *err.ErrorMessage, *err.ErrorCode)
		}
		os.Exit(1)
	}

	//////////////////////////////////////////
	// Read the TLS inspection policy    //
	//////////////////////////////////////////

	if addResult.Policy.TlsInspect.AddSection.Section != nil {
		section := addResult.Policy.TlsInspect.AddSection.Section
		sectionId := section.Section.ID
		sectionName := section.Section.Name

		fmt.Printf("\n======================================\n")
		fmt.Printf("Reading TLS Inspection Policy\n")
		fmt.Printf("======================================\n")

		// Query the TLS inspection policy to get the current state of all sections
		policyResult, err := catoClient.Tlsinspectpolicy(ctx, accountId)
		if err != nil {
			fmt.Println("error reading TLS inspection policy: ", err)
			os.Exit(1)
		}

		// Display the section details
		fmt.Printf("Section ID: %s\n", sectionId)
		fmt.Printf("Section Name: %s\n", sectionName)

		// Display sections if available
		if len(policyResult.Policy.TlsInspect.Policy.Sections) > 0 {
			fmt.Printf("Total sections in policy: %d\n", len(policyResult.Policy.TlsInspect.Policy.Sections))
			// Look for our specific section
			for _, sec := range policyResult.Policy.TlsInspect.Policy.Sections {
				if sec.Section.ID == sectionId {
					fmt.Printf("Found our section: %s\n", sec.Section.Name)
					break
				}
			}
		} else {
			fmt.Printf("No sections found in policy\n")
		}

		/////////////////////////////////////////////
		// Update the TLS inspection section   //
		/////////////////////////////////////////////

		fmt.Printf("\n======================================\n")
		fmt.Printf("Updating TLS Inspection Section\n")
		fmt.Printf("======================================\n")

		// Create update input with modified values
		updatedName := "Updated Example TLS Inspection Section"

		policyUpdateSectionInput := cato_models.PolicyUpdateSectionInput{
			ID: sectionId,
			Section: &cato_models.PolicyUpdateSectionDataInput{
				Name: &updatedName,
			},
		}

		// Perform the update
		updateResult, err := catoClient.PolicyTLSInspectUpdateSection(ctx, policyUpdateSectionInput, accountId)
		if err != nil {
			fmt.Println("error updating TLS inspection section: ", err)
			os.Exit(1)
		}

		// Print the update result
		updateResultJson, _ := json.MarshalIndent(updateResult, "", "  ")
		fmt.Println("TLS inspection section updated successfully:")
		fmt.Println(string(updateResultJson))

		// Access specific fields from update result
		if updateResult.Policy.TlsInspect.UpdateSection.Section != nil {
			updatedSection := updateResult.Policy.TlsInspect.UpdateSection.Section
			fmt.Printf("\nUpdated Section Details:\n")
			fmt.Printf("ID: %s\n", updatedSection.Section.ID)
			fmt.Printf("Name: %s\n", updatedSection.Section.Name)
			fmt.Printf("Updated by: %s\n", updatedSection.Audit.UpdatedBy)
			fmt.Printf("Updated time: %s\n", updatedSection.Audit.UpdatedTime)
		}

		// Check for any update errors
		if len(updateResult.Policy.TlsInspect.UpdateSection.Errors) > 0 {
			fmt.Printf("\nUpdate Errors:\n")
			for _, err := range updateResult.Policy.TlsInspect.UpdateSection.Errors {
				fmt.Printf("- %s (Code: %s)\n", *err.ErrorMessage, *err.ErrorCode)
			}
		} else {
			fmt.Printf("\n✓ TLS inspection section updated successfully!\n")
			fmt.Printf("  - Name changed from: %s\n", sectionName)
			fmt.Printf("  - Name changed to: %s\n", updatedName)
		}

		/////////////////////////////////////////////
		// Delete the TLS inspection section   //
		/////////////////////////////////////////////

		fmt.Printf("\n======================================\n")
		fmt.Printf("Deleting TLS Inspection Section\n")
		fmt.Printf("======================================\n")

		// Create remove input
		policyRemoveSectionInput := cato_models.PolicyRemoveSectionInput{
			ID: sectionId,
		}

		// Perform the delete operation
		deleteResult, err := catoClient.PolicyTLSInspectRemoveSection(ctx, policyRemoveSectionInput, accountId)
		if err != nil {
			fmt.Println("error deleting TLS inspection section: ", err)
			os.Exit(1)
		}

		// Print the delete result
		deleteResultJson, _ := json.MarshalIndent(deleteResult, "", "  ")
		fmt.Println("TLS inspection section deletion initiated:")
		fmt.Println(string(deleteResultJson))

		// Access specific fields from delete result
		fmt.Printf("\nDeletion Status: %s\n", deleteResult.Policy.TlsInspect.RemoveSection.Status)

		// Check for any delete errors
		if len(deleteResult.Policy.TlsInspect.RemoveSection.Errors) > 0 {
			fmt.Printf("\nDelete Errors:\n")
			for _, err := range deleteResult.Policy.TlsInspect.RemoveSection.Errors {
				fmt.Printf("- %s (Code: %s)\n", *err.ErrorMessage, *err.ErrorCode)
			}
		} else {
			fmt.Printf("\n✓ TLS inspection section deletion completed successfully!\n")
			fmt.Printf("  - Section ID %s has been marked for removal\n", sectionId)
			fmt.Printf("  - The section will be removed from the policy after publishing\n")
		}

		/////////////////////////////////////////////
		// Publish the TLS inspection policy   //
		/////////////////////////////////////////////

		fmt.Printf("\n======================================\n")
		fmt.Printf("Publishing TLS Inspection Policy\n")
		fmt.Printf("======================================\n")

		publishResult, err := catoClient.PolicyTlsInspectPublishPolicyRevision(ctx, accountId)
		if err != nil {
			fmt.Println("error publishing TLS inspection policy revision: ", err)
			os.Exit(1)
		}

		// Print the publish result
		publishResultJson, _ := json.MarshalIndent(publishResult, "", "  ")
		fmt.Println("\nTLS inspection policy revision published successfully:")
		fmt.Println(string(publishResultJson))

		// Access specific fields
		fmt.Printf("\nPublish Status: %s\n", publishResult.Policy.TlsInspect.PublishPolicyRevision.Status)

		// Check for any errors
		if len(publishResult.Policy.TlsInspect.PublishPolicyRevision.Errors) > 0 {
			fmt.Printf("\nPublish Errors:\n")
			errorsJson, _ := json.MarshalIndent(publishResult.Policy.TlsInspect.PublishPolicyRevision.Errors, "", "  ")
			fmt.Println(string(errorsJson))
		} else {
			fmt.Printf("\nThe TLS inspection policy revision has been successfully published and is now live.\n")
			fmt.Printf("All changes made to the draft revision are now active in the production environment.\n")
			fmt.Printf("\n🎉 Complete CRUD workflow finished successfully!\n")
			fmt.Printf("  ✓ Created TLS inspection section: %s\n", "Example TLS Inspection Section")
			fmt.Printf("  ✓ Read section from policy\n")
			fmt.Printf("  ✓ Updated section name to: %s\n", updatedName)
			fmt.Printf("  ✓ Deleted section with ID: %s\n", sectionId)
			fmt.Printf("  ✓ Published all changes to production\n")
		}

	}
}
