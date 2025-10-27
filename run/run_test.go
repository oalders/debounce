package run_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/oalders/debounce/run"
	"github.com/oalders/debounce/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Parallel()
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	args := &types.DebounceCommand{
		Quantity: "1",
		Unit:     "s",
		Command:  []string{"echo", "Hello, World!"},
		Debug:    false,
		CacheDir: tempDir,
	}
	success, output, err := run.Run(args, tempDir)
	assert.NoError(t, err)
	assert.True(t, success)
	assert.Equal(t, string(output), "Hello, World!\n")
}

func TestEnsureDir(t *testing.T) {
	t.Parallel()
	tempDir, err := os.MkdirTemp("", "test")
	defer os.RemoveAll(tempDir)
	require.NoError(t, err)

	dirName := "testDir"

	err = run.MaybeMakeCacheDir(tempDir, dirName)
	require.NoError(t, err, "first attempt to make dir")

	_, err = os.Stat(filepath.Join(tempDir, dirName))
	require.NoError(t, err)

	err = run.MaybeMakeCacheDir(tempDir, dirName)
	assert.NoError(t, err, "command is idempotent")
}

func TestGenerateCacheFileName(t *testing.T) {
	t.Parallel()
	expected := "46e878132d529376c3d0b2d19ca9d9ab34bf3a940a92ae484689e1a271a61e84"
	for i := 0; i < 2; i++ {
		actual, err := run.GenerateCacheFileName("arg/1 arg 2 arg3", false)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
}

func TestRunWithNonExistentCacheDir(t *testing.T) {
	t.Parallel()
	nonExistentDir := "/non/existent/dir"

	args := &types.DebounceCommand{
		Quantity: "1",
		Unit:     "s",
		Command:  []string{"echo", "Hello, World!"},
		Debug:    false,
		CacheDir: nonExistentDir,
	}
	success, output, err := run.Run(args, "")
	assert.Error(t, err)
	assert.False(t, success)
	assert.Empty(t, output)
}

func TestRunCreatesDefaultCacheDir(t *testing.T) {
	t.Parallel()
	// Create a temporary home directory
	tempHome, err := os.MkdirTemp("", "test-home")
	require.NoError(t, err)
	defer os.RemoveAll(tempHome)

	// Don't create the cache dir ahead of time - let Run create it
	args := &types.DebounceCommand{
		Quantity: "1",
		Unit:     "s",
		Command:  []string{"echo", "Hello, World!"},
		Debug:    false,
	}

	// Run should succeed even when cache dir doesn't exist
	success, output, err := run.Run(args, tempHome)
	assert.NoError(t, err)
	assert.True(t, success)
	assert.Equal(t, string(output), "Hello, World!\n")

	// Verify the cache directory was created
	cacheDir := filepath.Join(tempHome, ".cache", "debounce")
	_, err = os.Stat(cacheDir)
	assert.NoError(t, err, "cache directory should be created")

	// Verify the cache file was created
	entries, err := os.ReadDir(cacheDir)
	require.NoError(t, err)
	assert.Len(t, entries, 1, "should have created one cache file")
}
