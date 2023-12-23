package entity

import "errors"

var ErrNotFound = errors.New("entity not found")

func IsErrNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

var ErrInvalidEntity = errors.New("invalid entity")

func IsErrInvalidEntity(err error) bool {
	return errors.Is(err, ErrInvalidEntity)
}

var ErrInvalidID = errors.New("invalid ID")

func IsErrInvalidID(err error) bool {
	return errors.Is(err, ErrInvalidID)
}

var ErrAlreadyExists = errors.New("entity already exists")

func IsErrAlreadyExists(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}
