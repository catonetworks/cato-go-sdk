package cato_go_sdk_test

import (
	"context"
	"os"
	"testing"

	cato "github.com/catonetworks/cato-go-sdk"
)

// TestSubPolicyQueryFieldsLive validates that the regenerated client and the
// updated handwritten policy documents expose the sub-policy hierarchy fields
// (subPolicies, rule.ruleType, rule.subPolicy) for both firewalls. It is skipped
// unless CATO_TOKEN/CATO_BASEURL/CATO_ACCOUNT_ID are set, so it never runs in
// offline CI. Mutation lifecycle behaviour is covered by the provider
// acceptance tests, which build fully-populated rule inputs via the hydrators.
func TestSubPolicyQueryFieldsLive(t *testing.T) {
	token := os.Getenv("CATO_TOKEN")
	baseURL := os.Getenv("CATO_BASEURL")
	accountID := os.Getenv("CATO_ACCOUNT_ID")
	if token == "" || baseURL == "" || accountID == "" {
		t.Skip("CATO_TOKEN/CATO_BASEURL/CATO_ACCOUNT_ID not set; skipping live test")
	}

	client, err := cato.New(baseURL, token, accountID, nil, nil)
	if err != nil {
		t.Fatalf("client: %v", err)
	}
	ctx := context.Background()

	ifw, err := client.PolicyInternetFirewall(ctx, nil, accountID)
	if err != nil {
		t.Fatalf("PolicyInternetFirewall: %v", err)
	}
	ifp := ifw.GetPolicy().GetInternetFirewall().GetPolicy()
	// The selection must decode without error; fields are addressable.
	_ = ifp.GetSubPolicies()
	for _, rp := range ifp.GetRules() {
		_ = rp.GetRuleType()
		_ = rp.GetSubPolicy()
	}

	wan, err := client.PolicyWanFirewall(ctx, nil, accountID)
	if err != nil {
		t.Fatalf("PolicyWanFirewall: %v", err)
	}
	wfp := wan.GetPolicy().GetWanFirewall().GetPolicy()
	_ = wfp.GetSubPolicies()
	for _, rp := range wfp.GetRules() {
		_ = rp.GetRuleType()
		_ = rp.GetSubPolicy()
	}
}
