package pull_request

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

// mockgen -source="internal/http/server/handlers/pull_request/pull_request.go" -destination="internal/http/server/handlers/pull_request/mocks/mock_pull_request_service.go" -package=mocks PullRequestService
type PullRequestService interface {
	Create(context.Context, *dto.PullRequestCreateRequest) (*dto.PullRequestCreateResponse, error)
	Merge(context.Context, *dto.PullRequestMergeRequest) (*dto.PullRequestMergeResponse, error)
	ReassignReviewer(context.Context, *dto.PullRequestReassignRequest) (*dto.PullRequestReassignResponse, error)
}

type pullRequestHandler struct {
	prService PullRequestService
}

func NewPullRequestHandler(pullRequestService PullRequestService) *pullRequestHandler {
	return &pullRequestHandler{prService: pullRequestService}
}

func (h *pullRequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.PullRequestCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		error_wrapper.WriteError(w, codes.INVALID_JSON, server.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	resp, err := h.prService.Create(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			error_wrapper.WriteError(w, codes.NOT_FOUND, server.ErrUserNotFound, http.StatusNotFound)
		case errors.Is(err, service.ErrTeamNotFound):
			error_wrapper.WriteError(w, codes.NOT_FOUND, server.ErrTeamNotFound, http.StatusNotFound)
		case errors.Is(err, service.ErrPullRequestExists):
			error_wrapper.WriteError(w, codes.PR_EXISTS, server.ErrPRAlreadyExists, http.StatusConflict)
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

func (h *pullRequestHandler) Merge(w http.ResponseWriter, r *http.Request) {
	var req dto.PullRequestMergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		error_wrapper.WriteError(w, codes.INVALID_JSON, server.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	resp, err := h.prService.Merge(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPullRequestNotFound):
			error_wrapper.WriteError(w, codes.NOT_FOUND, server.ErrPRNotFound, http.StatusNotFound)
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

func (h *pullRequestHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	var req dto.PullRequestReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		error_wrapper.WriteError(w, codes.INVALID_JSON, server.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	resp, err := h.prService.ReassignReviewer(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPullRequestNotFound):
			error_wrapper.WriteError(w, codes.NOT_FOUND, server.ErrPRNotFound, http.StatusNotFound)
		case errors.Is(err, service.ErrUserNotFound):
			error_wrapper.WriteError(w, codes.NOT_FOUND, server.ErrUserNotFound, http.StatusNotFound)
		case errors.Is(err, service.ErrPullRequestMerged):
			error_wrapper.WriteError(w, codes.PR_MERGED, server.ErrPullRequestMerged, http.StatusConflict)
		case errors.Is(err, service.ErrReviewerNotAssigned):
			error_wrapper.WriteError(w, codes.NOT_ASSIGNED, server.ErrReviewerNotAssigned, http.StatusConflict)
		case errors.Is(err, service.ErrNoReplacementCandidate):
			error_wrapper.WriteError(w, codes.NO_CANDIDATE, server.ErrNoReplacementCandidate, http.StatusConflict)
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
