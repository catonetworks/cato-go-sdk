package cato_go_sdk

import (
	"context"

	"github.com/Yamashou/gqlgenc/clientv2"
	cato_models "github.com/cato-networks/cato-go-sdk/models"
)

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
