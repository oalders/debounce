// Package main contains the logic for the "cli" command
package age

import (
	"errors"
	"os"
	"time"

	"github.com/oalders/is/age"
)

func Compare(path, ageValue, ageUnit string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Join(errors.New("could not stat cache file"), err)
	}

	dur, err := age.StringToDuration(ageValue, ageUnit)
	if err != nil {
		return false, err
	}
	targetTime := time.Now().Add(*dur)

	return info.ModTime().Compare(targetTime) >= 0, nil
}
