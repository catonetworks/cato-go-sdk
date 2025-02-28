package cato_go_sdk

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Yamashou/gqlgenc/clientv2"
	cato_models "github.com/catonetworks/cato-go-sdk/models"
)

// func (c *Client) SiteCloudInterconnectPhysicalConnection(ctx context.Context, cloudInterconnectPhysicalConnectionInput cato_models.CloudInterconnectPhysicalConnectionInput, accountID string, interceptors ...clientv2.RequestInterceptor) (*Site, error) {
// 	vars := map[string]any{
// 		"cloudInterconnectPhysicalConnectionInput": cloudInterconnectPhysicalConnectionInput,
// 		"accountId": accountID,
// 	}

// 	var res Site
// 	if err := c.Client.Post(ctx, "site", SiteCloudInterconnectPhysicalConnectionDocument, &res, vars, interceptors...); err != nil {
// 		if c.Client.ParseDataWhenErrors {
// 			return &res, err
// 		}

// 		return nil, err
// 	}

// 	return &res, nil
// }

// const SiteCloudInterconnectPhysicalConnectionDocument = `query site ( $cloudInterconnectPhysicalConnectionInput:CloudInterconnectPhysicalConnectionInput! $accountId:ID! ) {
// 	site ( accountId:$accountId  ) {
// 		cloudInterconnectPhysicalConnection ( input:$cloudInterconnectPhysicalConnectionInput ) {
// 			id
// 			site {
// 				id
// 				name
// 			}
// 			haRole
// 			popLocation {
// 				id
// 				name
// 			}
// 			serviceProviderName
// 			encapsulationMethod
// 			subnet
// 			privateCatoIp
// 			privateSiteIp
// 			upstreamBwLimit
// 			downstreamBwLimit
// 			vlan
// 			sVlan
// 			cVlan
// 		}
// 	}
// }
// `

// func (c *Client) SiteCloudInterconnectPhysicalConnectionId(ctx context.Context, accountID string, interceptors ...clientv2.RequestInterceptor) (*Site, error) {
// 	vars := map[string]any{
// 		"accountId": accountID,
// 	}

// 	var res Site
// 	if err := c.Client.Post(ctx, "site", SiteCloudInterconnectPhysicalConnectionIdDocument, &res, vars, interceptors...); err != nil {
// 		if c.Client.ParseDataWhenErrors {
// 			return &res, err
// 		}

// 		return nil, err
// 	}

// 	return &res, nil
// }

// const SiteCloudInterconnectPhysicalConnectionIdDocument = `query site ( $cloudInterconnectPhysicalConnectionIdInput:CloudInterconnectPhysicalConnectionIdInput! $accountId:ID! ) {
// 	site ( accountId:$accountId  ) {

// 		cloudInterconnectPhysicalConnectionId ( input:$cloudInterconnectPhysicalConnectionIdInput  )  {
// 			id
// 		}

// 	}
// }
// `

// func (c *Client) SiteCloudInterconnectConnectionConnectivity(ctx context.Context, accountID string, interceptors ...clientv2.RequestInterceptor) (*Site, error) {
// 	vars := map[string]any{
// 		"accountId": accountID,
// 	}

// 	var res Site
// 	if err := c.Client.Post(ctx, "site", SiteCloudInterconnectConnectionConnectivityDocument, &res, vars, interceptors...); err != nil {
// 		if c.Client.ParseDataWhenErrors {
// 			return &res, err
// 		}

// 		return nil, err
// 	}

// 	return &res, nil
// }

// const SiteCloudInterconnectConnectionConnectivityDocument = `query site ( $accountId:ID! ) {
// 	site ( accountId:$accountId  ) {
// 		cloudInterconnectConnectionConnectivity  {
// 			success
// 		}
// 	}
// }
// `

