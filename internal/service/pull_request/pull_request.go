package pull_request

import (
	"context"
	"service-order-avito/internal/domain"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/service/error_wrapper"
)

// mockgen -source="internal/service/pull_request/pull_request.go" -destination="internal/service/pull_request/mocks/mock_pull_request_repository.go" -package=mocks PullRequestRepository
type PullRequestRepository interface {
	CreateWithReviewers(context.Context, domain.PullRequest) (*domain.PullRequestWithReviewers, error)
	Merge(context.Context, string) (*domain.PullRequestWithReviewers, error)
	ReassignReviewer(context.Context, string, string) (*domain.Reviewer, error)
}

type pullRequestService struct {
	repo PullRequestRepository
}

func NewPullRequestService(repo PullRequestRepository) *pullRequestService {
	return &pullRequestService{repo: repo}
}

func (s *pullRequestService) Create(ctx context.Context, req *dto.PullRequestCreateRequest) (*dto.PullRequestCreateResponse, error) {
	prDomain := domain.PullRequest{
		ID:       req.PullRequestID,
		Name:     req.PullRequestName,
		AuthorID: req.AuthorID,
		Status:   "OPEN",
	}

	prWithReviewers, err := s.repo.CreateWithReviewers(ctx, prDomain)
	if err != nil {
		return nil, error_wrapper.WrapRepositoryError(err)
	}

	resp := &dto.PullRequestCreateResponse{
		PullRequest: dto.PullRequestResponse{
			PullRequestID:     prWithReviewers.ID,
			PullRequestName:   prWithReviewers.Name,
			AuthorID:          prWithReviewers.AuthorID,
			Status:            prWithReviewers.Status,
			AssignedReviewers: prWithReviewers.AssignedReviewers,
		},
	}

	return resp, nil
}

func (s *pullRequestService) Merge(ctx context.Context, req *dto.PullRequestMergeRequest) (*dto.PullRequestMergeResponse, error) {
	prWithReviewers, err := s.repo.Merge(ctx, req.PullRequestID)
	if err != nil {
		return nil, error_wrapper.WrapRepositoryError(err)
	}

	resp := &dto.PullRequestMergeResponse{
		PullRequest: dto.PullRequestMergedResponse{
			PullRequestID:     prWithReviewers.ID,
			PullRequestName:   prWithReviewers.Name,
			AuthorID:          prWithReviewers.AuthorID,
			Status:            prWithReviewers.Status,
			AssignedReviewers: prWithReviewers.AssignedReviewers,
			MergedAt:          prWithReviewers.MergedAt,
		},
	}

	return resp, nil
}

func (s *pullRequestService) ReassignReviewer(ctx context.Context, req *dto.PullRequestReassignRequest) (*dto.PullRequestReassignResponse, error) {
	reviewer, err := s.repo.ReassignReviewer(ctx, req.PullRequestID, req.OldReviewerID)
	if err != nil {
		return nil, error_wrapper.WrapRepositoryError(err)
	}

	resp := &dto.PullRequestReassignResponse{
		ReplacedBy: reviewer.ID,
	}

	return resp, nil
}
