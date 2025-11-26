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

	/////////////////
	// Create Group //
	/////////////////
	fmt.Println("Creating group...")
	description := "Example group created via SDK"

	// Example: Create a group with no members initially
	createGroupInput := cato_models.CreateGroupInput{
		Name:        "Example_SDK_Group",
		Description: &description,
		Members: []*cato_models.GroupMemberRefTypedInput{
			{
				By:    cato_models.ObjectRefByID,
				Input: "1474041",
				Type:  cato_models.GroupMemberRefTypeFloatingSubnet,
			},
			{
				By:    cato_models.ObjectRefByName,
				Input: "global_ip_range",
				Type:  cato_models.GroupMemberRefTypeGlobalIPRange,
			},
			{
				By:    cato_models.ObjectRefByID,
				Input: "2528580",
				Type:  cato_models.GroupMemberRefTypeHost,
			},
			{
				By:    cato_models.ObjectRefByID,
				Input: "124986",
				Type:  cato_models.GroupMemberRefTypeNetworkInterface,
			},
			{
				By:    cato_models.ObjectRefByName,
				Input: "ipsec-dev-site",
				Type:  cato_models.GroupMemberRefTypeSite,
			},
			{
				By:    cato_models.ObjectRefByID,
				Input: "UzU4OTI1Mw==",
				Type:  cato_models.GroupMemberRefTypeSiteNetworkSubnet,
			},
		},
	}

	createResult, err := catoClient.GroupsCreateGroup(ctx, createGroupInput, accountId)
	if err != nil {
		fmt.Println("error creating group: ", err)
		os.Exit(1)
	}

	// Print the result
	resultJson, _ := json.MarshalIndent(createResult, "", "  ")
	fmt.Println("Group created successfully:")
	fmt.Println(string(resultJson))

	if createResult.Groups.CreateGroup.Group.ID == "" {
		fmt.Println("Error: Group ID is empty")
		os.Exit(1)
	}

	groupID := createResult.Groups.CreateGroup.Group.ID
	fmt.Printf("Group ID: %s\n", groupID)

	/////////////////
	// Read Group  //
	/////////////////
	fmt.Println("\nReading group...")
	groupRef := cato_models.GroupRefInput{
		By:    cato_models.ObjectRefByID,
		Input: groupID,
	}

	groupMembersListInput := cato_models.GroupMembersListInput{
		Filter: []*cato_models.GroupMembersListFilterInput{},
		Paging: &cato_models.PagingInput{
			Limit: 1000,
			From:  0,
		},
		Sort: &cato_models.GroupMembersListSortInput{},
	}

	readResult, err := catoClient.GroupsMembers(ctx, groupRef, groupMembersListInput, accountId)
	if err != nil {
		fmt.Println("error reading group: ", err)
	} else {
		readJson, _ := json.MarshalIndent(readResult, "", "  ")
		fmt.Println("Group read successfully:")
		fmt.Println(string(readJson))
	}

	/////////////////
	// Update Group //
	/////////////////
	fmt.Println("\nUpdating group...")

	newDescription := "Updated group description via SDK"
	newName := "Updated_Example_SDK_Group"

	updateGroupInput := cato_models.UpdateGroupInput{
		Group:       &cato_models.GroupRefInput{By: cato_models.ObjectRefByID, Input: groupID},
		Description: &newDescription,
		Name:        &newName,
	}

	updateResult, err := catoClient.GroupsUpdateGroup(ctx, updateGroupInput, accountId)
	if err != nil {
		fmt.Println("error updating group: ", err)
		os.Exit(1)
	}

	// Print the result
	updateJson, _ := json.MarshalIndent(updateResult, "", "  ")
	fmt.Println("Group updated successfully:")
	fmt.Println(string(updateJson))

	/////////////////
	// Delete Group //
	/////////////////
	fmt.Println("\nDeleting group...")

	deleteGroupInput := cato_models.GroupRefInput{
		By:    cato_models.ObjectRefByID,
		Input: groupID,
	}

	deleteResult, err := catoClient.GroupsDeleteGroup(ctx, deleteGroupInput, accountId)
	if err != nil {
		fmt.Println("error deleting group: ", err)
		os.Exit(1)
	}

	// Print the result
	deleteJson, _ := json.MarshalIndent(deleteResult, "", "  ")
	fmt.Println("Group deleted successfully:")
	fmt.Println(string(deleteJson))
}
