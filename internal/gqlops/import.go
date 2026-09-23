package gqlops

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

const (
	statusImported          = "imported"
	statusMapped            = "mapped"
	statusSDKOnly           = "sdk_only"
	operationNameMatchParts = 2
	fileNameSplitParts      = 2
)

var operationNamePattern = regexp.MustCompile(`(?m)^\s*(?:query|mutation)\s+([_A-Za-z][_0-9A-Za-z]*)`)

// ImportConfig configures migration from cato-cli into cato-go-sdk.
type ImportConfig struct {
	CLIRoot      string
	SDKRoot      string
	CLICommitSHA string
	Expected     int
	DryRun       bool
}

// ImportResult summarizes an operation migration.
type ImportResult struct {
	Canonical int
	Imported  int
	Mapped    int
	SDKOnly   int
}

type candidate struct {
	document    document
	destination string
}

// Import migrates validated CLI-only documents and writes the canonical manifest.
func Import(config ImportConfig) (ImportResult, error) {
	if err := validateImportConfig(config); err != nil {
		return ImportResult{}, err
	}

	schema, err := loadSchema(config.SDKRoot)
	if err != nil {
		return ImportResult{}, err
	}
	sdkDocuments, err := loadSDKDocuments(config.SDKRoot, schema)
	if err != nil {
		return ImportResult{}, err
	}
	if err := validateUniqueDocuments(sdkDocuments); err != nil {
		return ImportResult{}, fmt.Errorf("validate existing SDK operations: %w", err)
	}
	reservedGoNames, err := handwrittenGoNames(config.SDKRoot)
	if err != nil {
		return ImportResult{}, err
	}

	cliPaths, err := cliOperationPaths(config.CLIRoot)
	if err != nil {
		return ImportResult{}, err
	}

	manifest, candidates, result, err := buildImport(
		schema,
		config,
		sdkDocuments,
		cliPaths,
		reservedGoNames,
	)
	if err != nil {
		return ImportResult{}, err
	}
	manifestPath := filepath.Join(config.SDKRoot, "operations", "manifest.json")
	if err := preserveImportedStatuses(manifestPath, manifest); err != nil {
		return ImportResult{}, err
	}
	if config.Expected > 0 && result.Canonical != config.Expected {
		return ImportResult{}, fmt.Errorf("canonical operation count is %d; want %d", result.Canonical, config.Expected)
	}
	if config.DryRun {
		return result, nil
	}

	if err := writeCandidates(config.SDKRoot, candidates); err != nil {
		return ImportResult{}, err
	}
	if err := writeManifest(manifestPath, manifest); err != nil {
		return ImportResult{}, err
	}

	return result, nil
}

func preserveImportedStatuses(path string, manifest *Manifest) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return fmt.Errorf("inspect existing manifest %q: %w", path, err)
	}

	previous, err := loadManifest(path)
	if err != nil {
		return err
	}
	previousStatus := make(map[string]string, len(previous.Operations))
	for _, entry := range previous.Operations {
		previousStatus[entry.SDKFile] = entry.Status
	}
	for index := range manifest.Operations {
		entry := &manifest.Operations[index]
		if previousStatus[entry.SDKFile] == statusImported && entry.CLIFile != "" {
			entry.Status = statusImported
		}
	}

	return nil
}

