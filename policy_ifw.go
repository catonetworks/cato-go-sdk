package cato_go_sdk

import (
	"context"

	"github.com/gqlgo/gqlgenc/clientv2"
	cato_models "github.com/catonetworks/cato-go-sdk/models"
)

type IfwRulesIndexPolicy struct {
	Policy *IfwRulesIndexPolicy_Policy "json:\"policy,omitempty\" graphql:\"policy\""
}

type IfwRulesIndexPolicy_Policy struct {
	InternetFirewall *IfwRulesIndexPolicy_Policy_Policy "json:\"internetFirewall,omitempty\" graphql:\"internetFirewall\""
}

type IfwRulesIndexPolicy_Policy_Policy struct {
	Policy Policy_PIfwRulesIndexPolicy_Policy_Policy "json:\"policy\" graphql:\"policy\""
}

type Policy_PIfwRulesIndexPolicy_Policy_Policy struct {
	Enabled bool                                               "json:\"enabled\" graphql:\"enabled\""
	Rules   []*Policy_PIfwRulesIndexPolicy_Policy_Policy_Rules "json:\"rules\" graphql:\"rules\""
	// Sections []*Policy_Policy_InternetFirewall_Policy_Sections  "json:\"sections\" graphql:\"sections\""
}

type Policy_PIfwRulesIndexPolicy_Policy_Policy_Rules struct {
	Properties []cato_models.PolicyElementPropertiesEnum            "json:\"properties\" graphql:\"properties\""
	Rule       Policy_PIfwRulesIndexPolicy_Policy_Policy_Rules_Rule "json:\"rule\" graphql:\"rule\""
}

type Policy_PIfwRulesIndexPolicy_Policy_Policy_Rules_Rule struct {
	Description string                                                   "json:\"description\" graphql:\"description\""
	Enabled     bool                                                     "json:\"enabled\" graphql:\"enabled\""
	ID          string                                                   "json:\"id\" graphql:\"id\""
	Index       int64                                                    "json:\"index\" graphql:\"index\""
	Name        string                                                   "json:\"name\" graphql:\"name\""
	Section     Policy_Policy_InternetFirewall_Policy_Rules_Rule_Section "json:\"section\" graphql:\"section\""
}

type IfwSectionsIndexPolicy struct {
	Policy *IfwSectionsIndexPolicy_Policy "json:\"policy,omitempty\" graphql:\"policy\""
}

type IfwSectionsIndexPolicy_Policy struct {
	InternetFirewall *IfwSectionsIndexPolicy_Policy_Policy "json:\"internetFirewall,omitempty\" graphql:\"internetFirewall\""
}

type IfwSectionsIndexPolicy_Policy_Policy struct {
	Policy Policy_Policy_IfwSectionsIndexPolicy_Policy_Policy "json:\"policy\" graphql:\"policy\""
}

type Policy_Policy_IfwSectionsIndexPolicy_Policy_Policy struct {
	// Enabled bool                                               "json:\"enabled\" graphql:\"enabled\""
	// Rules   []*Policy_PIfwRulesIndexPolicy_Policy_Policy_Rules "json:\"rules\" graphql:\"rules\""
	Sections []*Policy_Policy_InternetFirewall_Policy_Sections "json:\"sections\" graphql:\"sections\""
}

func (c *Client) PolicyInternetFirewallRulesIndex(ctx context.Context, accountID string, interceptors ...clientv2.RequestInterceptor) (*IfwRulesIndexPolicy, error) {
	vars := map[string]any{
		"accountId": accountID,
	}

	var res IfwRulesIndexPolicy
	if err := c.Client.Post(ctx, "policy", PolicyDocumentInternetFirewallRulesIndex, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}

const PolicyDocumentInternetFirewallRulesIndex = `query policy ( $accountId:ID! ) {
	policy ( accountId:$accountId  ) {
		internetFirewall  {
			policy {
				enabled 
				rules  {
					rule {
						id 
						name 
						index
						description
						section  {
							id
							name
						}
						enabled 	
					}
					properties
				}
			}
		}
	}	
}`

func (c *Client) PolicyInternetFirewallSectionsIndex(ctx context.Context, accountID string, interceptors ...clientv2.RequestInterceptor) (*IfwSectionsIndexPolicy, error) {
	vars := map[string]any{
		"accountId": accountID,
	}

	var res IfwSectionsIndexPolicy
	if err := c.Client.Post(ctx, "policy", PolicyDocumentInternetFirewallSectionsIndex, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}

const PolicyDocumentInternetFirewallSectionsIndex = `query policy ( $accountId:ID! ) {
	policy ( accountId:$accountId  ) {
		internetFirewall  {
			policy {
				sections  {
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
			}
		}

	}	
}`

func (c *Client) PolicyInternetFirewall(ctx context.Context, internetFirewallPolicyInput *cato_models.InternetFirewallPolicyInput, accountID string, interceptors ...clientv2.RequestInterceptor) (*Policy, error) {
	vars := map[string]any{
		"internetFirewallPolicyInput": internetFirewallPolicyInput,
		"accountId":                   accountID,
	}

	var res Policy
	if err := c.Client.Post(ctx, "policy", PolicyDocumentInternetFirewall, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}

const PolicyDocumentInternetFirewall = `query policy ($internetFirewallPolicyInput: InternetFirewallPolicyInput, $accountId: ID!) {
	policy(accountId: $accountId) {
		internetFirewall {
			policy(input: $internetFirewallPolicyInput) {
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
							ip
							host {
								id
								name
							}
							site {
								id
								name
							}
							subnet
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
						deviceAttributes  {
							category
							type
							model
							manufacturer
							os
							osVersion
						}						
						destination {
							application {
								id
								name
							}
							customApp {
								id
								name
							}
							appCategory {
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
							country {
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
							remoteAsn
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

						exceptions {
							name
							source {
								ip
								host {
									id
									name
								}
								site {
									id
									name
								}
								subnet
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
						
							country {
								id
								name
							}
							device {
								id
								name
							}
							destination {
								application {
									id
									name
								}
								customApp {
									id
									name
								}
								appCategory {
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
								country {
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
								remoteAsn
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
			revisionsInternetFirewallPolicyQueries: revisions {
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
