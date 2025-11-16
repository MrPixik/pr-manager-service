package user

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

// mockgen -source="internal/http/server/handlers/user/user.go" -destination="internal/http/server/handlers/user/mocks/mock_user_service.go" -package=mocks UserService
type UserService interface {
	SetIsActive(ctx context.Context, req *dto.SetIsActiveRequest) (*dto.SetIsActiveResponse, error)
	GetReviewPullRequests(context.Context, *dto.GetReviewPRRequest) (*dto.GetReviewPRResponse, error)
}

type userHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *userHandler {
	return &userHandler{userService: userService}
}

func (h *userHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req dto.SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		error_wrapper.WriteError(w, codes.INVALID_JSON, server.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	resp, err := h.userService.SetIsActive(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			error_wrapper.WriteError(w, codes.NOT_FOUND, server.ErrUserNotFound, http.StatusNotFound)
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

func (h *userHandler) GetReviewPullRequests(w http.ResponseWriter, r *http.Request) {
	var req dto.GetReviewPRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		error_wrapper.WriteError(w, codes.INVALID_JSON, server.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	resp, err := h.userService.GetReviewPullRequests(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			error_wrapper.WriteError(w, codes.NOT_FOUND, server.ErrUserNotFound, http.StatusNotFound)
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
