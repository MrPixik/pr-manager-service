package dto

type TeamAddRequest struct {
	TeamName string              `json:"team_name"`
	Members  []TeamMemberRequest `json:"members"`
}

type TeamMemberRequest struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type GetTeamRequest struct {
	TeamName string `json:"team_name" validate:"required"`
}

type SetIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type PullRequestCreateRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type PullRequestMergeRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type PullRequestReassignRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldReviewerID string `json:"old_user_id"`
}

type GetReviewPRRequest struct {
	UserID string `json:"user_id"`
}

type GetTeamStatsRequest struct {
	TeamName string `json:"team_name"`
}
