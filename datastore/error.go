package datastore

import "errors"

var (
	ErrFailedPrecondition = errors.New("ErrFailedPrecondition")
	ErrNotFound           = errors.New("ErrNotFound")
	ErrAlreadyExists      = errors.New("ErrAlreadyExists")
	ErrInvalidArgument    = errors.New("ErrInvalidArgument")
)
