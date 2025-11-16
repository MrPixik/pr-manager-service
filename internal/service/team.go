package service

import (
	"context"
	"errors"
	"service-order-avito/internal/domain"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/repository"
	"service-order-avito/internal/domain/errors/service"
)

type TeamRepository interface {
	AddTeamWithMembers(context.Context, domain.Team, []domain.User) error
	GetTeamWithMembers(context.Context, string) (*domain.TeamWithUsers, error)
}

type teamService struct {
	repo TeamRepository
}

func NewTeamService(repo TeamRepository) *teamService {
	return &teamService{repo: repo}
}

func (s *teamService) AddTeam(ctx context.Context, req *dto.TeamAddRequest) (*dto.AddTeamResponse, error) {
	team := domain.Team{Name: req.TeamName}

	members := make([]domain.User, len(req.Members))
	for i, m := range req.Members {
		members[i] = domain.User{
			ID:       m.UserID,
			Username: m.Username,
			TeamName: req.TeamName,
			IsActive: m.IsActive,
		}
	}

	if err := s.repo.AddTeamWithMembers(ctx, team, members); err != nil {
		switch {
		case errors.Is(err, repository.ErrTeamAlreadyExists):
			return nil, service.ErrTeamAlreadyExists
		default:
			return nil, service.ErrInternalError
		}
	}

	respMembers := make([]dto.TeamMemberResponse, len(members))
	for i, u := range members {
		respMembers[i] = dto.TeamMemberResponse{
			UserID:   u.ID,
			Username: u.Username,
			IsActive: u.IsActive,
		}
	}

	resp := dto.AddTeamResponse{
		TeamName: team.Name,
		Members:  respMembers,
	}

	return &resp, nil
}

func (s *teamService) GetTeam(ctx context.Context, req *dto.GetTeamRequest) (*dto.GetTeamResponse, error) {
	agg, err := s.repo.GetTeamWithMembers(ctx, req.TeamName)
	if err != nil {
		switch err {
		case repository.ErrTeamNotFound:
			return nil, service.ErrTeamNotFound
		default:
			return nil, service.ErrInternalError
		}
	}

	members := make([]dto.TeamMemberResponse, len(agg.Members))
	for i, m := range agg.Members {
		members[i] = dto.TeamMemberResponse{
			UserID:   m.ID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}

	return &dto.GetTeamResponse{
		TeamName: agg.TeamName,
		Members:  members,
	}, nil
}
