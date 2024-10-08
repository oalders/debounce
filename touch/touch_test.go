package touch_test

import (
	"os"
	"testing"
	"time"

	"github.com/oalders/debounce/touch"
	"github.com/stretchr/testify/assert"
)

func TestTouch(t *testing.T) {
	t.Parallel()

	// Create a temporary file and get its last modification time
	tempFile, err := os.CreateTemp("", "test")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())
	before, err := tempFile.Stat()
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)
	err = touch.Touch(tempFile.Name())
	assert.NoError(t, err)

	after, err := tempFile.Stat()
	assert.NoError(t, err)
	assert.True(t, after.ModTime().After(before.ModTime()))

	// Call touch on a non-existent file
	nonExistentFile := tempFile.Name() + "nonexistent"
	err = touch.Touch(nonExistentFile)
	assert.NoError(t, err)
	defer os.Remove(nonExistentFile)

	// Check if the file was created
	_, err = os.Stat(nonExistentFile)
	assert.NoError(t, err)
}
