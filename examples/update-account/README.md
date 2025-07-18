# Update Account Example

This example demonstrates how to update an existing account using the Cato Go SDK.

## Prerequisites

- Go 1.18 or later
- Valid Cato API credentials
- Account with appropriate permissions to update account information

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
go run examples/update-account/main.go

# Or from the example directory
cd examples/update-account
go run main.go
```

## What This Example Does

1. **Updates the current account** (specified by `CATO_ACCOUNT_ID`) with a new description
2. **Displays the result** including updated account details:
   - Account ID
   - Account name
   - Account type and tenancy
   - Updated description
   - Creation timestamp and creator

## Sample Output

```json
{
  "accountManagement": {
    "updateAccount": {
      "id": "12345678-1234-1234-1234-123456789012",
      "name": "My Account",
      "type": "CUSTOMER",
      "tenancy": "SINGLE_TENANT",
      "timeZone": "UTC",
      "description": "Updated account description via SDK",
      "audit": {
        "createdBy": "admin@example.com",
        "createdTime": "2023-07-18T00:00:00Z"
      }
    }
  }
}
```

## Customization

You can modify the `updateAccountInput` struct to update different fields. Currently supported updates:

- **Description**: Update the account description

### Example Modifications

```go
// Update only description
updateAccountInput := cato_models.UpdateAccountInput{
    Description: &newDescription,
}

// Clear description (set to null)
updateAccountInput := cato_models.UpdateAccountInput{
    Description: nil,
}
```

## Error Handling

The example includes error handling for:
- Missing environment variables
- API authentication errors
- Account update failures
- Invalid input parameters

## Use Cases

This example is useful for:
- Updating account metadata
- Maintaining account descriptions
- Automated account management workflows
- Integration with external systems

## Notes

- Account updates require appropriate permissions
- Only the description field can be updated currently
- The account ID cannot be changed
- Other account properties (name, type, tenancy) are typically immutable
- Updates are applied immediately and are reflected in the response

## API Limitations

Based on the current API schema, the `UpdateAccountInput` only supports:
- `description` (String, optional)

Other account properties like name, type, tenancy, and timezone are not updatable through this endpoint and would need to be modified through other means or during account creation.
