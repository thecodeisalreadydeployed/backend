package gitapi

type GitAPIBackend interface {
	GetBranches(repoURL string) ([]string, error)
	GetFiles(repoURL string, branch string) ([]string, error)
	GetRaw(repoURL string, branch string, path string) (string, error)
}

type gitAPIBackend struct{}

func NewGitAPIBackend() GitAPIBackend {
	return &gitAPIBackend{}
}
