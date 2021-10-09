package model

type GitSource struct {
	CommitSHA        string `json:"commit_sha"`
	CommitMessage    string `json:"commit_message"`
	CommitAuthorName string `json:"commit_author_name"`
	RepositoryURL    string `json:"repository_url"`
	Branch           string `json:"branch"`
}
