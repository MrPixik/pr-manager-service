package repository

import "errors"

var (
	ErrTeamAlreadyExists      = errors.New("team already exists")
	ErrTeamNotFound           = errors.New("team not found")
	ErrUserNotFound           = errors.New("user not found")
	ErrInternalError          = errors.New("internal error")
	ErrPullRequestExists      = errors.New("pull request already exists")
	ErrPullRequestNotFound    = errors.New("pull request not found")
	ErrPullRequestMerged      = errors.New("pull request already merged")
	ErrReviewerNotAssigned    = errors.New("reviewer not assigned")
	ErrNoReplacementCandidate = errors.New("no candidate for reassignment")
)
