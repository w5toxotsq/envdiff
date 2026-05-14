package parser_test

import (
	"testing"

	"github.com/envdiff/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFile_ValueWithEquals(t *testing.T) {
	// Values that contain '=' should be preserved correctly
	path := writeEnvFile(t, "JDBC_URL=jdbc:postgresql://host/db?ssl=true\n")
	env, err := parser.ParseFile(path)
	require.NoError(t, err)
	assert.Equal(t, "jdbc:postgresql://host/db?ssl=true", env["JDBC_URL"])
}

func TestParseFile_WhitespaceAroundKeyValue(t *testing.T) {
	path := writeEnvFile(t, "  APP_ENV  =  production  \n")
	env, err := parser.ParseFile(path)
	require.NoError(t, err)
	assert.Equal(t, "production", env["APP_ENV"])
}

func TestParseFile_MultipleEntries(t *testing.T) {
	content := "HOST=localhost\nPORT=8080\nDEBUG=true\nSECRET=abc123\n"
	path := writeEnvFile(t, content)
	env, err := parser.ParseFile(path)
	require.NoError(t, err)
	assert.Len(t, env, 4)
	assert.Equal(t, "localhost", env["HOST"])
	assert.Equal(t, "8080", env["PORT"])
	assert.Equal(t, "true", env["DEBUG"])
	assert.Equal(t, "abc123", env["SECRET"])
}

func TestEnvMap_IsMap(t *testing.T) {
	path := writeEnvFile(t, "KEY=val\n")
	env, err := parser.ParseFile(path)
	require.NoError(t, err)
	// Ensure EnvMap behaves as a standard map
	var _ map[string]string = env
	assert.NotNil(t, env)
}
