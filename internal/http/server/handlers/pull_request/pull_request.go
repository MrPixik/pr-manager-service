package pull_request

import (
	"context"
	"encoding/json"
	"net/http"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/server"
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
		error_wrapper.WriteServiceError(w, err)
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
		error_wrapper.WriteServiceError(w, err)
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
		error_wrapper.WriteServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	return
}
