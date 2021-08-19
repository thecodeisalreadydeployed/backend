package model

type GitSource struct {
	Provider         string `json:"provider"`
	Organization     string `json:"organization"`
	CommitSHA        string `json:"commit_sha"`
	CommitMessage    string `json:"commit_message"`
	CommitAuthorName string `json:"commit_author_name"`
}
