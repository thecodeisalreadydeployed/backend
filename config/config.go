package config

import (
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/viper"
	"github.com/thecodeisalreadydeployed/constant"
)

func DefaultGitSignature() *object.Signature {
	return &object.Signature{
		Name:  "GitHub Action",
		Email: "action@github.com",
		When:  time.Now(),
	}
}

const (
	DefaultKanikoWorkingDirectory string = "/__w/kaniko/"
)

func DefaultUserspaceRepository() string {
	return viper.GetString(constant.UserspaceRepository)
}
