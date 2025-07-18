# WAN Network Section Management Example

This example demonstrates how to add, read, update, move, publish, and delete WAN network policy sections using the Cato Go SDK.

## Prerequisites

- Go 1.18 or later
- Valid Cato API credentials
- Account with WAN network policy management permissions

## Environment Variables

Set the following environment variables:

```bash
export CATO_API_KEY="your-api-key"
export CATO_ACCOUNT_ID="your-account-id"
export CATO_API_URL="https://api.catonetworks.com/api/v1/graphql2"
```

## Running the Example

```bash
# From the SDK root directory
go run examples/wan-network-section/main.go

# Or from the example directory
cd examples/wan-network-section
go run main.go
```

## What This Example Does

1. **Creates a new WAN network section** with the following properties:
   - Name: "Example WAN Section"
   - Position: Default (at the end)

2. **Reads the section** by:
   - Displaying the section information from the create response
   - Showing section ID, name, and audit information
   - Demonstrating how to work with API response data

3. **Deletes the section** by:
   - Removing the section from the policy
   - Demonstrating safe cleanup and resource management
   - Showing proper error handling for deletion operations

4. **Alternative workflow** (if deletion is commented out):
   - **Updates the section** by changing the name to "Updated WAN Section Name"
   - **Publishes the policy revision** to make changes live in the production environment
   - **Deletes the updated section** and publishes the deletion (following Internet Firewall pattern)
   - **Moves the section** to a new position within the policy (if section deletion is also commented out)

5. **Displays comprehensive results** including:
   - Section ID and name after each operation
   - Audit information (who created/updated/deleted and when)
   - Operation status and any error messages
   - Full API response for each step

## Sample Output

```json
{
  "policy": {
    "wanNetwork": {
      "addSection": {
        "policySectionPayloadSection": {
          "audit": {
            "updatedBy": "admin@example.com",
            "updatedTime": "2023-07-18T00:00:00Z"
          },
          "section": {
            "id": "12345678",
            "name": "Example WAN Section"
          }
        },
        "status": "SUCCESS"
      }
    }
  }
}
```

## Operations Covered

### 1. Add Section
- Creates a new WAN network section
- Uses `PolicyWanNetworkAddSection` API
- Demonstrates proper input structure and positioning

### 2. Read Section
- Displays section information from API responses
- Shows how to extract and work with section data
- Demonstrates data access patterns for WAN network sections
- Note: Uses create response data (WAN network query API not available in current SDK)

### 3. Delete Section
- Removes an existing section from the policy
- Uses `PolicyWanNetworkRemoveSection` API
- Shows proper cleanup and resource management
- Includes comprehensive error handling

### 4. Update Section
- Updates existing section properties (name)
- Uses `PolicyWanNetworkUpdateSection` API
- Shows how to modify section metadata

### 5. Publish Policy Revision
- Publishes the current policy revision to make changes live
- Uses `PolicyWanNetworkPublishPolicyRevision` API
- Demonstrates how to activate policy changes in production
- Essential step for making modifications active in the WAN network
- Includes comprehensive status checking and error handling
- Similar to the Internet Firewall publish operation pattern

### 6. Delete Section (After Publish)
- Demonstrates deletion of an updated section after publishing changes
- Uses `PolicyWanNetworkRemoveSection` API followed by `PolicyWanNetworkPublishPolicyRevision` API
- Follows the same pattern as Internet Firewall section deletion with publish
- Shows how to properly clean up published resources
- Includes comprehensive error handling for both delete and publish operations
- Essential for complete lifecycle management in production environments

### 7. Move Section
- Repositions section within the policy
- Uses `PolicyWanNetworkMoveSection` API
- Demonstrates position management

## Customization

You can modify the inputs to customize:
- Section name and description
- Position within the policy (using the `at` parameter)
- Section placement relative to other sections
- Policy revision targeting

## Error Handling

The example includes error handling for:
- Missing environment variables
- API authentication errors
- Section creation/update/move/delete failures
- Policy revision publish failures
- Invalid input parameters
- Policy mutation errors
- Concurrent policy modification issues
- Section deletion validation (e.g., sections with rules cannot be deleted)

## Notes

- Section names should be unique within the policy
- Sections are used to organize rules logically
- Operations are performed in the current policy revision
- Position can be specified using `LAST_IN_POLICY`, `FIRST_IN_POLICY`, `BEFORE_SECTION`, or `AFTER_SECTION`
- All operations return comprehensive audit information
- The example demonstrates a complete workflow from creation to deletion
- **Section deletion requirements:**
  - Sections must be empty (no rules) before they can be deleted
  - System sections cannot be deleted
  - Sections locked by another revision cannot be deleted
- The example currently deletes the section immediately after creation (before adding any rules)
- To test update/publish/move operations, comment out the delete section and uncomment the update/publish/move code
- **Policy Publishing:**
  - Publishing is required to make draft policy changes live in production
  - The publish operation follows the same pattern as Internet Firewall policy publishing
  - All policy modifications remain in draft state until published
  - Publishing activates changes across the entire WAN network infrastructure
- **Delete-After-Publish Pattern:**
  - The example demonstrates the Internet Firewall pattern of deleting a section after publishing
  - This pattern shows: Update → Publish → Delete → Publish (to make deletion live)
  - Each publish operation activates the changes in the production environment
  - This is essential for proper lifecycle management in production systems
  - The pattern ensures that both modifications and deletions are properly published
