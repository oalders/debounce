package run_test

import (
	"os"
	"path"
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

	_, err = os.Stat(path.Join(tempDir, dirName))
	require.NoError(t, err)

	err = run.MaybeMakeCacheDir(tempDir, dirName)
	assert.NoError(t, err, "command is idempotent")
}

func TestGenerateCacheFileName(t *testing.T) {
	t.Parallel()
	expected := "46e878132d529376c3d0b2d19ca9d9ab34bf3a940a92ae484689e1a271a61e84"
	for i := 0; i < 2; i++ {
		actual := run.GenerateCacheFileName("arg/1 arg 2 arg3")
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
