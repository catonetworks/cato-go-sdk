//nolint:gosec // Tests read files only from paths created beneath t.TempDir.
package gqlops

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

const (
	testOperationName = "foo"
	testOperationKey  = "query.foo"
	testOperationFile = "sources/query.foo.gql"
	testCLICommitSHA  = "0123456789abcdef0123456789abcdef01234567"
	stableSDKName     = "sdkExisting"
)

func TestImportAndCheck(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	sdkRoot := filepath.Join(root, "sdk")
	cliRoot := filepath.Join(root, "cli")
	mustMkdirAll(t, filepath.Join(sdkRoot, "sources"))
	mustMkdirAll(t, filepath.Join(cliRoot, "queryPayloads"))
	mustWrite(t, filepath.Join(sdkRoot, "cato_api.graphqls"), `
		schema { query: Query }
		type Query { existing: String!, added: String! }
	`)
	mustWrite(t, filepath.Join(sdkRoot, "sources", "query.existing.gql"), "query existing { existing }\n")
	mustWrite(t, filepath.Join(cliRoot, "queryPayloads", "query.existing.txt"), "query cliExisting { existing }\n")
	mustWrite(t, filepath.Join(cliRoot, "queryPayloads", "query.added.txt"), "query added { added }\n")

	result, err := Import(ImportConfig{
		CLIRoot:      cliRoot,
		SDKRoot:      sdkRoot,
		CLICommitSHA: testCLICommitSHA,
		Expected:     2,
	})
	if err != nil {
		t.Fatalf("Import() error = %v", err)
	}
	if result != (ImportResult{Canonical: 2, Imported: 1, Mapped: 1, SDKOnly: 0}) {
		t.Fatalf("Import() result = %#v", result)
	}
	replaced, err := os.ReadFile(filepath.Join(sdkRoot, "sources", "query.existing.gql"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(replaced), "query existing") ||
		strings.Contains(string(replaced), "cliExisting") {
		t.Fatalf("matched SDK source did not preserve its stable name:\n%s", replaced)
	}
	if err := Check(CheckConfig{SDKRoot: sdkRoot, Expected: 2}); err != nil {
		t.Fatalf("Check() error = %v", err)
	}
}

func TestImportPreservesStableMappedAPIAndIsIdempotent(t *testing.T) { //nolint:gocyclo // One end-to-end test verifies all preserved API properties and idempotence.
	t.Parallel()

	root := t.TempDir()
	sdkRoot := filepath.Join(root, "sdk")
	cliRoot := filepath.Join(root, "cli")
	mustMkdirAll(t, filepath.Join(sdkRoot, "sources"))
	mustMkdirAll(t, filepath.Join(cliRoot, "queryPayloads"))
	mustWrite(t, filepath.Join(sdkRoot, "cato_api.graphqls"), `
		schema { query: Query }
		type Query { existing(a: String!, z: String!, legacy: String!, added: String!, scope: String): Thing! }
		type Thing { name: String!, extra: String! }
	`)
	sourcePath := filepath.Join(sdkRoot, "sources", "query.existing.gql")
	mustWrite(t, sourcePath, `
		query sdkExisting($z: String!, $a: String!, $legacy: String!, $scope: String) {
			existing(z: $z, a: $a, legacy: $legacy, scope: $scope) {
				stableName: name
			}
		}
	`)
	mustWrite(t, filepath.Join(cliRoot, "queryPayloads", "query.existing.txt"), `
		query cliRenamed($added: String!, $a: String!, $z: String!, $cliValue: String!) {
			existing(a: $a, z: $z, legacy: $cliValue, added: $added) {
				cliName: name
				extra
			}
		}
	`)

	config := ImportConfig{
		CLIRoot:      cliRoot,
		SDKRoot:      sdkRoot,
		CLICommitSHA: testCLICommitSHA,
		Expected:     1,
	}
	if _, err := Import(config); err != nil {
		t.Fatalf("Import() error = %v", err)
	}

	content, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	query, err := parseQueryDocument(sourcePath, content)
	if err != nil {
		t.Fatal(err)
	}
	operation := query.Operations[0]
	if operation.Name != stableSDKName {
		t.Fatalf("operation name = %q; want sdkExisting", operation.Name)
	}
	gotVariables := make([]string, 0, len(operation.VariableDefinitions))
	for _, variable := range operation.VariableDefinitions {
		gotVariables = append(gotVariables, variable.Variable)
	}
	if !reflect.DeepEqual(gotVariables, []string{"z", "a", "legacy", "scope", "added"}) {
		t.Fatalf("variable order = %#v; want z, a, legacy, scope, added", gotVariables)
	}
	if !strings.Contains(string(content), "stableName: name") ||
		strings.Contains(string(content), "cliName: name") ||
		strings.Contains(string(content), "$cliValue") ||
		!strings.Contains(string(content), "scope: $scope") {
		t.Fatalf("stable response alias was not retained:\n%s", content)
	}

	manifestPath := filepath.Join(sdkRoot, "operations", "manifest.json")
	firstManifest, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	manifest, err := loadManifest(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if manifest.CLICommitSHA != testCLICommitSHA {
		t.Fatalf("CLI commit SHA = %q; want %q", manifest.CLICommitSHA, testCLICommitSHA)
	}
	entry := manifest.Operations[0]
	if entry.SDKName != stableSDKName || entry.CLIName != "cliRenamed" {
		t.Fatalf("manifest names = SDK %q, CLI %q", entry.SDKName, entry.CLIName)
	}
	if !reflect.DeepEqual(entry.Variables, []string{"z", "a", "legacy", "scope", "added"}) {
		t.Fatalf("manifest variables = %#v", entry.Variables)
	}

	if _, err := Import(config); err != nil {
		t.Fatalf("second Import() error = %v", err)
	}
	secondContent, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	secondManifest, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(content, secondContent) ||
		!reflect.DeepEqual(firstManifest, secondManifest) {
		t.Fatal("second import changed normalized output")
	}
}

func TestImportAvoidsHandwrittenCollisionAndRetainsSDKOnlySource(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	sdkRoot := filepath.Join(root, "sdk")
	cliRoot := filepath.Join(root, "cli")
	mustMkdirAll(t, filepath.Join(sdkRoot, "sources"))
	mustMkdirAll(t, filepath.Join(cliRoot, "queryPayloads"))
	mustWrite(t, filepath.Join(sdkRoot, "cato_api.graphqls"), `
		schema { query: Query }
		type Query { keep: String!, added(a: String!, z: String!): String! }
	`)
	sdkOnlyPath := filepath.Join(sdkRoot, "sources", "query.keep.gql")
	mustWrite(t, sdkOnlyPath, "query keep { keep }\n")
	helperPath := filepath.Join(sdkRoot, "helper.go")
	helperContent := "package cato_go_sdk\n\nfunc (c *Client) Added() {}\n"
	mustWrite(t, helperPath, helperContent)
	mustWrite(
		t,
		filepath.Join(cliRoot, "queryPayloads", "query.added.txt"),
		"query added($z: String!, $a: String!) { added(z: $z, a: $a) }\n",
	)

	result, err := Import(ImportConfig{
		CLIRoot:      cliRoot,
		SDKRoot:      sdkRoot,
		CLICommitSHA: testCLICommitSHA,
		Expected:     2,
	})
	if err != nil {
		t.Fatalf("Import() error = %v", err)
	}
	if result != (ImportResult{Canonical: 2, Imported: 1, SDKOnly: 1}) {
		t.Fatalf("Import() result = %#v", result)
	}

	imported, err := os.ReadFile(filepath.Join(sdkRoot, "sources", "query.added.gql"))
	if err != nil {
		t.Fatal(err)
	}
	query, err := parseQueryDocument("query.added.gql", imported)
	if err != nil {
		t.Fatal(err)
	}
	operation := query.Operations[0]
	if operation.Name != "queryAdded" ||
		operation.VariableDefinitions[0].Variable != "a" ||
		operation.VariableDefinitions[1].Variable != "z" {
		t.Fatalf("colliding operation was not normalized:\n%s", imported)
	}
	sdkOnly, err := os.ReadFile(sdkOnlyPath)
	if err != nil {
		t.Fatal(err)
	}
	helper, err := os.ReadFile(helperPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(sdkOnly) != "query keep { keep }\n" || string(helper) != helperContent {
		t.Fatal("SDK-only source or handwritten helper changed")
	}
}

func TestImportReplacesMappedSourceInvalidatedByNewSchema(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	sdkRoot := filepath.Join(root, "sdk")
	cliRoot := filepath.Join(root, "cli")
	mustMkdirAll(t, filepath.Join(sdkRoot, "sources"))
	mustMkdirAll(t, filepath.Join(cliRoot, "queryPayloads"))
	mustWrite(t, filepath.Join(sdkRoot, "cato_api.graphqls"), `
		schema { query: Query }
		type Query { user: User! }
		type User { userImportType: String! }
	`)
	mustWrite(
		t,
		filepath.Join(sdkRoot, "sources", "query.user.gql"),
		"query stableUser { user {\nimportType\nuserImportType\n} }\n",
	)
	mustWrite(
		t,
		filepath.Join(cliRoot, "queryPayloads", "query.user.txt"),
		"query cliUser { user {\nimportType\nuserImportType\n} }\n",
	)

	if _, err := Import(ImportConfig{
		CLIRoot:      cliRoot,
		SDKRoot:      sdkRoot,
		CLICommitSHA: testCLICommitSHA,
		Expected:     1,
	}); err != nil {
		t.Fatalf("Import() error = %v", err)
	}
	content, err := os.ReadFile(filepath.Join(sdkRoot, "sources", "query.user.gql"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(content), "\n\t\timportType\n") ||
		!strings.Contains(string(content), "query stableUser") {
		t.Fatalf("invalid mapped source was not safely replaced:\n%s", content)
	}
}

func TestParseOperationRejectsMalformedDocument(t *testing.T) {
	t.Parallel()

	schema, err := gqlparser.LoadSchema(&ast.Source{Input: "type Query { ping: String! }"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := parseOperation(schema, "bad.gql", []byte("query bad {")); err == nil {
		t.Fatal("parseOperation() accepted malformed document")
	}
}

func TestValidateUniqueDocuments(t *testing.T) {
	t.Parallel()

	documents := []document{
		{key: testOperationKey, name: testOperationName, relative: testOperationFile},
		{key: "QUERY.FOO", name: "bar", relative: "sources/query.bar.gql"},
		{key: "query.baz", name: testOperationName, relative: "sources/query.baz.gql"},
	}
	err := validateUniqueDocuments(documents)
	if err == nil {
		t.Fatal("validateUniqueDocuments() accepted duplicate key and name")
	}
	if !strings.Contains(err.Error(), "semantic key") || !strings.Contains(err.Error(), "operation name") {
		t.Fatalf("validateUniqueDocuments() error = %v", err)
	}
}

func TestCurateKnownCLIProblems(t *testing.T) { //nolint:gocyclo // One table-free test keeps each targeted repair readable.
	t.Parallel()

	headers := []byte("plainHeaders\nsecretHeaders\nplainHeaders {\n  name\n  value\n}\n")
	curated, err := curateCLIContent("query.notification.txt", headers)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(curated), "plainHeaders {") != 2 || !strings.Contains(string(curated), "secretHeaders {") {
		t.Fatalf("header curation failed:\n%s", curated)
	}

	devices := []byte("query devices { devices {\n{\n duplicate { id }\n}\nkeep\n{\n duplicate { id }\n}\n} }\n")
	curated, err = curateCLIContent("query.devices.txt", devices)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(curated), "duplicate") || !strings.Contains(string(curated), "keep") {
		t.Fatalf("devices curation failed:\n%s", curated)
	}

	fixedDisable := []byte("mutation disable { disableAccount ( accountId:$accountId ) { id } }\n")
	curated, err = curateCLIContent("mutation.accountManagement.disableAccount.txt", fixedDisable)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(curated, fixedDisable) {
		t.Fatalf("fixed account mutation changed:\n%s", curated)
	}

	fixedDevices := []byte("query devices { devices { keep } }\n")
	curated, err = curateCLIContent("query.devices.txt", fixedDevices)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(curated, fixedDevices) {
		t.Fatalf("fixed devices query changed:\n%s", curated)
	}

	user, err := curateCLIContent(
		"query.user.txt",
		[]byte("query user { user {\nimportType\nuserImportType\n} }\n"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(user), "\nimportType\n") ||
		!strings.Contains(string(user), "userImportType") {
		t.Fatalf("retired user field curation failed:\n%s", user)
	}

	socketLAN := []byte(`
		query policySocketLanPolicy {
			policy {
				socketLan {
					policy {
						sections {
							section {
								id
								name
							}
						}
					}
				}
			}
		}
	`)
	curated, err = curateCLIContent("query.policy.socketLan.policy.txt", socketLAN)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(curated), "subPolicyId") != 1 {
		t.Fatalf("socket LAN section owner curation failed:\n%s", curated)
	}
}

func TestConfinedPath(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	if err := confinedPath(root, filepath.Join(root, "sources", "query.gql")); err != nil {
		t.Fatalf("confinedPath() rejected child: %v", err)
	}
	if err := confinedPath(root, filepath.Join(root, "..", "outside.gql")); err == nil {
		t.Fatal("confinedPath() accepted escaping path")
	}
}

func TestManifestRoundTrip(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "operations", "manifest.json")
	want := &Manifest{Version: manifestVersion, Operations: []ManifestEntry{{
		Key:            testOperationKey,
		Kind:           operationKindQuery,
		SDKFile:        testOperationFile,
		SDKName:        testOperationName,
		DocumentSHA256: "hash",
		Status:         statusSDKOnly,
	}}}
	if err := writeManifest(path, want); err != nil {
		t.Fatal(err)
	}
	got, err := loadManifest(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("loadManifest() = %#v", got)
	}
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o750); err != nil {
		t.Fatal(err)
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
