# Add Account Example

This example demonstrates how to create a new account using the Cato Go SDK.

## Prerequisites

- Go 1.18 or later
- Valid Cato API credentials
- Partner-level account with account management permissions

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
go run examples/add-account/main.go

# Or from the example directory
cd examples/add-account
go run main.go
```

## What This Example Does

1. **Creates a new account** with the following properties:
   - Name: "Example Customer Account"
   - Type: Customer
   - Tenancy: Single Tenant
   - Description: "Example customer account created via SDK"
   - Timezone: UTC

2. **Displays the result** including:
   - Account ID
   - Account name
   - Account type and tenancy
   - Creation timestamp and creator
   - Description

## Sample Output

```json
{
  "accountManagement": {
    "addAccount": {
      "id": "12345678",
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
```

## Customization

You can modify the `addAccountInput` struct to customize:
- Account name
- Description
- Account type (Customer or Partner)
- Tenancy (Single or Multi-tenant)
- Timezone

## Error Handling

The example includes error handling for:
- Missing environment variables
- API authentication errors
- Account creation failures
- Invalid input parameters

## Notes

- Account creation requires partner-level privileges
- The account ID is automatically generated
- The timezone should be a valid timezone string (e.g., "UTC", "America/New_York")
- Account names should be unique within your organization
