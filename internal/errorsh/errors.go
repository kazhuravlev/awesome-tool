package errorsh

import "fmt"

// Wrap helps to wrap errors.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", msg, err)
}
