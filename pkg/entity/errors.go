package entity

import "errors"

var ErrNotFound = errors.New("entity not found")

func IsErrNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
