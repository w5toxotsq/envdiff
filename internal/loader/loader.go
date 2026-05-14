// Package loader provides utilities for loading and validating .env files
// from the filesystem before passing them to the parser and diff engine.
package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/envdiff/internal/parser"
)

// ErrFileNotFound is returned when a specified .env file does not exist.
type ErrFileNotFound struct {
	Path string
}

func (e *ErrFileNotFound) Error() string {
	return fmt.Sprintf("env file not found: %s", e.Path)
}

// ErrInvalidExtension is returned when a file does not have a recognised extension.
type ErrInvalidExtension struct {
	Path string
}

func (e *ErrInvalidExtension) Error() string {
	return fmt.Sprintf("file does not appear to be an env file: %s", e.Path)
}

// Options controls optional behaviour of Load.
type Options struct {
	// SkipExtensionCheck disables the .env extension validation.
	SkipExtensionCheck bool
}

// Load reads an env file from disk, optionally validates its extension, and
// returns the parsed key/value map.
func Load(path string, opts Options) (parser.EnvMap, error) {
	if !opts.SkipExtensionCheck {
		if err := validateExtension(path); err != nil {
			return nil, err
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, &ErrFileNotFound{Path: path}
	}

	env, err := parser.ParseFile(path)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	return env, nil
}

// validateExtension returns an error when path does not look like an env file.
func validateExtension(path string) error {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	// Accept .env, .env.* (e.g. .env.production) or files named exactly ".env".
	if ext == ".env" || base == ".env" {
		return nil
	}
	// Accept files whose name starts with ".env" (e.g. .env.local).
	if len(base) > 4 && base[:4] == ".env" {
		return nil
	}
	return &ErrInvalidExtension{Path: path}
}
