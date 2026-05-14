// Package config provides configuration loading for envdiff,
// supporting both CLI flags and optional config file (envdiff.yaml).
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the full runtime configuration for an envdiff run.
type Config struct {
	// Files lists the .env file paths to compare.
	Files []string `yaml:"files"`

	// Prefix restricts comparison to keys with this prefix.
	Prefix string `yaml:"prefix"`

	// OnlyMissing limits output to missing keys only.
	OnlyMissing bool `yaml:"only_missing"`

	// OnlyMismatched limits output to mismatched keys only.
	OnlyMismatched bool `yaml:"only_mismatched"`

	// Format selects the output format: "text" or "json".
	Format string `yaml:"format"`

	// ShowValues controls whether values are printed in the report.
	ShowValues bool `yaml:"show_values"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Format:     "text",
		ShowValues: false,
	}
}

// LoadFile reads a YAML config file from path and merges it into dst.
// Fields already set (non-zero) in dst are not overwritten.
func LoadFile(path string, dst *Config) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var file Config
	if err := yaml.NewDecoder(f).Decode(&file); err != nil {
		return err
	}

	merge(dst, file)
	return nil
}

// merge copies non-zero fields from src into dst only when dst's field is zero.
func merge(dst *Config, src Config) {
	if dst.Prefix == "" && src.Prefix != "" {
		dst.Prefix = src.Prefix
	}
	if dst.Format == "" && src.Format != "" {
		dst.Format = src.Format
	}
	if len(dst.Files) == 0 && len(src.Files) > 0 {
		dst.Files = src.Files
	}
	if !dst.OnlyMissing {
		dst.OnlyMissing = src.OnlyMissing
	}
	if !dst.OnlyMismatched {
		dst.OnlyMismatched = src.OnlyMismatched
	}
	if !dst.ShowValues {
		dst.ShowValues = src.ShowValues
	}
}
