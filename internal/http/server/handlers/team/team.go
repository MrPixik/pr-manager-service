package team

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/server"
	"service-order-avito/internal/domain/errors/service"
	"service-order-avito/internal/http/codes"
	"service-order-avito/pkg/http/error_wrapper"
)

// mockgen -source="internal/http/server/handlers/team/team.go" -destination="internal/http/server/handlers/team/mocks/mock_user_service.go" -package=mocks TeamService
type TeamService interface {
	AddTeam(context.Context, *dto.TeamAddRequest) (*dto.AddTeamResponse, error)
	GetTeam(context.Context, *dto.GetTeamRequest) (*dto.GetTeamResponse, error)
}

type teamHandler struct {
	teamService TeamService
}

func NewTeamHandler(teamService TeamService) *teamHandler {
	return &teamHandler{teamService: teamService}
}

func (h *teamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	var req dto.TeamAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		error_wrapper.WriteError(w, codes.INVALID_JSON, server.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	resp, err := h.teamService.AddTeam(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTeamAlreadyExists):
			error_wrapper.WriteError(w, codes.TEAM_EXISTS, server.ErrTeamAlreadyExists, http.StatusBadRequest)
		default:
			error_wrapper.WriteError(w, codes.INTERNAL_ERROR, server.ErrInternalError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
	return
}

func (h *teamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	var req dto.GetTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		error_wrapper.WriteError(w, codes.INVALID_JSON, server.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	resp, err := h.teamService.GetTeam(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTeamNotFound):
			error_wrapper.WriteError(w, codes.NOT_FOUND, server.ErrTeamNotFound, http.StatusNotFound)
		default:
			error_wrapper.WriteError(w, codes.INTERNAL_ERROR, server.ErrInternalError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	return
}
