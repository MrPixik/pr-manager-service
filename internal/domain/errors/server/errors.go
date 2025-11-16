package server

const (
	ErrInvalidJSON            = "invalid JSON"
	ErrTeamAlreadyExists      = "team already exists"
	ErrTeamNotFound           = "team not found"
	ErrUserNotFound           = "user not found"
	ErrPRAlreadyExists        = "PR id already exists"
	ErrPRNotFound             = "PR not found"
	ErrPullRequestMerged      = "cannot reassign on merged PR"
	ErrReviewerNotAssigned    = "reviewer is not assigned to this PR"
	ErrNoReplacementCandidate = "no candidate for reassignment"
	ErrRequestCanceled        = "request canceled"
	ErrInternalError          = "internal error"
)
