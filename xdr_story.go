package cato_go_sdk

import (
	"context"

	"github.com/Yamashou/gqlgenc/clientv2"
)

type Xdr_Story struct {
	Xdr Xdr_Xdr_Story "json:\"xdr\" graphql:\"xdr\""
}

type Xdr_Xdr_Story struct {
	Story *XdrStoryRespData "json:\"story,omitempty\" graphql:\"story\""
}

type XdrStoryRespData struct {
	AccountID    int64                             "json:\"accountId\" graphql:\"accountId\""
	AccountName  *string                           "json:\"accountName,omitempty\" graphql:\"accountName\""
	AnalystEmail *string                           "json:\"analystEmail,omitempty\" graphql:\"analystEmail\""
	AnalystName  *string                           "json:\"analystName,omitempty\" graphql:\"analystName\""
	CreatedAt    string                            "json:\"createdAt\" graphql:\"createdAt\""
	ID           string                            "json:\"id\" graphql:\"id\""
	Incident     Xdr_Xdr_Stories_Items_Incident    "json:\"incident\" graphql:\"incident\""
	Playbook     *string                           "json:\"playbook,omitempty\" graphql:\"playbook\""
	Summary      *string                           "json:\"summary,omitempty\" graphql:\"summary\""
	Timeline     []*Xdr_Xdr_Stories_Items_Timeline "json:\"timeline\" graphql:\"timeline\""
	UpdatedAt    string                            "json:\"updatedAt\" graphql:\"updatedAt\""
}

