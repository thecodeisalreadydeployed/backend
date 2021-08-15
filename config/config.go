package config

import (
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
)

func DefaultGitSignature() object.Signature {
	return object.Signature{
		Name:  "GitHub Action",
		Email: "action@github.com",
		When:  time.Now(),
	}
}
