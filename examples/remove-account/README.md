# Remove Account Example

This example demonstrates how to remove an existing account using the Cato Go SDK.

## Prerequisites

- Go 1.18 or later
- Valid Cato API credentials
- Partner-level account with account management permissions
- An existing account ID to remove

## Environment Variables

Set the following environment variables:

```bash
export CATO_API_KEY="your-api-key"
export CATO_ACCOUNT_ID="your-account-id"
export CATO_API_URL="https://api.catonetworks.com/api/v1/graphql2"
export CATO_ACCOUNT_TO_REMOVE="account-id-to-remove"
```

## Running the Example

```bash
# From the SDK root directory
go run examples/remove-account/main.go

# Or from the example directory
cd examples/remove-account
go run main.go
```

## What This Example Does

1. **Removes an existing account** specified by the `CATO_ACCOUNT_TO_REMOVE` environment variable
2. **Displays the result** including details of the removed account:
   - Account ID
   - Account name
   - Account type and tenancy
   - Creation timestamp and creator
   - Description

## Sample Output

```json
{
  "accountManagement": {
    "removeAccount": {
      "accountInfo": {
        "id": "12345678-1234-1234-1234-123456789012",
        "name": "Example Customer Account",
        "type": "CUSTOMER",
        "tenancy": "SINGLE_TENANT",
        "timeZone": "UTC",
        "description": "Example customer account created via SDK",
        "audit": {
          "createdBy": "admin@example.com",
          "createdTime": "2023-07-18T00:00:00Z"
        }
      }
    }
  }
}
```

## Error Handling

The example includes error handling for:
- Missing environment variables
- API authentication errors
- Account removal failures
- Invalid account IDs

## Important Notes

⚠️ **WARNING**: This operation will permanently remove the account and cannot be undone!

- Account removal requires partner-level privileges
- The account status will become "Disabled" and it will be scheduled for deletion
- All data associated with the account will be lost
- Make sure you have the correct account ID before running this example
- Consider backing up any important data before account removal

## Safety Considerations

1. **Double-check the account ID** before running the removal
2. **Verify permissions** - ensure you have the right to remove the account
3. **Data backup** - backup any critical data before removal
4. **Test environment** - test with non-production accounts first

## Use Cases

This example is useful for:
- Automated account lifecycle management
- Cleanup of test/demo accounts
- Decommissioning unused accounts
- Integration with provisioning systems
