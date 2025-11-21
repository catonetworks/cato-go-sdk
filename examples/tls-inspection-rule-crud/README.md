# WAN Network Add Rule Example

This example demonstrates how to add a new WAN network policy rule using the Cato Go SDK.

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
go run examples/wan-network-add-rule/main.go

# Or from the example directory
cd examples/wan-network-add-rule
go run main.go
```

## What This Example Does

1. **Creates a new WAN network rule** with the following properties:
   - Name: "Example WAN Rule"
   - Type: Route
   - Status: Enabled
   - Section: Reference to an existing section

2. **Displays the result** including:
   - Rule ID
   - Rule name and description
   - Rule type and status
   - Section information
   - Audit information (who created/updated and when)

## Sample Output

```json
{
  "policy": {
    "wanNetwork": {
      "addRule": {
        "wanNetworkRulePayloadRule": {
          "audit": {
            "updatedBy": "admin@example.com",
            "updatedTime": "2023-07-18T00:00:00Z"
          },
          "rule": {
            "id": "12345678",
            "name": "Example WAN Rule",
            "description": "Example WAN network rule created via SDK",
            "enabled": true,
            "ruleType": "ROUTE"
          }
        }
      }
    }
  }
}
```

## Customization

You can modify the `wanNetworkAddRuleInput` struct to customize:
- Rule name and description
- Rule type (Route, Bypass, etc.)
- Source and destination definitions
- Application specifications
- Configuration settings
- Priority and positioning

## Error Handling

The example includes error handling for:
- Missing environment variables
- API authentication errors
- Rule creation failures
- Invalid input parameters

## Notes

- The section ID must reference an existing section in your WAN network policy
- Rule names should be unique within the policy
- Some rule types may require additional configuration parameters
- The rule will be created in the current policy revision
