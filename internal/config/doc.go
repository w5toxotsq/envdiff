// Package config handles runtime configuration for envdiff.
//
// Configuration can originate from two sources, applied in priority order:
//
//  1. CLI flags parsed by cmd/envdiff/flags.go (highest priority).
//  2. An optional YAML config file (e.g. envdiff.yaml) in the working directory.
//
// Use DefaultConfig to obtain a Config with sensible defaults, then
// optionally call LoadFile to merge values from a YAML file. Fields
// already populated (non-zero) in the Config are never overwritten by
// LoadFile, ensuring CLI flags always win.
//
// Example YAML config file:
//
//	format: json
//	show_values: true
//	prefix: APP_
//	files:
//	  - .env.staging
//	  - .env.production
package config
