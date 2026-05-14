package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents the key-value pairs parsed from a .env file.
type EnvMap map[string]string

// ParseFile reads a .env file and returns an EnvMap of its contents.
// Lines starting with '#' are treated as comments and ignored.
// Empty lines are also ignored.
func ParseFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", path, err)
	}
	defer f.Close()

	env := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d in %q: %w", lineNum, path, err)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning file %q: %w", path, err)
	}

	return env, nil
}

// parseLine splits a single line into a key and value.
// Supports KEY=VALUE and KEY="VALUE" formats.
func parseLine(line string) (string, string, error) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid format %q: expected KEY=VALUE", line)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	if key == "" {
		return "", "", fmt.Errorf("empty key in line %q", line)
	}

	// Strip surrounding quotes if present
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'') {
			value = value[1 : len(value)-1]
		}
	}

	return key, value, nil
}
