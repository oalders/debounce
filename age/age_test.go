package age_test

import (
	"os"
	"testing"
	"time"

	"github.com/oalders/debounce/age"
	"github.com/stretchr/testify/require"
)

func TestCompare(t *testing.T) {
	t.Parallel()
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test")
	require.NoError(t, err, "Failed to create temp file")
	defer os.Remove(tempFile.Name()) // clean up

	// Test Compare function
	tooSoon, err := age.Compare(tempFile.Name(), "1", "s")
	require.NoError(t, err, "Error in Compare function")

	// The file was just created, so it should be too soon
	require.True(t, tooSoon, "Expected true, got false")

	// Wait for 2 seconds
	time.Sleep(2 * time.Second)

	// Now it should not be too soon
	tooSoon, err = age.Compare(tempFile.Name(), "1", "s")
	require.NoError(t, err, "Error in Compare function after sleep")
	require.False(t, tooSoon, "Expected false, got true")
}

func TestCompareNonExistentFile(t *testing.T) {
	t.Parallel()

	tooSoon, err := age.Compare("idonotexist", "1", "s")
	require.Nil(t, err, "no error if file does not exist")

	// The file does not exist, so it should not be too soon
	require.False(t, tooSoon, "Expected false for non-existent file")
}