func (c *Client) SiteBgpPeer(ctx context.Context, bgpPeerRefInput cato_models.BgpPeerRefInput, accountID string, interceptors ...clientv2.RequestInterceptor) (*Site, error) {
	vars := map[string]any{
		"bgpPeerRefInput": bgpPeerRefInput,
		"accountId":       accountID,
	}

	var res Site
	if err := c.Client.Post(ctx, "site", SiteBgpPeerDocument, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}

const SiteBgpPeerDocument = `query site ( $bgpPeerRefInput:BgpPeerRefInput! $accountId:ID! ) {
	site ( accountId:$accountId  ) {
		bgpPeer ( input:$bgpPeerRefInput  )  {
			site {
				id 
				name 
			}
			id
			name
			peerAsn
			catoAsn
			peerIp
			catoIp
			advertiseDefaultRoute
			advertiseAllRoutes
			advertiseSummaryRoutes
			summaryRoute {
				id 
				route 
				community  {
					from
					to
				}

			}
			defaultAction
			performNat
			md5AuthKey
			metric
			holdTime
			keepaliveInterval
			bfdEnabled
			bfdSettingsBgpPeer: bfdSettings {
				transmitInterval 
				receiveInterval 
				multiplier 
			}
			trackingBgpPeer: tracking {
				id 
				enabled 
				alertFrequency 
				subscriptionId 
			}
		}

	}	
}
`

func (c *Client) SiteBgpPeerList(ctx context.Context, bgpPeerListInput cato_models.BgpPeerListInput, accountID string, interceptors ...clientv2.RequestInterceptor) (*Site, error) {
	vars := map[string]any{
		"bgpPeerListInput": bgpPeerListInput,
		"accountId":        accountID,
	}

	var res Site
	if err := c.Client.Post(ctx, "site", SiteBgpPeerListDocument, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}

const SiteBgpPeerListDocument = `query site ( $bgpPeerListInput:BgpPeerListInput! $accountId:ID! ) {
	site ( accountId:$accountId  ) {
		bgpPeerList ( input:$bgpPeerListInput  )  {
			bgpPeerBgpPeerListPayload: bgpPeer {
				site  {
					id
					name
				}

				id 
				name 
				peerAsn 
				catoAsn 
				peerIp 
				catoIp 
				advertiseDefaultRoute 
				advertiseAllRoutes 
				advertiseSummaryRoutes 
				summaryRoute  {
					id
					route
					community {
						from 
						to 
					}
				}

				defaultAction 
				performNat 
				md5AuthKey 
				metric 
				holdTime 
				keepaliveInterval 
				bfdEnabled 
				bfdSettings  {
					transmitInterval
					receiveInterval
					multiplier
				}

				tracking  {
					id
					enabled
					alertFrequency
					subscriptionId
				}

			}
			total
		}
	}	
}
`

func (c *Client) SiteBgpStatus(ctx context.Context, siteBgpStatusInput cato_models.SiteBgpStatusInput, accountID string, interceptors ...clientv2.RequestInterceptor) (*Site, []*SiteBgpStatusResult, error) {
	vars := map[string]any{
		"siteBgpStatusInput": siteBgpStatusInput,
		"accountId":          accountID,
	}

	var res Site
	if err := c.Client.Post(ctx, "site", SiteBgpStatusDocument, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, nil, err
		}

		return nil, nil, err
	}

	siteBgpStatusResult := []*SiteBgpStatusResult{}
	for _, val := range res.Site.SiteBgpStatus.RawStatus {

		tmpVal := SiteBgpStatusResult{}
		err2 := json.Unmarshal([]byte(val), &tmpVal)
		if err2 != nil {
			fmt.Println("Error parsing JSON:", err2)
		}
		siteBgpStatusResult = append(siteBgpStatusResult, &tmpVal)
	}

	return &res, siteBgpStatusResult, nil
}

const SiteBgpStatusDocument = `query site ( $siteBgpStatusInput:SiteBgpStatusInput! $accountId:ID! ) {
	site ( accountId:$accountId  ) {
		siteBgpStatus ( input:$siteBgpStatusInput  )  {
			status {
				remoteIp 
				bgpSession 
				bfdSession 
				routesFromPeer 
				routesToPeer 
				rejectedRoutesFromPeer  {
					subnet
					type
					community {
						from 
						to 
					}
					rule
					lastPublishAttempt
				}

			}
			rawStatus
		}
	}	
}
`
