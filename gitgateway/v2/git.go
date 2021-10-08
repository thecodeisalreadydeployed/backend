package gitgateway

type GitGateway interface {
	Checkout(branch string) error
	OpenFile(filePath string) error
	WriteFile(filePath string, data string) error
	Commit(files []string) (string, error)
	Pull() error
	Diff(oldCommit string, currentCommit string) ([]string, error)
}
