package error_wrapper

import (
	"encoding/json"
	"net/http"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/server"
	"service-order-avito/internal/domain/errors/service"
	"service-order-avito/internal/http/codes"
)

type errorMeta struct {
	Code    string
	Message string
	Status  int
}

var serviceErrorMap = map[error]errorMeta{
	service.ErrInternalError:          {codes.INTERNAL_ERROR, server.ErrInternalError, http.StatusInternalServerError},
	service.ErrTeamAlreadyExists:      {codes.TEAM_EXISTS, server.ErrTeamAlreadyExists, http.StatusBadRequest},
	service.ErrTeamNotFound:           {codes.NOT_FOUND, server.ErrTeamNotFound, http.StatusNotFound},
	service.ErrUserNotFound:           {codes.NOT_FOUND, server.ErrUserNotFound, http.StatusNotFound},
	service.ErrPullRequestExists:      {codes.PR_EXISTS, server.ErrPRAlreadyExists, http.StatusConflict},
	service.ErrPullRequestNotFound:    {codes.NOT_FOUND, server.ErrPRNotFound, http.StatusNotFound},
	service.ErrPullRequestMerged:      {codes.PR_MERGED, server.ErrPullRequestMerged, http.StatusBadRequest},
	service.ErrReviewerNotAssigned:    {codes.NOT_ASSIGNED, server.ErrReviewerNotAssigned, http.StatusBadRequest},
	service.ErrNoReplacementCandidate: {codes.NO_CANDIDATE, server.ErrNoReplacementCandidate, http.StatusBadRequest},
}

// WriteServiceError принимает ошибку уровня service и пишет ошибку уровня контроллера в ResponseWriter
func WriteServiceError(w http.ResponseWriter, err error) {
	meta, ok := serviceErrorMap[err]
	if !ok {
		meta = serviceErrorMap[service.ErrInternalError]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(meta.Status)
	_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
		Error: dto.ErrorDetail{
			Code:    meta.Code,
			Message: meta.Message,
		},
	})
}

func WriteError(w http.ResponseWriter, code string, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
		Error: dto.ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}
