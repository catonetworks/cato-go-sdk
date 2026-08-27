package gqlops

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

const (
	testOperationName = "foo"
	testOperationKey  = "query.foo"
	testOperationFile = "sources/query.foo.gql"
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

	result, err := Import(ImportConfig{CLIRoot: cliRoot, SDKRoot: sdkRoot, Expected: 2})
	if err != nil {
		t.Fatalf("Import() error = %v", err)
	}
	if result != (ImportResult{Canonical: 2, Imported: 1, Mapped: 1, SDKOnly: 0}) {
		t.Fatalf("Import() result = %#v", result)
	}
	if err := Check(CheckConfig{SDKRoot: sdkRoot, Expected: 2}); err != nil {
		t.Fatalf("Check() error = %v", err)
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

func TestCurateKnownCLIProblems(t *testing.T) {
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

	xdr, err := curateCLIContent("mutation.xdr.analystFeedback.txt", []byte("invalid generated selection"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(xdr), "story {\n        id") {
		t.Fatalf("XDR curation failed:\n%s", xdr)
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
	if len(got.Operations) != 1 || got.Operations[0] != want.Operations[0] {
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
