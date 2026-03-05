package main

import (
	"context"
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

	// Test scenario: Create a socket site and then update with invalid DHCP range
	fmt.Println("=== Starting Socket Site CRUD Test ===")

	// Step 1: Create the socket site
	fmt.Println("\nStep 1: Creating socket site...")
	siteID, nativeRangeID, err := createSocketSite(ctx, catoClient, accountId)
	if err != nil {
		fmt.Printf("Error creating socket site: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Socket site created successfully. Site ID: %s, Native Range ID: %s\n", siteID, nativeRangeID)

	// Step 2: Update the site with new network range and DHCP settings
	fmt.Println("\nStep 2: Updating site with new network range and DHCP settings...")
	err = updateSiteNetworkRange(ctx, catoClient, accountId, siteID, nativeRangeID)
	if err != nil {
		fmt.Printf("Error updating site: %v\n", err)
		// This is expected to fail due to invalid DHCP range
		fmt.Println("\n⚠ Update failed as expected - DHCP range not within native network range")
	} else {
		fmt.Println("✓ Site updated successfully")
	}

	// Step 3: Clean up - remove the site
	fmt.Println("\nStep 3: Cleaning up - removing site...")
	err = deleteSite(ctx, catoClient, accountId, siteID)
	if err != nil {
		fmt.Printf("Error deleting site: %v\n", err)
	} else {
		fmt.Println("✓ Site deleted successfully")
	}

	fmt.Println("\n=== Test Complete ===")
}

// createSocketSite creates a socket site with initial configuration
func createSocketSite(ctx context.Context, client *cato.Client, accountId string) (string, string, error) {
	// Initial configuration matching the Terraform resource
	city := "San Diego"
	address := "555 That Way"
	stateCode := "US-CA"

	input := cato_models.AddSocketSiteInput{
		Name:           "SOCKET_X1500_default",
		SiteType:       cato_models.SiteTypeDatacenter,
		ConnectionType: cato_models.SiteConnectionTypeEnumSocketX1500,
		NativeNetworkRange: "192.160.151.0/24",
		SiteLocation: &cato_models.AddSiteLocationInput{
			City:        &city,
			CountryCode: "US",
			StateCode:   &stateCode,
			Timezone:    "America/Los_Angeles",
			Address:     &address,
		},
	}

	fmt.Printf("Creating site with config: %+v\n", input)
	resp, err := client.SiteAddSocketSite(ctx, input, accountId)
	if err != nil {
		return "", "", fmt.Errorf("failed to create socket site: %w", err)
	}

	if resp == nil || resp.Site.AddSocketSite == nil {
		return "", "", fmt.Errorf("nil response from AddSocketSite API")
	}

	siteID := resp.Site.AddSocketSite.SiteID

	// Get the native range ID
	nativeRangeID, err := getNativeRangeID(ctx, client, accountId, siteID)
	if err != nil {
		return siteID, "", fmt.Errorf("failed to get native range ID: %w", err)
	}

	// Update network range with initial settings
	localIP := "192.160.151.1"
	subnet := "192.160.151.0/24"
	dhcpType := cato_models.DhcpTypeAccountDefault

	updateNetworkRangeInput := cato_models.UpdateNetworkRangeInput{
		LocalIP: &localIP,
		Subnet:  &subnet,
		DhcpSettings: &cato_models.NetworkDhcpSettingsInput{
			DhcpType: dhcpType,
		},
	}

	fmt.Printf("Updating network range with: %+v\n", updateNetworkRangeInput)
	_, err = client.SiteUpdateNetworkRange(ctx, nativeRangeID, updateNetworkRangeInput, accountId)
	if err != nil {
		return siteID, nativeRangeID, fmt.Errorf("failed to update network range: %w", err)
	}

	// Update socket interface with LAG configuration
	interfaceName := "LAN1"
	lagMinLinks := int64(2)

	updateInterfaceInput := cato_models.UpdateSocketInterfaceInput{
		Name:     &interfaceName,
		DestType: cato_models.SocketInterfaceDestTypeLanLagMaster,
		Lag: &cato_models.SocketInterfaceLagInput{
			MinLinks: lagMinLinks,
		},
		Lan: &cato_models.SocketInterfaceLanInput{
			LocalIP: localIP,
			Subnet:  subnet,
		},
	}

	fmt.Printf("Updating socket interface with: %+v\n", updateInterfaceInput)
	_, err = client.SiteUpdateSocketInterface(ctx, siteID, cato_models.SocketInterfaceIDEnumLan1, updateInterfaceInput, accountId)
	if err != nil {
		return siteID, nativeRangeID, fmt.Errorf("failed to update socket interface: %w", err)
	}

	return siteID, nativeRangeID, nil
}

// updateSiteNetworkRange updates the site with new configuration that has invalid DHCP range
func updateSiteNetworkRange(ctx context.Context, client *cato.Client, accountId, siteID, nativeRangeID string) error {
	// Update configuration - intentionally using DHCP range outside native network range
	localIP := "192.159.151.1"
	subnet := "192.159.151.0/24"
	dhcpType := cato_models.DhcpTypeDhcpRange
	// Note: This DHCP range (192.159.153.x) is NOT within the native range (192.159.151.0/24)
	// This should trigger an API error
	ipRange := "192.159.153.100-192.159.153.150"
	dhcpMicrosegmentation := false

	updateNetworkRangeInput := cato_models.UpdateNetworkRangeInput{
		LocalIP: &localIP,
		Subnet:  &subnet,
		DhcpSettings: &cato_models.NetworkDhcpSettingsInput{
			DhcpType:              dhcpType,
			IPRange:               &ipRange,
			DhcpMicrosegmentation: &dhcpMicrosegmentation,
		},
	}

	fmt.Printf("Updating network range with: %+v\n", updateNetworkRangeInput)
	resp, err := client.SiteUpdateNetworkRange(ctx, nativeRangeID, updateNetworkRangeInput, accountId)
	
	// Check response even if no error
	if resp != nil {
		fmt.Printf("API Response: %+v\n", resp)
		if resp.Site.UpdateNetworkRange == nil {
			return fmt.Errorf("API returned null updateNetworkRange - this indicates a GraphQL error (likely: DHCP Range should be included in the native range)")
		}
	}
	
	if err != nil {
		return fmt.Errorf("failed to update network range: %w", err)
	}

	// Update socket interface to LAN (from LAN_LAG_MASTER)
	interfaceName := "LAN1"
	updateInterfaceInput := cato_models.UpdateSocketInterfaceInput{
		Name:     &interfaceName,
		DestType: cato_models.SocketInterfaceDestTypeLan,
		Lan: &cato_models.SocketInterfaceLanInput{
			LocalIP: localIP,
			Subnet:  subnet,
		},
	}

	fmt.Printf("Updating socket interface with: %+v\n", updateInterfaceInput)
	_, err = client.SiteUpdateSocketInterface(ctx, siteID, cato_models.SocketInterfaceIDEnumLan1, updateInterfaceInput, accountId)
	if err != nil {
		return fmt.Errorf("failed to update socket interface: %w", err)
	}

	return nil
}

// getNativeRangeID retrieves the native network range ID for a site
func getNativeRangeID(ctx context.Context, client *cato.Client, accountId, siteID string) (string, error) {
	siteEntity := &cato_models.EntityInput{
		Type: "site",
		ID:   siteID,
	}
	zeroInt64 := int64(0)

	resp, err := client.EntityLookup(ctx, accountId, cato_models.EntityTypeSiteRange, &zeroInt64, nil, siteEntity, nil, nil, nil, nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to lookup site ranges: %w", err)
	}

	if resp == nil || len(resp.EntityLookup.Items) == 0 {
		return "", fmt.Errorf("no ranges found for site")
	}

	// Find the native range (first range is typically the native range)
	for _, item := range resp.EntityLookup.Items {
		if item != nil && item.Entity.ID != "" {
			return item.Entity.ID, nil
		}
	}

	return "", fmt.Errorf("could not find native range ID")
}

// deleteSite removes the test site
func deleteSite(ctx context.Context, client *cato.Client, accountId, siteID string) error {
	_, err := client.SiteRemoveSite(ctx, siteID, accountId)
	if err != nil {
		return fmt.Errorf("failed to delete site: %w", err)
	}
	return nil
}