func buildImport( //nolint:funlen // Keeping candidate collection in one pass guarantees validation completes before any write.
	schema *ast.Schema,
	config ImportConfig,
	sdkDocuments []document,
	cliPaths []string,
	reservedGoNames map[string]struct{},
) (*Manifest, []candidate, ImportResult, error) {
	records := make([]ManifestEntry, 0, len(sdkDocuments)+len(cliPaths))
	documents := append([]document(nil), sdkDocuments...)
	indexByKey := make(map[string]int, len(sdkDocuments))
	indexByName := make(map[string]int, len(sdkDocuments))

	for _, sdkDocument := range sdkDocuments {
		indexByKey[normalizedKey(sdkDocument.key)] = len(records)
		indexByName[sdkDocument.name] = len(records)
		reservedGoNames[goIdentifier(sdkDocument.name)] = struct{}{}
		records = append(records, entryFromSDK(sdkDocument))
	}

	var candidates []candidate
	var validationErrors []error
	result := ImportResult{SDKOnly: len(sdkDocuments)}

	for _, cliPath := range cliPaths {
		cliContent, err := readRegularFile(cliPath, maxOperationBytes)
		if err != nil {
			validationErrors = append(validationErrors, err)
			continue
		}
		cliKind := strings.SplitN(filepath.Base(cliPath), ".", fileNameSplitParts)[0]
		key, err := semanticKey(cliPath, cliKind)
		if err != nil {
			validationErrors = append(validationErrors, err)
			continue
		}

		curatedContent, err := curateCLIContent(cliPath, cliContent)
		if err != nil {
			validationErrors = append(validationErrors, err)
			continue
		}
		parsed, err := parseOperation(schema, cliPath, curatedContent)
		if err != nil {
			validationErrors = append(validationErrors, err)
			continue
		}
		if parsed.kind != cliKind {
			validationErrors = append(validationErrors, fmt.Errorf(
				"CLI operation %q contains %s; filename declares %s",
				cliPath,
				parsed.kind,
				cliKind,
			))
			continue
		}

		cliRelative := filepath.ToSlash(filepath.Join("queryPayloads", filepath.Base(cliPath)))
		existingIndex, exists := indexByKey[normalizedKey(key)]
		if exists {
			existing := documents[existingIndex]
			if parsed.kind != existing.kind {
				validationErrors = append(validationErrors, fmt.Errorf(
					"CLI operation %q kind %q does not match SDK operation %q kind %q",
					cliPath,
					parsed.kind,
					existing.relative,
					existing.kind,
				))
				continue
			}

			cliName := parsed.name
			normalizedContent, normalizeErr := normalizeMappedOperation(
				schema,
				cliPath,
				existing.content,
				curatedContent,
				existing.name,
			)
			if normalizeErr != nil {
				validationErrors = append(validationErrors, normalizeErr)
				continue
			}
			parsed, err = parseOperation(schema, cliPath, normalizedContent)
			if err != nil {
				validationErrors = append(validationErrors, err)
				continue
			}
			parsed.key = existing.key
			parsed.path = existing.path
			parsed.relative = existing.relative
			parsed.content = normalizedContent
			parsed.hash = contentHash(normalizedContent)
			if err := mapCLIEntry(&records[existingIndex], cliRelative, cliContent, cliName, parsed); err != nil {
				validationErrors = append(validationErrors, err)
			} else {
				documents[existingIndex] = parsed
				indexByKey[normalizedKey(key)] = existingIndex
				candidates = append(candidates, candidate{document: parsed, destination: existing.path})
				result.Mapped++
				result.SDKOnly--
			}
			continue
		}

		destination := filepath.Join(config.SDKRoot, "sources", strings.TrimSuffix(filepath.Base(cliPath), ".txt")+".gql")
		if _, err := os.Lstat(destination); err == nil {
			validationErrors = append(validationErrors, fmt.Errorf("refuse to overwrite unindexed SDK operation %q", destination))
			continue
		} else if !errors.Is(err, os.ErrNotExist) {
			validationErrors = append(validationErrors, fmt.Errorf("inspect destination %q: %w", destination, err))
			continue
		}

		cliName := parsed.name
		stableName := stableOperationName(cliName, key, indexByName, reservedGoNames)
		normalizedContent, normalizeErr := normalizeNewOperation(
			schema,
			cliPath,
			curatedContent,
			stableName,
		)
		if normalizeErr != nil {
			validationErrors = append(validationErrors, normalizeErr)
			continue
		}
		parsed, err = parseOperation(schema, cliPath, normalizedContent)
		if err != nil {
			validationErrors = append(validationErrors, err)
			continue
		}
		parsed.key = key
		parsed.path = destination
		parsed.relative = filepath.ToSlash(filepath.Join("sources", filepath.Base(destination)))
		parsed.content = normalizedContent
		parsed.hash = contentHash(normalizedContent)
		documents = append(documents, parsed)
		indexByKey[normalizedKey(key)] = len(records)
		indexByName[parsed.name] = len(records)
		reservedGoNames[goIdentifier(parsed.name)] = struct{}{}
		records = append(records, ManifestEntry{
			Key:            key,
			Kind:           parsed.kind,
			SDKFile:        parsed.relative,
			SDKName:        parsed.name,
			Variables:      parsed.variables,
			DocumentSHA256: parsed.hash,
			CLIFile:        cliRelative,
			CLIName:        cliName,
			CLISHA256:      contentHash(cliContent),
			Status:         statusImported,
		})
		candidates = append(candidates, candidate{document: parsed, destination: destination})
		result.Imported++
	}

	if err := errors.Join(validationErrors...); err != nil {
		return nil, nil, ImportResult{}, err
	}
	if err := validateUniqueDocuments(documents); err != nil {
		return nil, nil, ImportResult{}, fmt.Errorf("validate canonical operations: %w", err)
	}

	sort.Slice(records, func(left, right int) bool {
		return normalizedKey(records[left].Key) < normalizedKey(records[right].Key)
	})
	result.Canonical = len(records)

	return &Manifest{
		Version:      manifestVersion,
		CLICommitSHA: config.CLICommitSHA,
		Operations:   records,
	}, candidates, result, nil
}

