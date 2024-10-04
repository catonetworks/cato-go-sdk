package cato_go_sdk

import (
	"context"

	"github.com/Yamashou/gqlgenc/clientv2"
	cato_models "github.com/routebyintuition/cato-go-sdk/models"
)

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
							container {
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
								container {
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
