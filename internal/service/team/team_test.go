package team

import (
	"context"
	"errors"
	"testing"

	"service-order-avito/internal/domain"
	"service-order-avito/internal/domain/dto"
	repoErr "service-order-avito/internal/domain/errors/repository"
	serviceErr "service-order-avito/internal/domain/errors/service"
	"service-order-avito/internal/service/team/mocks"

	"github.com/golang/mock/gomock"
)

func TestTeamService_AddTeam(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.TeamAddRequest
		mockErr error
		wantErr error
	}{
		{
			name: "success",
			req: &dto.TeamAddRequest{
				TeamName: "backend",
				Members: []dto.TeamMemberRequest{
					{UserID: "u1", Username: "Alice", IsActive: true},
					{UserID: "u2", Username: "Bob", IsActive: false},
				},
			},
			mockErr: nil,
			wantErr: nil,
		},
		{
			name:    "team already exists",
			req:     &dto.TeamAddRequest{TeamName: "backend"},
			mockErr: repoErr.ErrTeamAlreadyExists,
			wantErr: serviceErr.ErrTeamAlreadyExists,
		},
		{
			name:    "internal repository error",
			req:     &dto.TeamAddRequest{TeamName: "backend"},
			mockErr: errors.New("unknown repo error"),
			wantErr: serviceErr.ErrInternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mocks.NewMockTeamRepository(ctrl)
			svc := NewTeamService(mockRepo)

			ctx := context.Background()

			mockRepo.EXPECT().
				AddTeamWithMembers(ctx, gomock.Any(), gomock.Any()).
				Return(tt.mockErr).
				Times(1)

			resp, err := svc.AddTeam(ctx, tt.req)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if tt.wantErr == nil {
				if resp.TeamName != tt.req.TeamName {
					t.Fatalf("expected TeamName=%s, got %s", tt.req.TeamName, resp.TeamName)
				}
				if len(resp.Members) != len(tt.req.Members) {
					t.Fatalf("members length mismatch: expected %d, got %d",
						len(tt.req.Members), len(resp.Members))
				}
			}
		})
	}
}

func TestTeamService_GetTeam(t *testing.T) {
	tests := []struct {
		name       string
		req        *dto.GetTeamRequest
		mockReturn *domain.TeamWithUsers
		mockErr    error
		wantErr    error
	}{
		{
			name: "success",
			req:  &dto.GetTeamRequest{TeamName: "backend"},
			mockReturn: &domain.TeamWithUsers{
				TeamName: "backend",
				Members: []domain.User{
					{ID: "u1", Username: "Alice", IsActive: true},
					{ID: "u2", Username: "Bob", IsActive: false},
				},
			},
			mockErr: nil,
			wantErr: nil,
		},
		{
			name:       "team not found",
			req:        &dto.GetTeamRequest{TeamName: "backend"},
			mockReturn: nil,
			mockErr:    repoErr.ErrTeamNotFound,
			wantErr:    serviceErr.ErrTeamNotFound,
		},
		{
			name:       "internal repository error",
			req:        &dto.GetTeamRequest{TeamName: "backend"},
			mockReturn: nil,
			mockErr:    errors.New("unknown error"),
			wantErr:    serviceErr.ErrInternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mocks.NewMockTeamRepository(ctrl)
			svc := NewTeamService(mockRepo)

			ctx := context.Background()

			mockRepo.EXPECT().
				GetTeamWithMembers(ctx, tt.req.TeamName).
				Return(tt.mockReturn, tt.mockErr).
				Times(1)

			resp, err := svc.GetTeam(ctx, tt.req)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if tt.wantErr == nil {
				if resp.TeamName != tt.mockReturn.TeamName {
					t.Fatalf("wrong team name: expected %s, got %s", tt.mockReturn.TeamName, resp.TeamName)
				}

				expectedMembers := len(tt.mockReturn.Members)
				if len(resp.Members) != expectedMembers {
					t.Fatalf("expected %d members, got %d", expectedMembers, len(resp.Members))
				}

				for i := range resp.Members {
					if resp.Members[i].UserID != tt.mockReturn.Members[i].ID ||
						resp.Members[i].Username != tt.mockReturn.Members[i].Username ||
						resp.Members[i].IsActive != tt.mockReturn.Members[i].IsActive {
						t.Fatalf("member mismatch at index %d", i)
					}
				}
			}
		})
	}
}
