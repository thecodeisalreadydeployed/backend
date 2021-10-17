package model

type GitSource struct {
	CommitSHA        string `json:"commitSHA"`
	CommitMessage    string `json:"commitMessage"`
	CommitAuthorName string `json:"commitAuthorName"`
	RepositoryURL    string `json:"repositoryURL"`
	Branch           string `json:"branch"`
}
