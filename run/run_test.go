package run_test

import (
	"os"
	"path"
	"testing"

	"github.com/oalders/debounce/run"
	"github.com/oalders/debounce/types"
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
	}
	success, output, err := run.Run(args, tempDir)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !success {
		t.Errorf("Expected success, got %v", success)
	}

	expectedOutput := "Hello, World!\n"
	if string(output) != expectedOutput {
		t.Errorf("Expected output %q, got %q", expectedOutput, output)
	}
}

func TestEnsureDir(t *testing.T) {
	t.Parallel()
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Errorf("Expected directory to exist at %s", tempDir)
	}

	dirName := "testDir"

	err = run.MaybeMakeCacheDir(tempDir, dirName)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if _, err := os.Stat(path.Join(tempDir, dirName)); os.IsNotExist(err) {
		t.Errorf("Expected directory to exist at %s", dirName)
	}
}
