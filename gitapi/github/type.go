package github

type Branch struct {
	Name      string `json:"name"`
	Protected bool   `json:"protected"`
	Commit    Commit
}

type Commit struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

type File struct {
	SHA  string `json:"sha"`
	URL  string `json:"url"`
	Tree []Tree
}

type Tree struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	SHA  string `json:"sha"`
	URL  string `json:"url"`
}

type CommitDetails struct {
	SHA    string            `json:"sha"`
	NodeId string            `json:"node_id"`
	Commit CommitMoreDetails `json:"commit"`
}

type CommitMoreDetails struct {
	Author    Person `json:"author"`
	Committer Person `json:"committer"`
	Message   string `json:"message"`
}

type Person struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}
