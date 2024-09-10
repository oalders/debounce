package run_test

import (
	"os"
	"path"
	"testing"

	"github.com/oalders/debounce/run"
	"github.com/oalders/debounce/types"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err, "first attempt to make dir")

	if _, err := os.Stat(path.Join(tempDir, dirName)); os.IsNotExist(err) {
		t.Errorf("Expected directory to exist at %s", dirName)
	}

	err = run.MaybeMakeCacheDir(tempDir, dirName)
	assert.NoError(t, err, "command is idempotent")
}

func TestGenerateCmdName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want string
		args []string
	}{
		{
			name: "No special characters",
			args: []string{"arg1", "arg2", "arg3"},
			want: "arg1+arg2+arg3",
		},
		{
			name: "Contains slash",
			args: []string{"arg/1", "arg2", "arg3"},
			want: "arg%2F1+arg2+arg3",
		},
		{
			name: "Contains space",
			args: []string{"arg 1", "arg2", "arg3"},
			want: "arg+1+arg2+arg3",
		},
		{
			name: "Contains slash and space",
			args: []string{"arg/1", "arg 2", "arg3"},
			want: "arg%2F1+arg+2+arg3",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, run.GenerateCmdName(tt.args))
	}
}
