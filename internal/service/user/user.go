package user

import (
	"context"
	"service-order-avito/internal/domain"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/service/error_wrapper"
)

// mockgen -source="internal/service/user/user.go" -destination="internal/service/user/mocks/mock_user_repository.go" -package=mocks UserRepository
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
		return nil, error_wrapper.WrapRepositoryError(err)
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
		return nil, error_wrapper.WrapRepositoryError(err)
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