func (c *Client) XdrStory(ctx context.Context, storyId *string, incidentId *string, accountID string, interceptors ...clientv2.RequestInterceptor) (*Xdr_Story, error) {
	vars := map[string]any{
		"storyId":    *storyId,
		"incidentId": nil,
		"accountID":  accountID,
	}

	var res Xdr_Story
	if err := c.Client.Post(ctx, "xdr", XdrStoryDocument, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}

const XdrStoryDocument = `query xdr ( $storyId:ID $incidentId:ID $accountID:ID! ) {
	xdr ( accountID:$accountID  ) {
		story ( storyId:$storyId  incidentId:$incidentId  )  {
			id
			accountId
			analystName
			analystEmail
			accountName
			updatedAt
			createdAt
			playbook
			summary
			incident {
				id
				firstSignal
				lastSignal
				engineTypeMergedIncident: engineType
				vendorMergedIncident: vendor
				producer
				producerName
				connectionTypeMergedIncident: connectionType
				indication
				queryName
				criticality
				source
				ticket
				statusMergedIncident: status
				research
				siteName
				storyDuration
				description
				sourceIp
				analystFeedbackMergedIncident: analystFeedback {
					verdict
					severity
					threatType {
						name
						recommendedAction
						details
					}
					threatClassification
					additionalInfo
				}
				siteMergedIncident: site {
					id
					name
				}
				userMergedIncident: user {
					id
					name
				}
				predictedVerdictMergedIncident: predictedVerdict
				predictedThreatType
				... on MicrosoftEndpoint {
					similarStoriesData {
						storyId
						threatTypeName
						verdict
						threatClassification
						similarityPercentage
						indication
					}
					device {
						id
						deviceName
						osDetailsMicrosoftDeviceDetails: osDetails {
							osType
							osBuild
							osVersion
						}
						loggedOnUsersMicrosoftDeviceDetails: loggedOnUsers {
							id
							name
							... on MicrosoftEndpointUser {
								userSid
								accountName
								domainName
								principalName
							}
							... on CatoEndpointUser {
								id
								name
							}
						}
						firstSeenDateTime
						avStatusMicrosoftDeviceDetails: avStatus
						healthStatusMicrosoftDeviceDetails: healthStatus
						rbacGroupMicrosoftDeviceDetails: rbacGroup {
							id
							name
						}
						ipInterfaces
						azureAdDeviceId
						onboardingStatusMicrosoftDeviceDetails: onboardingStatus
					}
					alerts {
						id
						title
						description
						threatName
						mitreTechnique {
							id
							name
						}
						mitreSubTechnique {
							id
							name
						}
						createdDateTime
						resources {
							id
							createdDateTime
							remediationStatus
							remediationStatusDetails
							tags
							roles
							verdict
							... on MicrosoftProcessResource {
								processId
								processCommandLine
								imageFile {
									name
									path
									size
									sha1
									sha256
									md5
									issuer
									signer
									publisher
								}
								userAccount {
									id
									name
									... on MicrosoftEndpointUser {
										userSid
										accountName
										domainName
										principalName
									}
									... on CatoEndpointUser {
										id
										name
									}
								}
							}
							... on MicrosoftFileResource {
								fileDetails {
									name
									path
									size
									sha1
									sha256
									md5
									issuer
									signer
									publisher
								}
								detectionStatus
							}
							... on MicrosoftRegistryResource {
								hive
								key
								value
								valueName
								valueType
							}
						}
						activities {
							id
							resourceId
							parentResourceId
							action
							firstActivityDateTime
							lastActivityDateTime
						}
						criticality
						comments
						recommendedActions
						category
						ownerName
						threatFamilyName
						threatType
						resolvedDateTime
						firstActivityDateTime
						lastActivityDateTime
						lastUpdateDateTime
						localIp
						destinationIp
						destinationUrl
						statusMicrosoftDefenderEndpointAlert: status
						providerAlertId
						alertWebUrl
						determinationMicrosoftDefenderEndpointAlert: determination
						detectionSourceMicrosoftDefenderEndpointAlert: detectionSource
						classificationMicrosoftDefenderEndpointAlert: classification
					}
				}
				... on AnomalyStats {
					similarStoriesData {
						storyId
						threatTypeName
						verdict
						threatClassification
						similarityPercentage
						indication
					}
					srcSiteId
					os
					deviceName
					macAddress
					logonName
					clientClass
					drillDownFilter {
						name
						value
					}
					breakdownField
					subjectType
					extra {
						name
						type
						value
					}
					gaussian {
						std
						ss
						z_score
						avg
						n
					}
					metric {
						name
						value
					}
					metricDetails {
						name
						units
					}
					mitres {
						id
						name
					}
					rules
					timeSeries {
						data
						label
						sum
						unitsIncidentTimeseries: units
						info
						keyIncidentTimeseries: key {
							measureFieldName
							dimensions {
								fieldName
								value
							}
						}
					}
					targets {
						typeIncidentTargetRep: type
						name
						analysisScore
						infectionSource
						threatReference
						catoPopularity
						threatFeeds
						creationTime
						categories
						countryOfRegistration
						searchHits
						engines
						eventData {
							signatureId
							eventType
							threatType
							threatName
							severity
							action
							ruleId
							virusName
							scanResult
							appId
							appName
							dnsProtectionCategory
						}
					}
					direction
				}
				... on AnomalyEvents {
					similarStoriesData {
						storyId
						threatTypeName
						verdict
						threatClassification
						similarityPercentage
						indication
					}
					srcSiteId
					os
					deviceName
					macAddress
					logonName
					clientClass
					drillDownFilter {
						name
						value
					}
					breakdownField
					subjectType
					extra {
						name
						type
						value
					}
					gaussian {
						std
						ss
						z_score
						avg
						n
					}
					metric {
						name
						value
					}
					metricDetails {
						name
						units
					}
					mitres {
						id
						name
					}
					rules
					timeSeries {
						data
						label
						sum
						unitsIncidentTimeseries: units
						info
						keyIncidentTimeseries: key {
							measureFieldName
							dimensions {
								fieldName
								value
							}
						}
					}
					targets {
						typeIncidentTargetRep: type
						name
						analysisScore
						infectionSource
						threatReference
						catoPopularity
						threatFeeds
						creationTime
						categories
						countryOfRegistration
						searchHits
						engines
						eventData {
							signatureId
							eventType
							threatType
							threatName
							severity
							action
							ruleId
							virusName
							scanResult
							appId
							appName
							dnsProtectionCategory
						}
					}
					direction
				}
				... on Threat {
					similarStoriesData {
						storyId
						threatTypeName
						verdict
						threatClassification
						similarityPercentage
						indication
					}
					srcSiteId
					flowsCardinality
					riskLevel
					os
					deviceName
					macAddress
					logonName
					direction
					clientClass
					events {
						signatureId
						eventType
						threatType
						threatName
						severity
						action
						ruleId
						virusName
						scanResultEvent: scanResult
						appId
						appName
						dnsProtectionCategory
					}
					mitres {
						id
						name
					}
					timeSeries {
						data
						label
						sum
						unitsIncidentTimeseries: units
						info
						keyIncidentTimeseries: key {
							measureFieldName
							dimensions {
								fieldName
								value
							}
						}
					}
					targets {
						typeIncidentTargetRep: type
						name
						analysisScore
						infectionSource
						threatReference
						catoPopularity
						threatFeeds
						creationTime
						categories
						countryOfRegistration
						searchHits
						engines
						eventData {
							signatureId
							eventType
							threatType
							threatName
							severity
							action
							ruleId
							virusName
							scanResult
							appId
							appName
							dnsProtectionCategory
						}
					}
					flows {
						appName
						clientClass
						sourceIp
						sourcePort
						destinationCountry
						destinationIp
						destinationPort
						direction
						createdAt
						referer
						userAgent
						method
						url
						target
						domain
						sourceGeolocation
						destinationGeolocation
						tunnelGeolocation
						httpResponseCode
						dnsResponseIP
						smbFileName
						user
						fileHash
						ja3
					}
				}
				... on ThreatPrevention {
					similarStoriesData {
						storyId
						threatTypeName
						verdict
						threatClassification
						similarityPercentage
						indication
					}
					srcSiteId
					flowsCardinality
					riskLevel
					os
					deviceName
					macAddress
					logonName
					direction
					clientClass
					events {
						signatureId
						eventType
						threatType
						threatName
						severity
						action
						ruleId
						virusName
						scanResultEvent: scanResult
						appId
						appName
						dnsProtectionCategory
					}
					mitres {
						id
						name
					}
					timeSeries {
						data
						label
						sum
						unitsIncidentTimeseries: units
						info
						keyIncidentTimeseries: key {
							measureFieldName
							dimensions {
								fieldName
								value
							}
						}
					}
					targets {
						typeIncidentTargetRep: type
						name
						analysisScore
						infectionSource
						threatReference
						catoPopularity
						threatFeeds
						creationTime
						categories
						countryOfRegistration
						searchHits
						engines
						eventData {
							signatureId
							eventType
							threatType
							threatName
							severity
							action
							ruleId
							virusName
							scanResult
							appId
							appName
							dnsProtectionCategory
						}
					}
					threatPreventionsEvents {
						appName
						clientClass
						sourceIp
						sourcePort
						destinationCountry
						destinationIp
						destinationPort
						direction
						createdAt
						method
						url
						target
						domain
						sourceGeolocation
						destinationGeolocation
						tunnelGeolocation
						dnsResponseIP
						smbFileName
						user
						userAgent
						fileHash
						ja3
						referrer
						httpResponseCode
					}
				}
				... on NetworkXDRIncident {
					similarStoriesData {
						storyId
						threatTypeName
						verdict
						threatClassification
						similarityPercentage
						indication
					}
					networkIncidentTimeline {
						created
						validated
						description
						eventTypeNetworkTimelineEvent: eventType
						incidentId
						networkEventSourceNetworkTimelineEvent: networkEventSource
						eventIds
						acknowledged
						linkId
						linkName
						linkConfigPrecedenceNetworkTimelineEvent: linkConfigPrecedence
						linkStatusNetworkTimelineEvent: linkStatus
						linkConfigBandwidth
						deviceConfigHaRoleNetworkTimelineEvent: deviceConfigHaRole
						deviceHaRoleStateNetworkTimelineEvent: deviceHaRoleState
						socketSerialId
						pop
						isp
						bgpConnectionNetworkTimelineEvent: bgpConnection {
							connectionName
							peerIp
							peerAsn
							catoIp
							catoAsn
						}
						linkQualityIssueNetworkTimelineEvent: linkQualityIssue {
							issueType
							direction
							current
							threshold
						}
						hostIp
						ruleName
						tunnelResetCount
						muted
					}
					storyType
					occurrences
					siteConnectionType
					siteConfigLocation
					acknowledged
					linkId
					linkName
					linkConfigPrecedence
					deviceConfigHaRole
					licenseRegion
					licenseBandwidth
					pop
					isp
					bgpConnection {
						connectionName
						peerIp
						peerAsn
						catoIp
						catoAsn
					}
					hostIp
					ruleName
					muted
					ilmmDetails {
						linkDetailsIlmmDetails: linkDetails {
							linkId
							description
							ispLinkId
							comments
							onboardingStatus
							activeLicense
						}
						ispDetailsIlmmDetails: ispDetails {
							name
							ispAccountId
							supportEmail
							supportPhone
							description
							countryCode
							loaFile {
								fileName
								fileHash
								uploadedAt
							}
						}
						contacts {
							name
							phone
							email
						}
					}
				}
				... on CatoEndpoint {
					similarStoriesData {
						storyId
						threatTypeName
						verdict
						threatClassification
						similarityPercentage
						indication
					}
					device {
						id
						deviceName
						osDetailsCatoEndpointDeviceDetails: osDetails {
							osType
							osBuild
							osVersion
						}
						loggedOnUsersCatoEndpointDeviceDetails: loggedOnUsers {
							id
							name
							... on MicrosoftEndpointUser {
								userSid
								accountName
								domainName
								principalName
							}
							... on CatoEndpointUser {
								id
								name
							}
						}
						macAddress
					}
					alerts {
						id
						title
						description
						threatName
						mitreTechnique {
							id
							name
						}
						mitreSubTechnique {
							id
							name
						}
						createdDateTime
						resources {
							id
							createdDateTime
							remediationStatus
							... on CatoProcessResource {
								processId
								processCommandLine
								imageFile {
									name
									path
									size
									sha1
									sha256
									md5
									issuer
									signer
									publisher
								}
								userAccount {
									id
									name
									... on MicrosoftEndpointUser {
										userSid
										accountName
										domainName
										principalName
									}
									... on CatoEndpointUser {
										id
										name
									}
								}
							}
							... on CatoFileResource {
								fileDetails {
									name
									path
									size
									sha1
									sha256
									md5
									issuer
									signer
									publisher
								}
								detectionStatus
							}
						}
						activities {
							id
							resourceId
							parentResourceId
						}
						criticality
						engineTypeCatoEndpointAlert: engineType
						statusCatoEndpointAlert: status
						endpointProtectionProfile
					}
				}
			}
			timeline {
				createdAt
				description
				context
				type
				descriptions
				categoryTimelineItem: category
				additionalInfo
				analystInfoTimelineItem: analystInfo {
					name
					email
				}
			}
		}
	}	
}
`
