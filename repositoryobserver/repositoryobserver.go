package repositoryobserver

import (
	"github.com/thecodeisalreadydeployed/gitgateway"
)

type Repository struct {
	SourceCode    *gitgateway.GitGateway
	LastCommitSHA string
}

func (r *Repository) Init(path string) {
	it := gitgateway.NewGitGateway(path)
	r.SourceCode = &it
}

func (r *Repository) HasChanges() bool {
	old := r.SourceCode.GetCommit(r.LastCommitSHA)
	current := r.SourceCode.GetCurrentCommit()
	return gitgateway.HasProperDiff(old, current)
}
