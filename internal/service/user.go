package service

import (
	"context"
	"errors"
	"service-order-avito/internal/domain"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/repository"
	"service-order-avito/internal/domain/errors/service"
)

type UserRepository interface {
	SetIsActive(context.Context, string, bool) (*domain.User, error)
	GetReviewPullRequests(context.Context, string) ([]domain.PullRequest, error)
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *userService {
	return &userService{repo: repo}
}

func (s *userService) SetIsActive(ctx context.Context, req *dto.SetIsActiveRequest) (*dto.SetIsActiveResponse, error) {
	user, err := s.repo.SetIsActive(ctx, req.UserID, req.IsActive)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			return nil, service.ErrUserNotFound
		default:
			return nil, service.ErrInternalError
		}
	}

	return &dto.SetIsActiveResponse{
		User: dto.UserResponse{
			UserID:   user.ID,
			Username: user.Username,
			TeamName: user.TeamName,
			IsActive: user.IsActive,
		},
	}, nil
}

func (s *userService) GetReviewPullRequests(ctx context.Context, req *dto.GetReviewPRRequest) (*dto.GetReviewPRResponse, error) {
	prs, err := s.repo.GetReviewPullRequests(ctx, req.UserID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			return nil, service.ErrUserNotFound
		default:
			return nil, service.ErrInternalError
		}
	}

	respPRs := make([]dto.PullRequestShortResponse, len(prs))
	for i, pr := range prs {
		respPRs[i] = dto.PullRequestShortResponse{
			PullRequestID:   pr.ID,
			PullRequestName: pr.Name,
			AuthorID:        pr.AuthorID,
			Status:          pr.Status,
		}
	}

	return &dto.GetReviewPRResponse{
		UserID:       req.UserID,
		PullRequests: respPRs,
	}, nil
}
