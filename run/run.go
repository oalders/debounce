package run

import (
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

func Run(args *types.DebounceCommand, home string) (bool, []byte, error) {
	command := args.Command[0]
	arguments := args.Command[1:]

	prettyName := strings.Join(args.Command, " ")

	cmdAsFilename := GenerateCacheFileName(prettyName)

	cacheDir := filepath.Join(".cache", "debounce")
	err := MaybeMakeCacheDir(home, cacheDir)
	if err != nil {
		return false, []byte{}, err
	}

	filename := filepath.Join(home, cacheDir, cmdAsFilename)
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
	fmt.Printf("Running command: %s %s\n", command, strings.Join(arguments, " "))
	cmd := exec.Command(command, arguments...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, output, errors.Join(errors.New("run command"), err)
	}
	err = touch.Touch(filename)
	if err != nil {
		return false, output, errors.Join(errors.New("touch"), err)
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

func GenerateCacheFileName(args string) string {
	hash := sha256.Sum256([]byte(args))
	return hex.EncodeToString(hash[:])
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
	if os.IsNotExist(err) {
		fmt.Println("Cache file does not exist. Command will run on next debounce")
		return false, []byte{}, nil
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
		fmt.Printf("🚀 \"%s\" will run on next invocation", prettyName)
	}
	return true, []byte{}, nil
}
