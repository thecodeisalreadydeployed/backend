package provider

type GitProvider interface {
	GetBranches() ([]string, error)
	GetFiles(branch string) ([]string, error)
	GetRaw(branch string, path string) (string, error)
}
