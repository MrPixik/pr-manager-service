package error_wrapper

import (
	"service-order-avito/internal/domain/errors/repository"
	"service-order-avito/internal/domain/errors/service"
)

// сделал эту мапу, чтобы не передавать ошибки с уровня репозитория наверх к уровню контроллеров
// в принципе, мне кажется, можно было бы и передавать, но я захотел реализовать более чистую архитектуру, полностью изолировав контроллер от репозитория
var repoToServiceMap = map[string]error{
	repository.ErrTeamAlreadyExists.Error():      service.ErrTeamAlreadyExists,
	repository.ErrTeamNotFound.Error():           service.ErrTeamNotFound,
	repository.ErrUserNotFound.Error():           service.ErrUserNotFound,
	repository.ErrInternalError.Error():          service.ErrInternalError,
	repository.ErrPullRequestExists.Error():      service.ErrPullRequestExists,
	repository.ErrPullRequestNotFound.Error():    service.ErrPullRequestNotFound,
	repository.ErrPullRequestMerged.Error():      service.ErrPullRequestMerged,
	repository.ErrReviewerNotAssigned.Error():    service.ErrReviewerNotAssigned,
	repository.ErrNoReplacementCandidate.Error(): service.ErrNoReplacementCandidate,
}

// WrapRepositoryError возвращает ошибку сервиса по ошибке репозитория
func WrapRepositoryError(err error) error {
	if err == nil {
		return nil
	}

	if mapped, ok := repoToServiceMap[err.Error()]; ok {
		return mapped
	}

	return service.ErrInternalError
}
