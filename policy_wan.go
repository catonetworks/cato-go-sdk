package cato_go_sdk

import (
	"context"

	"github.com/Yamashou/gqlgenc/clientv2"
	cato_models "github.com/catonetworks/cato-go-sdk/models"
)

type WanRulesIndexPolicy struct {
	Policy *WanRulesIndexPolicy_Policy "json:\"policy,omitempty\" graphql:\"policy\""
}

type WanRulesIndexPolicy_Policy struct {
	WanFirewall *WanRulesIndexPolicy_Policy_WanFirewall "json:\"wanFirewall,omitempty\" graphql:\"wanFirewall\""
}

type WanRulesIndexPolicy_Policy_WanFirewall struct {
	Policy WanRulesIndexPolicy_Policy_WanFirewall_Policy "json:\"policy\" graphql:\"policy\""
}

type WanRulesIndexPolicy_Policy_WanFirewall_Policy struct {
	Rules []*WanRulesIndexPolicy_Policy_WanFirewall_Policy_Rules "json:\"rules\" graphql:\"rules\""
}

type WanRulesIndexPolicy_Policy_WanFirewall_Policy_Rules struct {
	Properties []cato_models.PolicyElementPropertiesEnum                "json:\"properties\" graphql:\"properties\""
	Rule       WanRulesIndexPolicy_Policy_WanFirewall_Policy_Rules_Rule "json:\"rule\" graphql:\"rule\""
}

type WanRulesIndexPolicy_Policy_WanFirewall_Policy_Rules_Rule struct {
	Action      cato_models.WanFirewallActionEnum                   "json:\"action\" graphql:\"action\""
	Description string                                              "json:\"description\" graphql:\"description\""
	Direction   cato_models.WanFirewallDirectionEnum                "json:\"direction\" graphql:\"direction\""
	Enabled     bool                                                "json:\"enabled\" graphql:\"enabled\""
	ID          string                                              "json:\"id\" graphql:\"id\""
	Index       int64                                               "json:\"index\" graphql:\"index\""
	Name        string                                              "json:\"name\" graphql:\"name\""
	Section     Policy_Policy_WanFirewall_Policy_Rules_Rule_Section "json:\"section\" graphql:\"section\""
}

type WanSectionsIndexPolicy struct {
	Policy *WanSectionsIndexPolicy_Policy "json:\"policy,omitempty\" graphql:\"policy\""
}

type WanSectionsIndexPolicy_Policy struct {
	WanFirewall *Policy_Policy_WanFirewall "json:\"wanFirewall,omitempty\" graphql:\"wanFirewall\""
}

type WanSectionsIndexPolicy_Policy_WanFirewall struct {
	Policy WanSectionsIndexPolicy_Policy_WanFirewall_Policy "json:\"policy\" graphql:\"policy\""
}

type WanSectionsIndexPolicy_Policy_WanFirewall_Policy struct {
	// Rules    []*Policy_Policy_WanFirewall_Policy_Rules    "json:\"rules\" graphql:\"rules\""
	Sections []*Policy_Policy_WanFirewall_Policy_Sections "json:\"sections\" graphql:\"sections\""
}

