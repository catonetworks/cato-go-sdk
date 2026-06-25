---
name: sdk-schema-update
description: >-
  Regenerate the Cato Go SDK from the latest GraphQL schema. Handles fetching
  the upstream schema, applying patches, running gqlgenc codegen, resolving new
  custom scalar types, and verifying the build. Use when the user asks to update
  the schema, run codegen, regenerate the SDK, sync the GraphQL schema, or
  mentions schema-update / gqlgenc / cato_api.graphqls.
---

# SDK Schema Update

## Workflow

Run steps in order. Stop and report if any step fails.

### 1. Fetch & normalize schema

```bash
make schema-update
```

Archives the current `cato_api.graphqls` to `archives/` and replaces it with the normalized upstream schema.

### 2. Triage and update patches

For **each** `.patch` file in `schema-patches/`, perform the following checks and actions before running `make apply-patches`.

#### 2a. Update stale patch context

If the patch fails to apply (context lines no longer match), the agent must regenerate the patch:

1. Read the patch to extract the **intent** — what lines are being added and roughly where.
2. Locate the correct position in the updated `cato_api.graphqls` by searching for stable surrounding identifiers (type names, field names, enum values).
3. Produce a new unified diff for the same change against the current schema using `git diff --no-index` or by writing the updated patch manually.
4. Overwrite the patch file with the updated diff.
5. Verify: `git apply --check schema-patches/<name>.patch`

Notify the user of each patch that was regenerated: `"Updated context in <name>.patch."`

#### 2b. Check if already in base schema

Extract all lines added by the patch (lines starting with `+`, excluding the `+++ b/` header). Check whether those lines already exist verbatim in `cato_api.graphqls`.

```bash
# Quick check: see if patch applies in reverse (i.e. is already applied)
git apply --reverse --check schema-patches/<name>.patch 2>/dev/null
```

If the patch content **is already present** in the base schema:
- Notify the user: `"Patch <name>.patch is now included in the upstream schema — deleting."`
- Delete the patch file.
- Skip to the next patch.

#### 2c. Check if patch applies cleanly

```bash
git apply --check schema-patches/<name>.patch 2>&1
```

If it applies cleanly → no action needed, move on.

#### 2d. Run apply-patches

Once all patches are triaged:

```bash
make apply-patches
```

### 3. Run codegen

```bash
make generate
```

Regenerates `client.go` and `models/models.go` via `go tool gqlgenc`.

**If codegen fails with an unknown scalar type**, see [Resolving new custom scalars](#resolving-new-custom-scalars) below.

### 4. Verify build

```bash
go build ./...
```

Fix any compile errors before proceeding.

### 5. Commit

Stage and commit:
- `cato_api.graphqls` (updated schema)
- `archives/cato_api-<date>.graphqls` (dated backup created by `make schema-update`)
- `client.go` and `models/models.go` (regenerated)
- any deleted patch files from `schema-patches/`
- any new files in `scalars/` or `sources/`
- `.gqlgenc.yml` if you added a scalar mapping

---

## Resolving new custom scalars

When `make generate` fails because a scalar type is unmapped:

1. **Create** `scalars/<scalar_name_snake_case>.go` modelled after an existing scalar, e.g. `scalars/port.go`.
2. **Register** the mapping in `.gqlgenc.yml` under `models:`:

```yaml
models:
  MyNewScalar:
    model: github.com/catonetworks/cato-go-sdk/scalars.MyNewScalar
```

3. Re-run `make generate`.

---

## Key files

| File / Dir | Purpose |
|---|---|
| `cato_api.graphqls` | GraphQL schema (source of truth) |
| `.gqlgenc.yml` | Codegen config (scalar mappings, output paths) |
| `client.go` | Generated GraphQL client |
| `models/models.go` | Generated Go types |
| `scalars/` | Hand-written scalar type implementations |
| `schema-patches/` | Patches applied on top of the upstream schema |
| `sources/*.gql` | GraphQL operation files fed to gqlgenc |
| `archives/` | Dated backups of previous schemas |

---

## Schema URL

Default: `https://system.cc.catonetworks.com/api/schema?with_undocumented=true`

Override with `SCHEMA_CURL_URL=<url> make schema-update`.
