package errs

import "errors"

var (
	ErrNoNote = errors.New("note doesn't exist")
)
