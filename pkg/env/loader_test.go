//go:build unit

package env

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigPath(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Mock current working directory
	os.Chdir(tempDir)

	// Test case 1: Directory ends with "config"
	path, err := GetConfigPath("gateway")
	assert.NoError(t, err, "GetConfigPath should not return an error")
	expectedPath := filepath.Join(tempDir, "gateway/config/config.staging.yaml")
	assert.Equal(t, expectedPath, path, "GetConfigPath result should match expected path")

	// Test case 2: Directory ends with specified suffix
	path, err = GetConfigPath("auth_service")
	assert.NoError(t, err, "GetConfigPath should not return an error")
	expectedPath = filepath.Join(tempDir, "auth_service/config/config.staging.yaml")
	assert.Equal(t, expectedPath, path, "GetConfigPath result should match expected path")

	// Test case 3: Directory ends with "interview"
	os.MkdirAll(filepath.Join(tempDir, "cmd", "suffix", "config"), 0755)
	os.Chdir(filepath.Join(tempDir, "interview"))
	path, err = GetConfigPath("gateway")
	assert.NoError(t, err, "GetConfigPath should not return an error")
	expectedPath = filepath.Join(tempDir, "gateway", "config", "config.staging.yaml")
	assert.Equal(t, expectedPath, path, "GetConfigPath result should match expected path")

}