func entryFromSDK(sdkDocument document) ManifestEntry {
	return ManifestEntry{
		Key:            sdkDocument.key,
		Kind:           sdkDocument.kind,
		SDKFile:        sdkDocument.relative,
		SDKName:        sdkDocument.name,
		Variables:      sdkDocument.variables,
		DocumentSHA256: sdkDocument.hash,
		Status:         statusSDKOnly,
	}
}

func updateMappedEntry(
	entry *ManifestEntry,
	cliRelative string,
	cliContent []byte,
	cliName string,
	document document,
) {
	entry.Key = document.key
	entry.Kind = document.kind
	entry.SDKName = document.name
	entry.Variables = document.variables
	entry.DocumentSHA256 = document.hash
	entry.CLIFile = cliRelative
	entry.CLIName = cliName
	entry.CLISHA256 = contentHash(cliContent)
	entry.Status = statusMapped
}

func mapCLIEntry(
	entry *ManifestEntry,
	cliRelative string,
	cliContent []byte,
	cliName string,
	document document,
) error {
	if entry.Status != statusSDKOnly {
		return fmt.Errorf(
			"CLI operations %q and %q both map to %q",
			entry.CLIFile,
			cliRelative,
			entry.SDKFile,
		)
	}
	updateMappedEntry(entry, cliRelative, cliContent, cliName, document)
	return nil
}

func operationName(content []byte) string {
	match := operationNamePattern.FindSubmatch(content)
	if len(match) != operationNameMatchParts {
		return ""
	}
	return string(match[1])
}

func cliOperationPaths(cliRoot string) ([]string, error) {
	payloadRoot := filepath.Join(cliRoot, "queryPayloads")
	if err := confinedPath(cliRoot, payloadRoot); err != nil {
		return nil, err
	}

	var paths []string
	for _, pattern := range []string{"query.*.txt", "mutation.*.txt"} {
		matches, err := filepath.Glob(filepath.Join(payloadRoot, pattern))
		if err != nil {
			return nil, fmt.Errorf("scan CLI operations with %q: %w", pattern, err)
		}
		paths = append(paths, matches...)
	}
	sort.Strings(paths)

	return paths, nil
}

func writeCandidates(sdkRoot string, candidates []candidate) error {
	stageRoot, err := os.MkdirTemp(sdkRoot, ".gqlops-stage-")
	if err != nil {
		return fmt.Errorf("create operation staging directory: %w", err)
	}
	defer func() {
		_ = os.RemoveAll(stageRoot)
	}()

	for _, current := range candidates {
		stagedPath := filepath.Join(stageRoot, filepath.Base(current.destination))
		// #nosec G306 -- GraphQL documents contain no secrets and are committed source files.
		if err := os.WriteFile(stagedPath, current.document.content, 0o644); err != nil {
			return fmt.Errorf("stage operation %q: %w", current.destination, err)
		}
	}
	for _, current := range candidates {
		stagedPath := filepath.Join(stageRoot, filepath.Base(current.destination))
		if err := os.Rename(stagedPath, current.destination); err != nil {
			return fmt.Errorf("install operation %q: %w", current.destination, err)
		}
	}

	return nil
}

func validateImportConfig(config ImportConfig) error {
	if config.CLIRoot == "" {
		return errors.New("CLI root is required")
	}
	if config.SDKRoot == "" {
		return errors.New("SDK root is required")
	}
	if !commitSHAPattern.MatchString(config.CLICommitSHA) {
		return fmt.Errorf("CLI commit SHA %q is not a full commit SHA", config.CLICommitSHA)
	}
	if err := confinedPath(config.CLIRoot, filepath.Join(config.CLIRoot, "queryPayloads")); err != nil {
		return err
	}
	if err := confinedPath(config.SDKRoot, filepath.Join(config.SDKRoot, "sources")); err != nil {
		return err
	}
	return nil
}
