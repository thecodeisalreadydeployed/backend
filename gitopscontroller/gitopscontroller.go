package gitopscontroller

import "github.com/thecodeisalreadydeployed/gitinteractor"

type GitOpsController struct {
	gitInteractor *gitinteractor.GitInteractor
}

func (c *GitOpsController) Init() {
	git := gitinteractor.NewGitInteractor()
	c.gitInteractor = &git
}