func (c *Client) PolicyWanFirewallRulesIndex(ctx context.Context, accountID string, interceptors ...clientv2.RequestInterceptor) (*WanRulesIndexPolicy, error) {
	vars := map[string]any{
		"accountId": accountID,
	}

	var res WanRulesIndexPolicy
	if err := c.Client.Post(ctx, "policy", PolicyDocumentWanFirewallRulessIndex, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}

const PolicyDocumentWanFirewallRulessIndex = `query policy ( $accountId:ID! ) {
	policy ( accountId:$accountId  ) {
		wanFirewall  {
			policy {
				rules  {
					rule {
						id 
						name 
						description 
						index 
						section  {
							id
							name
						}
						enabled 
						action 
						direction 
				}
				properties
				}
			}
		}
	}	
}`

func (c *Client) PolicyWanFirewallSectionsIndex(ctx context.Context, accountID string, interceptors ...clientv2.RequestInterceptor) (*WanSectionsIndexPolicy, error) {
	vars := map[string]any{
		"accountId": accountID,
	}

	var res WanSectionsIndexPolicy
	if err := c.Client.Post(ctx, "policy", PolicyDocumentWanFirewallSectionsIndex, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}

const PolicyDocumentWanFirewallSectionsIndex = `query policy ( $accountId:ID! ) {
	policy ( accountId:$accountId  ) {
		wanFirewall  {
			policy {
				enabled 
				sections  {
					section {
						id 
						name 
					}
					properties
				}
			}
		}

	}	
}`

func (c *Client) PolicyWanFirewall(ctx context.Context, wanFirewallPolicyInput *cato_models.WanFirewallPolicyInput, accountID string, interceptors ...clientv2.RequestInterceptor) (*Policy, error) {
	vars := map[string]any{
		"wanFirewallPolicyInput": wanFirewallPolicyInput,
		"accountId":              accountID,
	}

	var res Policy
	if err := c.Client.Post(ctx, "policy", PolicyDocumentWanFirewall, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}

const PolicyDocumentWanFirewall = `query policy ($wanFirewallPolicyInput: WanFirewallPolicyInput, $accountId: ID!) {
	policy(accountId: $accountId) {
		wanFirewall {
			policy(input: $wanFirewallPolicyInput) {
				enabled
				rules {
					audit {
						updatedTime
						updatedBy
					}
					rule {
						id
						name
						description
						index
						section {
							id
							name
						}
						enabled
						source {
							host {
								id
								name
							}
							site {
								id
								name
							}
							subnet
							ip
							ipRange {
								from
								to
							}
							globalIpRange {
								id
								name
							}
							networkInterface {
								id
								name
							}
							siteNetworkSubnet {
								id
								name
							}
							floatingSubnet {
								id
								name
							}
							user {
								id
								name
							}
							usersGroup {
								id
								name
							}
							group {
								id
								name
							}
							systemGroup {
								id
								name
							}
						}
						connectionOrigin
						country {
							id
							name
						}
						device {
							id
							name
						}
						deviceOS
						destination {
							host {
								id
								name
							}
							site {
								id
								name
							}
							subnet
							ip
							ipRange {
								from
								to
							}
							globalIpRange {
								id
								name
							}
							networkInterface {
								id
								name
							}
							siteNetworkSubnet {
								id
								name
							}
							floatingSubnet {
								id
								name
							}
							user {
								id
								name
							}
							usersGroup {
								id
								name
							}
							group {
								id
								name
							}
							systemGroup {
								id
								name
							}
						}
						application {
							application {
								id
								name
							}
							appCategory {
								id
								name
							}
							customApp {
								id
								name
							}
							customCategory {
								id
								name
							}
							sanctionedAppsCategory {
								id
								name
							}
							domain
							fqdn
							ip
							subnet
							ipRange {
								from
								to
							}
							globalIpRange {
								id
								name
							}
						}
						service {
							standard {
								id
								name
							}
							custom {
								port
								portRange {
									from
									to
								}
								protocol
							}
						}
						action
						tracking {
							event {
								enabled
							}
							alert {
								enabled
								frequency
								subscriptionGroup {
									id
									name
								}
								webhook {
									id
									name
								}
								mailingList {
									id
									name
								}
							}
						}
						schedule {
							activeOn
							customTimeframePolicySchedule: customTimeframe {
								from
								to
							}
							customRecurringPolicySchedule: customRecurring {
								from
								to
								days
							}
						}

						activePeriod {
              				useEffectiveFrom
              				effectiveFrom
              				useExpiresAt
              				expiresAt
            			}

						direction
						exceptions {
							name
							source {
								host {
									id
									name
								}
								site {
									id
									name
								}
								subnet
								ip
								ipRange {
									from
									to
								}
								globalIpRange {
									id
									name
								}
								networkInterface {
									id
									name
								}
								siteNetworkSubnet {
									id
									name
								}
								floatingSubnet {
									id
									name
								}
								user {
									id
									name
								}
								usersGroup {
									id
									name
								}
								group {
									id
									name
								}
								systemGroup {
									id
									name
								}
							}
							deviceOS
							destination {
								host {
									id
									name
								}
								site {
									id
									name
								}
								subnet
								ip
								ipRange {
									from
									to
								}
								globalIpRange {
									id
									name
								}
								networkInterface {
									id
									name
								}
								siteNetworkSubnet {
									id
									name
								}
								floatingSubnet {
									id
									name
								}
								user {
									id
									name
								}
								usersGroup {
									id
									name
								}
								group {
									id
									name
								}
								systemGroup {
									id
									name
								}
							}
							country {
								id
								name
							}
							device {
								id
								name
							}
							application {
								application {
									id
									name
								}
								appCategory {
									id
									name
								}
								customApp {
									id
									name
								}
								customCategory {
									id
									name
								}
								sanctionedAppsCategory {
									id
									name
								}
								domain
								fqdn
								ip
								subnet
								ipRange {
									from
									to
								}
								globalIpRange {
									id
									name
								}
							}
							service {
								standard {
									id
									name
								}
								custom {
									port
									portRangeCustomService: portRange {
										from
										to
									}
									protocol
								}
							}
							connectionOrigin
							direction
						}
					}
					properties
				}
				sections {
					audit {
						updatedTime
						updatedBy
					}
					section {
						id
						name
					}
					properties
				}
				audit {
					publishedTime
					publishedBy
				}
				revision {
					id
					name
					description
					changes
					createdTime
					updatedTime
				}
			}
			revisionsWanFirewallPolicyQueries: revisions {
				revision {
					id
					name
					description
					changes
					createdTime
					updatedTime
				}
			}
		}
	}
}
`
