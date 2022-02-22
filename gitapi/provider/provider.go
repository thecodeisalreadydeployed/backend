package provider

import "github.com/thecodeisalreadydeployed/model"

type GitProvider interface {
	GetBranches() ([]string, error)
	GetFiles(branch string) ([]string, error)
	GetRaw(branch string, path string) (string, error)
	FillGitSource(gs *model.GitSource) (*model.GitSource, error)
}
