package repositoryobserver

import (
	"github.com/thecodeisalreadydeployed/gitinteractor"
)

type Repository struct {
	SourceCode    *gitinteractor.GitInteractor
	LastCommitSHA string
}

func (r *Repository) Init(path string) {
	it := gitinteractor.NewGitInteractor(path)
	r.SourceCode = &it
}

func (r *Repository) HasChanges() bool {
	old := r.SourceCode.GetCommit(r.LastCommitSHA)
	current := r.SourceCode.GetCurrentCommit()
	return gitinteractor.HasProperDiff(old, current)
}
