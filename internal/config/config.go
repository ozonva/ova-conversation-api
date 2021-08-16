package config

import (
	"errors"
	"fmt"
	"os"
	"time"
)

const numOfAttempts = 5

func Update(path string, duration int64) error {
	update := func(path string) error {
		f, err := os.Open(path)
		if err != nil {
			textErr := fmt.Sprintf("fail to open file %q, error: %q", path, err)
			return errors.New(textErr)
		}

		var closeErr error
		defer func() {
			closeErr := f.Close()
			if closeErr != nil {
				textErr := fmt.Sprintf("fail to close file %q, error: %q", path, closeErr)
				closeErr = errors.New(textErr)
			}
		}()

		return closeErr
	}

	for i := 0; i < numOfAttempts; i++ {
		err := update(path)
		if err != nil {
			return err
		}
		time.Sleep(time.Duration(duration) * time.Second)
	}

	return nil
}
