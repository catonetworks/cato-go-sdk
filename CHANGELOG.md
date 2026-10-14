# Changelog

## [1.0.0] - 2026-10-14

First major release of the Cato Networks Go SDK. This entry describes the
changes since the `v0.4.2` tag, including the generated SDK updates in
[#60](https://github.com/catonetworks/cato-go-sdk/pull/60),
[#62](https://github.com/catonetworks/cato-go-sdk/pull/62), and
[#64](https://github.com/catonetworks/cato-go-sdk/pull/64).

### Release contents

- 628 canonical GraphQL operations: 493 mutations and 135 queries, with typed
  client methods, response types, document constants, and generated models.
- Compared with `v0.4.2`, 15 generated client methods were added, no existing
  generated client methods were removed, and 99 existing method signatures
  changed. Generated response types and selections also changed, so retaining
  a method name does not guarantee source compatibility.
- The Go module remains `github.com/catonetworks/cato-go-sdk`; version 1 does
  not add a `/v1` suffix. The module requires Go 1.26.9 or newer.
- Default HTTP requests identify the release as `Cato_Go_SDK/v1.0.0`.

### Added

- **SDP static IP allocation policy:** 13 `PolicySdpStaticIPAllocation*`
  methods cover policy reads; rule and section creation, updates, movement,
  and removal; policy updates; and revision creation, discard, and publication.
  Generated inputs and payloads include user-based source selection and static
  IP allocation. This API namespace is marked beta in the schema.
- **Licensing service trials:** `LicensingStartServiceTrial`, its input and
  payload types, and trial service, plan, status, start date, and expiration
  fields. This operation is marked beta in the schema.
- **Internet Firewall revision discard:**
  `MutationPolicyInternetFirewallDiscardPolicyRevision` provides the imported
  CLI operation alongside the existing SDK operation with a distinct name.
- **App Connector pooled bandwidth:** allocation inputs and payloads for
  connector creation and connector query/mutation responses, plus `SocketRef`.
- **Posture metadata:** effective severity, suggested score, criticality,
  blast radius, findings signal, and score contribution metadata; aggregate
  posture summary fields; and expanded check-result and finding sort inputs.
- **Networking and policy fields:** BFD `MultiHop`, DHCP lease-time override
  settings, Site Web Proxy `AllowedDestinationProtocol`, and TLS inspection
  `ConnectionsOriginList`.
- **Additional schema models:** AI Security agent events and sessions,
  Business Platform account overviews, device crown-jewel classification and
  filtering, and application catalog FQDN/IP filtering. Schema models alone
  do not imply a dedicated callable SDK method for every nested API field.
- **Pagination totals:** `TotalCount` in application statistics, events, and
  Socket port metrics models and updated operation selections. Schema fields
  `total` for application statistics and events are deprecated in favor of
  `totalCount`.

### Breaking changes

- **Go toolchain:** the minimum Go version increases from 1.26.6 to 1.26.9.
  Update local and CI toolchains before building the release.
- **Policy mutation arguments:** affected Internet Firewall, WAN Firewall,
  WAN Network, TLS inspection, Socket LAN, Application Control, App Tenant
  Restriction, Private Access, and Remote Port Forwarding methods accept
  additional namespace input arguments. Many are optional GraphQL inputs,
  but callers must still supply the corresponding Go argument, typically
  `nil` when no policy context is needed. Check each generated signature:
  argument placement differs between operations.
- **Policy query arguments:** `InternetFirewallPolicy`, `WanFirewallPolicy`,
  `WanNetworkPolicy`, `Tlsinspectpolicy`, `ApplicationControlPolicy`,
  `AppTenantRestrictionPolicy`, `RemotePortFwdPolicy`, and
  `PolicyReadPrivateAccessPolicy` have updated input arguments.
- **Group mutations:** `GroupsCreateGroup`, `GroupsUpdateGroup`, and
  `GroupsDeleteGroup` now require a `GroupMembersListInput` argument for the
  selected members response. An empty input is shown in the group CRUD example.
- **Analytics and inventory arguments:** `AccountMetrics`, `AppStats`,
  `AppStatsTimeSeries`, `Events`, `EventsTimeSeries`, `AiSecurity`, and
  `HardwareManagement` have additional arguments required by their refreshed
  selections, including nested query inputs and analytics options.
- **User import type:** the retired `ImportType` field was removed from
  `User`, `UserFilterInput`, and `UserSortInput`. Use `UserImportType` and the
  corresponding generated getters, filter inputs, and sort fields.
- **Generated response shapes:** nested selections, aliases, union/interface
  projections, and generated type names reflect the normalized CLI documents.
  Recompile response readers and mocks, even where the client call is unchanged.

### Fixed and changed

- Raised the minimum Go patch version to 1.26.9 to address nine reachable
  standard-library advisories reported by `govulncheck` on Go 1.26.6:
  GO-2026-6603, GO-2026-6605, GO-2026-6607, GO-2026-6608, GO-2026-6610,
  GO-2026-6611, GO-2026-6612, GO-2026-6613, and GO-2026-6617.
- Regenerated the client and models from the refreshed GraphQL schema and
  normalized CLI operations, including current policy metadata, connector,
  posture, user, and networking selections.
- Preserved stable SDK operation names, compatible variable names and order,
  response aliases, SDK-only documents, and handwritten helpers during import.
  New variables are appended deterministically where compatibility allows.
- Reserved historical and handwritten Go identifiers and resolved new operation
  collisions with stable names instead of removing existing APIs.
- Preserved valid existing variable bindings when the current schema still
  accepts them; ambiguous variable renames now fail validation.
- Preserved the established `sites` response alias for matching site mutations
  and explicit Socket LAN section ownership.
- Accepted user documents that already omit retired `importType` selections;
  removed obsolete XDR analyst-feedback curation so valid CLI selections survive.
- Removed obsolete DHCP lease schema patches now covered by the upstream schema.
- Updated firewall, group, Socket LAN, TLS inspection, Remote Port Forwarding,
  and WAN Network examples for changed input and response shapes.

### Generation and CI

- Manifest version 2 records semantic operation keys, stable SDK names,
  ordered variables, CLI/document hashes, and the exact CLI producer commit.
  The release candidate records
  `ca604f63c4cd345c28cdc54b89cf6dabf07a0d0d` in
  `operations/manifest.json`. Its Git tree matches the squash-merged CLI
  `main` commit `9a7dd5fae5bfc5107f0195148cbee188ef6e596f` from CLI PR #26;
  the producer SHA is retained as the immutable source revision.
- The importer validates the complete candidate operation set before installing
  files. Operation counts are discovered from the CLI; explicit
  `EXPECTED_OPERATIONS` checks remain available.
- Schema-update CI requires and verifies a full 40-character CLI commit SHA,
  runs the shared operation-sync flow, and records the CLI dependency in the
  generated SDK PR. The CLI dependency must be reconciled with CLI `main`
  before tagging or publishing the SDK.
- `make operations-check` validates canonical documents and manifest coverage;
  `make generate-check` verifies deterministic client and model generation.
- Jenkins acceptance-test skip configuration was synchronized. A successful
  SDK build or unit test run does not establish live API acceptance coverage.

### Release preparation validation

- On Go 1.26.9, canonical operation validation, deterministic regeneration,
  unit tests, the SDK build, and all 50 example program builds passed.
- Lint against the release-preparation diff passed with zero issues.
- `govulncheck` reported zero reachable vulnerabilities. It also reported one
  advisory in imported packages and two in required modules whose vulnerable
  symbols are not called by this code.
- Live API tests were excluded from the local release checks. The existing
  `f406357` SDK commit has a successful `jenkins/sdk-acctest` status; that
  status predates these release-preparation changes.

### Upgrading from 0.x

1. Update the dependency with
   `go get github.com/catonetworks/cato-go-sdk@v1.0.0` once the tag is published.
2. Update affected call sites to match the generated signatures. For example,
   `PolicyInternetFirewallAddRule(ctx, input, accountID)` becomes
   `PolicyInternetFirewallAddRule(ctx, input, accountID, nil)` when no policy
   mutation context is needed.
3. Replace retired user import-type fields, update response accessors, and
   regenerate mocks against the new client interface. Prefer keyed model
   literals so added fields do not break positional struct initialization.
4. Build and test the consuming application and any examples it uses. The SDK's
   root module intentionally ignores `examples`, so `go test ./...` does not
   compile those programs.
5. Validate the operations you use against the intended account and endpoint.
   API permissions, feature entitlements, beta availability, and deployed
   schema support remain account/environment dependent.

## 2026-03-04

- Added latest methods in schema excluding sub-policy and AI Security operations.

[1.0.0]: https://github.com/catonetworks/cato-go-sdk/compare/v0.4.2...v1.0.0
