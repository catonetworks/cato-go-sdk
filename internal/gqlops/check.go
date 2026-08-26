package gqlops

import (
	"errors"
	"fmt"
	"path/filepath"
)

// CheckConfig configures canonical SDK operation validation.
type CheckConfig struct {
	SDKRoot  string
	Expected int
}

// Check validates all SDK documents and their manifest without modifying files.
func Check(config CheckConfig) error {
	if config.SDKRoot == "" {
		return errors.New("SDK root is required")
	}

	schema, err := loadSchema(config.SDKRoot)
	if err != nil {
		return err
	}
	documents, err := loadSDKDocuments(config.SDKRoot, schema)
	if err != nil {
		return err
	}
	if err := validateUniqueDocuments(documents); err != nil {
		return err
	}
	if config.Expected > 0 && len(documents) != config.Expected {
		return fmt.Errorf("SDK operation count is %d; want %d", len(documents), config.Expected)
	}

	manifestPath := filepath.Join(config.SDKRoot, "operations", "manifest.json")
	manifest, err := loadManifest(manifestPath)
	if err != nil {
		return err
	}
	if len(manifest.Operations) != len(documents) {
		return fmt.Errorf(
			"manifest contains %d operations; SDK contains %d",
			len(manifest.Operations),
			len(documents),
		)
	}

	return compareManifest(config.SDKRoot, manifest, documents)
}

func compareManifest(sdkRoot string, manifest *Manifest, documents []document) error {
	documentByFile := make(map[string]document, len(documents))
	for _, current := range documents {
		documentByFile[current.relative] = current
	}

	seenFiles := make(map[string]struct{}, len(manifest.Operations))
	seenKeys := make(map[string]string, len(manifest.Operations))
	var validationErrors []error

	for _, entry := range manifest.Operations {
		if err := validateManifestEntry(sdkRoot, entry); err != nil {
			validationErrors = append(validationErrors, err)
			continue
		}
		if _, exists := seenFiles[entry.SDKFile]; exists {
			validationErrors = append(validationErrors, fmt.Errorf("manifest repeats SDK file %q", entry.SDKFile))
			continue
		}
		seenFiles[entry.SDKFile] = struct{}{}

		key := normalizedKey(entry.Key)
		if previous, exists := seenKeys[key]; exists {
			validationErrors = append(validationErrors, fmt.Errorf("manifest key %q repeats %q", entry.Key, previous))
			continue
		}
		seenKeys[key] = entry.SDKFile

		current, exists := documentByFile[entry.SDKFile]
		if !exists {
			validationErrors = append(validationErrors, fmt.Errorf("manifest SDK file %q does not exist", entry.SDKFile))
			continue
		}
		if entry.Key != current.key {
			validationErrors = append(validationErrors, fmt.Errorf("%s key is %q; want %q", entry.SDKFile, entry.Key, current.key))
		}
		if entry.Kind != current.kind {
			validationErrors = append(validationErrors, fmt.Errorf("%s kind is %q; want %q", entry.SDKFile, entry.Kind, current.kind))
		}
		if entry.SDKName != current.name {
			validationErrors = append(validationErrors, fmt.Errorf("%s name is %q; want %q", entry.SDKFile, entry.SDKName, current.name))
		}
		if entry.DocumentSHA256 != current.hash {
			validationErrors = append(validationErrors, fmt.Errorf("%s content hash does not match manifest", entry.SDKFile))
		}
	}

	return errors.Join(validationErrors...)
}

func validateManifestEntry(sdkRoot string, entry ManifestEntry) error {
	if entry.Key == "" || entry.Kind == "" || entry.SDKFile == "" || entry.SDKName == "" || entry.DocumentSHA256 == "" {
		return fmt.Errorf("manifest entry for %q has missing required fields", entry.SDKFile)
	}
	if err := validateManifestStatus(entry); err != nil {
		return err
	}

	expectedPath := filepath.Join(sdkRoot, filepath.FromSlash(entry.SDKFile))
	if err := confinedPath(filepath.Join(sdkRoot, "sources"), expectedPath); err != nil {
		return fmt.Errorf("manifest entry %q: %w", entry.SDKFile, err)
	}
	if filepath.Dir(entry.SDKFile) != "sources" {
		return fmt.Errorf("manifest entry %q must reference a direct child of sources", entry.SDKFile)
	}

	return nil
}

func validateManifestStatus(entry ManifestEntry) error {
	switch entry.Status {
	case statusImported:
		if entry.CLIFile == "" || entry.CLIName == "" || entry.CLISHA256 == "" {
			return fmt.Errorf("imported manifest entry %q has incomplete CLI metadata", entry.SDKFile)
		}
	case statusMapped:
		if entry.CLIFile == "" || entry.CLISHA256 == "" {
			return fmt.Errorf("mapped manifest entry %q has incomplete CLI metadata", entry.SDKFile)
		}
	case statusSDKOnly:
		if entry.CLIFile != "" || entry.CLIName != "" || entry.CLISHA256 != "" {
			return fmt.Errorf("SDK-only manifest entry %q unexpectedly has CLI metadata", entry.SDKFile)
		}
	default:
		return fmt.Errorf("manifest entry %q has invalid status %q", entry.SDKFile, entry.Status)
	}

	return nil
}
