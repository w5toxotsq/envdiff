package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envdiff/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)
	return path
}

func TestParseFile_BasicKeyValue(t *testing.T) {
	path := writeEnvFile(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	env, err := parser.ParseFile(path)
	require.NoError(t, err)
	assert.Equal(t, "localhost", env["DB_HOST"])
	assert.Equal(t, "5432", env["DB_PORT"])
}

func TestParseFile_IgnoresComments(t *testing.T) {
	path := writeEnvFile(t, "# This is a comment\nAPP_NAME=envdiff\n")
	env, err := parser.ParseFile(path)
	require.NoError(t, err)
	assert.Len(t, env, 1)
	assert.Equal(t, "envdiff", env["APP_NAME"])
}

func TestParseFile_IgnoresEmptyLines(t *testing.T) {
	path := writeEnvFile(t, "\nKEY=value\n\n")
	env, err := parser.ParseFile(path)
	require.NoError(t, err)
	assert.Len(t, env, 1)
}

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeEnvFile(t, `SECRET="my secret value"` + "\n" + `TOKEN='abc123'` + "\n")
	env, err := parser.ParseFile(path)
	require.NoError(t, err)
	assert.Equal(t, "my secret value", env["SECRET"])
	assert.Equal(t, "abc123", env["TOKEN"])
}

func TestParseFile_InvalidLine(t *testing.T) {
	path := writeEnvFile(t, "INVALID_LINE_NO_EQUALS\n")
	_, err := parser.ParseFile(path)
	assert.Error(t, err)
}

func TestParseFile_FileNotFound(t *testing.T) {
	_, err := parser.ParseFile("/nonexistent/.env")
	assert.Error(t, err)
}

func TestParseFile_EmptyKey(t *testing.T) {
	path := writeEnvFile(t, "=value\n")
	_, err := parser.ParseFile(path)
	assert.Error(t, err)
}
