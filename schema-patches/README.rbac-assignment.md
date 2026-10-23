# RBAC assignment contract

Source: server-services commit `77770519fd3c531ef830eb86b0fb16c9077bc8d8`,
`services/iam/src/main/resources/schema/rbac.public.graphqls` and
`RBACAssignmentDataService.java` / `RBACAssignmentConverter.java`.

Adds the undocumented public assignment operations under `rbac(accountId:)`.
Private fields and enum values are excluded. `paging.from` is a zero-based page
number in the assignment service (Spring PageRequest), not a row offset.
Updates replace role/entity grants and omit create-only principal arrays.
An omitted obfuscation preference leaves it unchanged. Legacy ADMIN assignments
reject updates and deletion.

The patch contains schema only. Client/models are generated from the four
SDK-only canonical operations; existing CLI provenance is retained.

Two SDK-only role read projections include namespace/paging typename markers to
distinguish absent objects from malformed non-null namespaces without changing
existing operation documents.
