package run

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/oalders/debounce/age"
	"github.com/oalders/debounce/touch"
	"github.com/oalders/debounce/types"
	is_age "github.com/oalders/is/age"
)

func Run(args *types.DebounceCommand, home string) (bool, []byte, error) { //nolint:cyclop
	command := args.Command[0]
	arguments := args.Command[1:]

	prettyName := strings.Join(args.Command, " ")

	cmdAsFilename, err := GenerateCacheFileName(prettyName, args.Local)
	if err != nil {
		return false, []byte{}, err
	}

	var cacheDir string
	if args.CacheDir != "" {
		cacheDir = args.CacheDir
		if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
			return false, []byte{}, fmt.Errorf("provided cache directory does not exist: %s", cacheDir)
		}
	} else {
		cacheDir = filepath.Join(home, ".cache", "debounce")
		err := MaybeMakeCacheDir(home, filepath.Join(".cache", "debounce"))
		if err != nil {
			return false, []byte{}, err
		}
	}

	filename := filepath.Join(cacheDir, cmdAsFilename)
	if args.Debug {
		fmt.Printf("Looking for debounce file \"%s\"\n", filename)
	}

	tooSoon, err := age.Compare(filename, args.Quantity, args.Unit)
	if err != nil {
		return false, []byte{}, errors.Join(errors.New(`checking last modified time`), err)
	}
	if args.Status {
		return HandleStatus(args, filename, tooSoon, prettyName)
	}

	if tooSoon {
		TooSoon(args, prettyName)
		return true, []byte{}, nil
	}

	// This is just like running any other user command, so assume user has
	// already sanitized inputs.
	if args.Debug {
		fmt.Printf("🚀 Running command: %s %s\n", command, strings.Join(arguments, " "))
	}
	cmd := exec.CommandContext(context.Background(), command, arguments...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, output, errors.Join(
			fmt.Errorf("running command: %s %s", command, strings.Join(arguments, " ")),
			err,
		)
	}
	err = touch.Touch(filename)
	if err != nil {
		return false, output, errors.Join(errors.New("touch "+filename), err)
	}
	return true, output, nil
}

func MaybeMakeCacheDir(parent, cache string) error {
	fullPath := filepath.Join(parent, cache)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		err = os.MkdirAll(fullPath, 0o755)
		if err != nil {
			return errors.Join(errors.New("make cache dir"), err)
		}
	}

	return nil
}

func GenerateCacheFileName(args string, local bool) (string, error) {
	if local {
		cwd, err := os.Getwd()
		if err != nil {
			return "", errors.Join(errors.New("get current working directory"), err)
		}
		args = cwd + " " + args
	}
	hash := sha256.Sum256([]byte(args))
	return hex.EncodeToString(hash[:]), nil
}

func TooSoon(args *types.DebounceCommand, cmd string) {
	fmt.Printf(
		"🚥 will not run \"%s\" more than once every %s %s\n",
		cmd,
		args.Quantity,
		args.Unit,
	)
}

func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func HandleStatus(
	args *types.DebounceCommand,
	filename string,
	tooSoon bool,
	prettyName string,
) (bool, []byte, error) {
	fileInfo, err := os.Stat(filename)
	if args.Local {
		cwd, err := os.Getwd()
		if err != nil {
			return false, []byte{}, errors.Join(errors.New("get current working directory"), err)
		}
		fmt.Printf("🏘️ command localized to %s\n", cwd)
	}
	if os.IsNotExist(err) {
		fmt.Printf("Cache file does not exist. \"%s\" will run on next debounce\n", prettyName)
		return true, []byte{}, nil
	} else if err != nil {
		return false, []byte{}, errors.Join(errors.New("stat file"), err)
	}

	fileAge := time.Since(fileInfo.ModTime())
	debounceInterval, err := is_age.StringToDuration(args.Quantity, args.Unit)
	if err != nil {
		return false, []byte{}, err
	}

	fmt.Printf("📁 cache location: %s\n", filename)
	fmt.Printf("🚧 cache last modified: %s\n", fileInfo.ModTime().Format(time.RFC1123))
	fmt.Printf("⏲️ debounce interval: %s\n", FormatDuration(debounceInterval.Abs()))
	fmt.Printf("🕰️ cache age: %s\n", FormatDuration(fileAge))
	if tooSoon {
		waitTime := -1**debounceInterval - time.Since(fileInfo.ModTime())
		fmt.Printf("⏳ time remaining: %s\n", FormatDuration(waitTime))
	} else {
		fmt.Printf("🚀 \"%s\" will run on next invocation\n", prettyName)
	}
	return true, []byte{}, nil
}
