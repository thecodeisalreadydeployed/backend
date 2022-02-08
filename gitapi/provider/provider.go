package provider

type GitProvider interface {
	GetBranches(owner string, repo string) ([]string, error)
	GetFiles(owner string, repo string, branch string) ([]string, error)
	GetRaw(owner string, repo string, branch string, path string) (string, error)
}
