package error

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
