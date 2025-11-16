package service

import (
	"context"
	"errors"
	"service-order-avito/internal/domain"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/repository"
	"service-order-avito/internal/domain/errors/service"
)

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
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			return nil, service.ErrUserNotFound
		case errors.Is(err, repository.ErrTeamNotFound):
			return nil, service.ErrTeamNotFound
		case errors.Is(err, repository.ErrPullRequestExists):
			return nil, service.ErrPullRequestExists
		default:
			return nil, service.ErrInternalError
		}
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
		switch {
		case errors.Is(err, repository.ErrPullRequestNotFound):
			return nil, service.ErrPullRequestNotFound
		default:
			return nil, service.ErrInternalError
		}
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
		switch {
		case errors.Is(err, repository.ErrPullRequestNotFound):
			return nil, service.ErrPullRequestNotFound
		case errors.Is(err, repository.ErrUserNotFound):
			return nil, service.ErrUserNotFound
		case errors.Is(err, repository.ErrPullRequestMerged):
			return nil, service.ErrPullRequestMerged
		case errors.Is(err, repository.ErrReviewerNotAssigned):
			return nil, service.ErrReviewerNotAssigned
		case errors.Is(err, repository.ErrNoReplacementCandidate):
			return nil, service.ErrNoReplacementCandidate
		default:
			return nil, service.ErrInternalError
		}
	}

	resp := &dto.PullRequestReassignResponse{
		ReplacedBy: reviewer.ID,
	}

	return resp, nil
}
