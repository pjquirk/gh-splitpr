package cmd

type PullRequest struct {
	Number int
	Title  string
	Author string
}

type Commit struct {
	Sha     string
	Comment string
}
