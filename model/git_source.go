package model

import "github.com/thecodeisalreadydeployed/gitgateway"

type GitSource struct {
	Provider              string                 `json:"provider"`
	Organization          string                 `json:"organization"`
	CommitSHA             string                 `json:"commit_sha"`
	CommitMessage         string                 `json:"commit_message"`
	CommitAuthorName      string                 `json:"commit_author_name"`
	RepositoryURL         string                 `json:"repository_url"`
	LastObservedCommitSHA string                 `json:"last_observed_commit_sha"`
	SourceCode            *gitgateway.GitGateway `json:"source_code"`
}
