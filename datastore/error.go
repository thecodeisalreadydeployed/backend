package datastore

import "errors"

var (
	ErrFailedPrecondition = errors.New("ErrFailedPrecondition")
	ErrNotFound           = errors.New("ErrNotFound")
	ErrAlreadyExists      = errors.New("ErrAlreadyExists")
	ErrInvalidArgument    = errors.New("ErrInvalidArgument")
	ErrCannotCreate       = errors.New("ErrCannotCreate")
)

var (
	MsgProjectPrefix    = "Project must have a prj_ prefix."
	MsgAppPrefix        = "App must have an app_ prefix."
	MsgDeploymentPrefix = "Deployment must have a dpl_ prefix."
)
