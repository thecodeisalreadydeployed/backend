package model

type GitSource struct {
	Provider         string
	Organization     string
	CommitSHA        string
	CommitMessage    string
	CommitAuthorName string
}
