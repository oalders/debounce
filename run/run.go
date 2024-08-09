package run

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/oalders/debounce/age"
	"github.com/oalders/debounce/touch"
	"github.com/oalders/debounce/types"
)

func Run(args *types.DebounceCommand, home string) (bool, []byte, error) {
	command := args.Command[0]
	arguments := args.Command[1:]

	prettyName := strings.Join(args.Command, " ")

	cmdName := GenerateCmdName(args.Command)

	cacheDir := filepath.Join(".cache", "debounce")
	err := MaybeMakeCacheDir(home, cacheDir)
	if err != nil {
		return false, []byte{}, err
	}

	fileName := filepath.Join(home, cacheDir, cmdName)
	if args.Debug {
		fmt.Printf("Looking for debounce file \"%s\"\n", fileName)
	}

	tooSoon, err := age.Compare(fileName, args.Quantity, args.Unit)
	if err != nil {
		return false, []byte{}, errors.Join(errors.New(`checking last modified time`), err)
	}
	if tooSoon {
		fmt.Printf(
			"🚥 will not run \"%s\" more than once every %s %s\n",
			prettyName,
			args.Quantity,
			args.Unit,
		)
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
	err = touch.Touch(fileName)
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

func GenerateCmdName(args []string) string {
	cmdName := strings.Join(args, "-")
	cmdName = strings.ReplaceAll(cmdName, "/", "-")
	cmdName = strings.ReplaceAll(cmdName, " ", "-")
	return cmdName
}
