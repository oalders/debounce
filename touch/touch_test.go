package touch_test

import (
	"os"
	"testing"

	"github.com/oalders/debounce/touch"
	"github.com/stretchr/testify/assert"
)

func TestTouch(t *testing.T) {
	t.Parallel()

	// Create a temporary file and get its last modification time
	tempFile, err := os.CreateTemp("", "test")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	fileInfo, err := tempFile.Stat()
	assert.NoError(t, err)
	lastModTime := fileInfo.ModTime()

	// Call touch on the temporary file
	err = touch.Touch(tempFile.Name())
	assert.NoError(t, err)

	// Check if the last modification time of the file has been updated
	fileInfo, err = tempFile.Stat()
	assert.NoError(t, err)
	assert.True(t, fileInfo.ModTime().After(lastModTime))

	// Call touch on a non-existent file
	nonExistentFile := tempFile.Name() + "nonexistent"
	err = touch.Touch(nonExistentFile)
	assert.NoError(t, err)
	defer os.Remove(nonExistentFile)

	// Check if the file was created
	_, err = os.Stat(nonExistentFile)
	assert.NoError(t, err)
}
