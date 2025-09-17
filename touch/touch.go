package touch

import (
	"context"
	"errors"
	"os/exec"
)

func Touch(filename string) error {
	cmd := exec.CommandContext(context.Background(), "touch", filename)
	err := cmd.Start()
	if err != nil {
		return errors.Join(errors.New("start command"), err)
	}

	err = cmd.Wait()
	if err != nil {
		return errors.Join(errors.New("wait for command"), err)
	}
	return nil
}
