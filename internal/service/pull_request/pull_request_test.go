package pull_request

import (
	"context"
	"service-order-avito/internal/domain/errors/repository"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"service-order-avito/internal/domain"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/service/error_wrapper"
	"service-order-avito/internal/service/pull_request/mocks"
)

func TestPullRequestService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPullRequestRepository(ctrl)
	service := NewPullRequestService(mockRepo)

	req := &dto.PullRequestCreateRequest{
		PullRequestID:   "pr1",
		PullRequestName: "Fix bug",
		AuthorID:        "user1",
	}

	expectedDomain := domain.PullRequest{
		ID:       "pr1",
		Name:     "Fix bug",
		AuthorID: "user1",
		Status:   "OPEN",
	}

	now := time.Now()
	prWithReviewers := &domain.PullRequestWithReviewers{
		PullRequest: domain.PullRequest{
			ID:        "pr1",
			Name:      "Fix bug",
			AuthorID:  "user1",
			Status:    "OPEN",
			CreatedAt: now,
			MergedAt:  nil,
		},
		AssignedReviewers: []string{"rev1", "rev2"},
	}

	mockRepo.
		EXPECT().
		CreateWithReviewers(gomock.Any(), expectedDomain).
		Return(prWithReviewers, nil)

	resp, err := service.Create(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Equal(t, "pr1", resp.PullRequest.PullRequestID)
	require.Equal(t, "Fix bug", resp.PullRequest.PullRequestName)
	require.Equal(t, "user1", resp.PullRequest.AuthorID)
	require.Equal(t, "OPEN", resp.PullRequest.Status)
	require.Equal(t, []string{"rev1", "rev2"}, resp.PullRequest.AssignedReviewers)
}

func TestPullRequestService_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPullRequestRepository(ctrl)
	service := NewPullRequestService(mockRepo)

	repoErr := repository.ErrPullRequestExists

	mockRepo.
		EXPECT().
		CreateWithReviewers(gomock.Any(), gomock.Any()).
		Return(nil, repoErr)

	_, err := service.Create(context.Background(), &dto.PullRequestCreateRequest{})
	require.Error(t, err)
	require.Equal(t, error_wrapper.WrapRepositoryError(repoErr), err)
}

func TestPullRequestService_Merge_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPullRequestRepository(ctrl)
	service := NewPullRequestService(mockRepo)

	req := &dto.PullRequestMergeRequest{
		PullRequestID: "pr1",
	}

	merged := &domain.PullRequestWithReviewers{
		PullRequest: domain.PullRequest{
			ID:       "pr1",
			Name:     "Feature X",
			AuthorID: "user1",
			Status:   "MERGED",
		},
		AssignedReviewers: []string{"rev1"},
	}

	mockRepo.
		EXPECT().
		Merge(gomock.Any(), "pr1").
		Return(merged, nil)

	resp, err := service.Merge(context.Background(), req)
	require.NoError(t, err)

	require.Equal(t, "pr1", resp.PullRequest.PullRequestID)
	require.Equal(t, "Feature X", resp.PullRequest.PullRequestName)
	require.Equal(t, "user1", resp.PullRequest.AuthorID)
	require.Equal(t, "MERGED", resp.PullRequest.Status)
	require.Equal(t, []string{"rev1"}, resp.PullRequest.AssignedReviewers)
}

func TestPullRequestService_Merge_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPullRequestRepository(ctrl)
	service := NewPullRequestService(mockRepo)

	repoErr := repository.ErrPullRequestNotFound

	mockRepo.
		EXPECT().
		Merge(gomock.Any(), "pr1").
		Return(nil, repoErr)

	_, err := service.Merge(context.Background(),
		&dto.PullRequestMergeRequest{PullRequestID: "pr1"},
	)

	require.Error(t, err)
	require.Equal(t, error_wrapper.WrapRepositoryError(repoErr), err)
}

func TestPullRequestService_ReassignReviewer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPullRequestRepository(ctrl)
	service := NewPullRequestService(mockRepo)

	req := &dto.PullRequestReassignRequest{
		PullRequestID: "pr1",
		OldReviewerID: "rev_old",
	}

	expectedReviewer := &domain.Reviewer{ID: "rev_new"}

	mockRepo.
		EXPECT().
		ReassignReviewer(gomock.Any(), "pr1", "rev_old").
		Return(expectedReviewer, nil)

	resp, err := service.ReassignReviewer(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, "rev_new", resp.ReplacedBy)
}

func TestPullRequestService_ReassignReviewer_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPullRequestRepository(ctrl)
	service := NewPullRequestService(mockRepo)

	repoErr := repository.ErrNoReplacementCandidate

	mockRepo.
		EXPECT().
		ReassignReviewer(gomock.Any(), "pr1", "rev1").
		Return(nil, repoErr)

	_, err := service.ReassignReviewer(context.Background(),
		&dto.PullRequestReassignRequest{
			PullRequestID: "pr1",
			OldReviewerID: "rev1",
		},
	)

	require.Error(t, err)
	require.Equal(t, error_wrapper.WrapRepositoryError(repoErr), err)
}
