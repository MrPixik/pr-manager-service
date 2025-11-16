package domain

import "time"

const (
	PRStatusOpen   = "OPEN"
	PRStatusMerged = "MERGED"
)

type PullRequest struct {
	ID        string
	Name      string
	AuthorID  string
	Status    string
	CreatedAt time.Time
	MergedAt  *time.Time
}

type PullRequestWithReviewers struct {
	PullRequest
	AssignedReviewers []string
}

type Reviewer struct {
	ID string
}
