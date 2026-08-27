package gqlops

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const manifestVersion = 1
const maxManifestBytes = 8 << 20

// Manifest describes every canonical GraphQL operation in the SDK.
type Manifest struct {
	Version    int             `json:"version"`
	Operations []ManifestEntry `json:"operations"`
}

// ManifestEntry links a canonical SDK document to its former CLI source.
type ManifestEntry struct {
	Key            string `json:"key"`
	Kind           string `json:"kind"`
	SDKFile        string `json:"sdk_file"`
	SDKName        string `json:"sdk_name"`
	DocumentSHA256 string `json:"document_sha256"`
	CLIFile        string `json:"cli_file,omitempty"`
	CLIName        string `json:"cli_name,omitempty"`
	CLISHA256      string `json:"cli_sha256,omitempty"`
	Status         string `json:"status"`
}

func loadManifest(path string) (*Manifest, error) {
	data, err := readRegularFile(path, maxManifestBytes)
	if err != nil {
		return nil, fmt.Errorf("read manifest %q: %w", path, err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("decode manifest %q: %w", path, err)
	}
	if manifest.Version != manifestVersion {
		return nil, fmt.Errorf("manifest %q has version %d; want %d", path, manifest.Version, manifestVersion)
	}

	return &manifest, nil
}

func writeManifest(path string, manifest *Manifest) error {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("encode manifest: %w", err)
	}
	data = append(data, '\n')

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return fmt.Errorf("create manifest directory: %w", err)
	}

	temp, err := os.CreateTemp(filepath.Dir(path), ".manifest-*.json")
	if err != nil {
		return fmt.Errorf("create temporary manifest: %w", err)
	}
	tempName := temp.Name()
	defer func() {
		_ = os.Remove(tempName)
	}()

	if err := temp.Chmod(0o644); err != nil {
		_ = temp.Close()
		return fmt.Errorf("set temporary manifest permissions: %w", err)
	}
	if _, err := temp.Write(data); err != nil {
		_ = temp.Close()
		return fmt.Errorf("write temporary manifest: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close temporary manifest: %w", err)
	}
	if err := os.Rename(tempName, path); err != nil {
		return fmt.Errorf("replace manifest %q: %w", path, err)
	}

	return nil
}
