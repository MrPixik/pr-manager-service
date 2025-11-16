package team

import (
	"context"
	"encoding/json"
	"net/http"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/server"
	"service-order-avito/internal/http/codes"
	"service-order-avito/pkg/http/error_wrapper"
)

// mockgen -source="internal/http/server/handlers/team/team.go" -destination="internal/http/server/handlers/team/mocks/mock_user_service.go" -package=mocks TeamService
type TeamService interface {
	AddTeam(context.Context, *dto.TeamAddRequest) (*dto.AddTeamResponse, error)
	GetTeam(context.Context, *dto.GetTeamRequest) (*dto.GetTeamResponse, error)
	GetTeamStats(context.Context, *dto.GetTeamStatsRequest) (*dto.TeamStatsResponse, error)
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
		error_wrapper.WriteServiceError(w, err)
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
		error_wrapper.WriteServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	return
}

func (h *teamHandler) GetTeamStats(w http.ResponseWriter, r *http.Request) {
	var req dto.GetTeamStatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		error_wrapper.WriteError(w, codes.INVALID_JSON, server.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	resp, err := h.teamService.GetTeamStats(r.Context(), &req)
	if err != nil {
		error_wrapper.WriteServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	return
}
