package errutil

import (
	"errors"
)

// Definitions of common error types used throughout thecodeisalreadydeployed/backend.
// All errors returned by most packages will map into one of these errors.

var (
	ErrUnknown            = errors.New("unknown")
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrNotFound           = errors.New("not found")
	ErrAlreadyExists      = errors.New("already exists")
	ErrFailedPrecondition = errors.New("failed precondition")
	ErrUnavailable        = errors.New("unavailable")
	ErrNotImplemented     = errors.New("not implemented")
)

// IsInvalidArgument returns true if the error is due to an invalid argument.
func IsInvalidArgument(err error) bool {
	return errors.Is(err, ErrInvalidArgument)
}

// IsNotFound returns true if the error is due to a missing object.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsAlreadyExists returns true if the error is due to an already existing item.
func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}

// IsFailedPrecondition returns true if an operation could not proceed due to
// the lack of a particular condition.
func IsFailedPrecondition(err error) bool {
	return errors.Is(err, ErrFailedPrecondition)
}

// IsUnavilable returns true if the error is due to a resource being unavailable.
func IsUnavilable(err error) bool {
	return errors.Is(err, ErrUnavailable)
}

// IsNotImplemented returns true if the error is due to not being implemented.
func IsNotImplemented(err error) bool {
	return errors.Is(err, ErrNotImplemented)
}
