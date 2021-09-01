package main

import (
	"fmt"
	"github.com/thecodeisalreadydeployed/gitinteractor"
	"github.com/thecodeisalreadydeployed/repositoryobserver"
)

func main() {
	path := "/home/jean/Desktop/gittest"
	it := gitinteractor.NewGitInteractor(path)
	repo := repositoryobserver.Repository{
		SourceCode:    &it,
		LastCommitSHA: "a7402afdac1147ed8908055e1d511f11418714c7",
	}
	fmt.Println(repo.HasChanges())
}
