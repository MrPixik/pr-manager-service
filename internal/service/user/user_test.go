package user

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"service-order-avito/internal/domain"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/repository"
	"service-order-avito/internal/domain/errors/service"
	"service-order-avito/internal/service/user/mocks"
	"testing"
)

func TestUserService_SetIsActive(t *testing.T) {
	tests := []struct {
		name          string
		req           *dto.SetIsActiveRequest
		mockUser      *domain.User
		mockErr       error
		expectedResp  *dto.SetIsActiveResponse
		expectedError error
	}{
		{
			name: "success",
			req: &dto.SetIsActiveRequest{
				UserID:   "u1",
				IsActive: true,
			},
			mockUser: &domain.User{
				ID:       "u1",
				Username: "Alice",
				TeamName: "team1",
				IsActive: true,
			},
			mockErr: nil,
			expectedResp: &dto.SetIsActiveResponse{
				User: dto.UserResponse{
					UserID:   "u1",
					Username: "Alice",
					TeamName: "team1",
					IsActive: true,
				},
			},
			expectedError: nil,
		},
		{
			name: "user not found",
			req: &dto.SetIsActiveRequest{
				UserID:   "u999",
				IsActive: true,
			},
			mockUser:      nil,
			mockErr:       repository.ErrUserNotFound,
			expectedResp:  nil,
			expectedError: service.ErrUserNotFound,
		},
		{
			name: "internal error",
			req: &dto.SetIsActiveRequest{
				UserID:   "u2",
				IsActive: false,
			},
			mockUser:      nil,
			mockErr:       errors.New("db error"),
			expectedResp:  nil,
			expectedError: service.ErrInternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockUserRepository(ctrl)
			mockRepo.EXPECT().
				SetIsActive(gomock.Any(), tt.req.UserID, tt.req.IsActive).
				Return(tt.mockUser, tt.mockErr)

			svc := NewUserService(mockRepo)
			resp, err := svc.SetIsActive(context.Background(), tt.req)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResp, resp)
		})
	}
}

func TestUserService_GetReviewPullRequests(t *testing.T) {
	tests := []struct {
		name          string
		req           *dto.GetReviewPRRequest
		mockPRs       []domain.PullRequest
		mockErr       error
		expectedResp  *dto.GetReviewPRResponse
		expectedError error
	}{
		{
			name: "success",
			req: &dto.GetReviewPRRequest{
				UserID: "u1",
			},
			mockPRs: []domain.PullRequest{
				{ID: "pr1", Name: "Add feature", AuthorID: "u2", Status: "OPEN"},
				{ID: "pr2", Name: "Fix bug", AuthorID: "u3", Status: "MERGED"},
			},
			mockErr: nil,
			expectedResp: &dto.GetReviewPRResponse{
				UserID: "u1",
				PullRequests: []dto.PullRequestShortResponse{
					{PullRequestID: "pr1", PullRequestName: "Add feature", AuthorID: "u2", Status: "OPEN"},
					{PullRequestID: "pr2", PullRequestName: "Fix bug", AuthorID: "u3", Status: "MERGED"},
				},
			},
			expectedError: nil,
		},
		{
			name: "user not found",
			req: &dto.GetReviewPRRequest{
				UserID: "u999",
			},
			mockPRs:       nil,
			mockErr:       repository.ErrUserNotFound,
			expectedResp:  nil,
			expectedError: service.ErrUserNotFound,
		},
		{
			name: "internal error",
			req: &dto.GetReviewPRRequest{
				UserID: "u2",
			},
			mockPRs:       nil,
			mockErr:       errors.New("db error"),
			expectedResp:  nil,
			expectedError: service.ErrInternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockUserRepository(ctrl)
			mockRepo.EXPECT().
				GetReviewPullRequests(gomock.Any(), tt.req.UserID).
				Return(tt.mockPRs, tt.mockErr)

			svc := NewUserService(mockRepo)
			resp, err := svc.GetReviewPullRequests(context.Background(), tt.req)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResp, resp)
		})
	}
}
