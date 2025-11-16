package pull_request

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/service"
	"service-order-avito/internal/http/server/handlers/pull_request/mocks"
	"testing"
)

func TestPullRequestHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPullRequestService(ctrl)
	handler := NewPullRequestHandler(mockService)

	tests := []struct {
		name           string
		reqBody        interface{}
		mockReturnResp *dto.PullRequestCreateResponse
		mockReturnErr  error
		expectedCode   int
	}{
		{
			name: "success",
			reqBody: &dto.PullRequestCreateRequest{
				PullRequestID:   "pr1",
				PullRequestName: "Add feature",
				AuthorID:        "u1",
			},
			mockReturnResp: &dto.PullRequestCreateResponse{PullRequest: dto.PullRequestResponse{
				PullRequestID:   "pr1",
				PullRequestName: "Add feature",
				AuthorID:        "u1",
				Status:          "OPEN",
			}},
			mockReturnErr: nil,
			expectedCode:  http.StatusCreated,
		},
		{
			name: "user not found",
			reqBody: &dto.PullRequestCreateRequest{
				PullRequestID:   "pr2",
				PullRequestName: "Fix bug",
				AuthorID:        "u999",
			},
			mockReturnResp: nil,
			mockReturnErr:  service.ErrUserNotFound,
			expectedCode:   http.StatusNotFound,
		},
		{
			name:           "invalid json",
			reqBody:        "invalid json",
			mockReturnResp: nil,
			mockReturnErr:  nil,
			expectedCode:   http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte
			if str, ok := tt.reqBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.reqBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			if tt.name != "invalid json" {
				mockService.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(tt.mockReturnResp, tt.mockReturnErr)
			}

			handler.Create(w, req)

			assert.Equal(t, tt.expectedCode, w.Result().StatusCode)
		})
	}
}

func TestPullRequestHandler_Merge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPullRequestService(ctrl)
	handler := NewPullRequestHandler(mockService)

	reqBody := &dto.PullRequestMergeRequest{PullRequestID: "pr1"}
	respBody := &dto.PullRequestMergeResponse{PullRequest: dto.PullRequestMergedResponse{
		PullRequestID: "pr1",
		Status:        "MERGED",
	}}

	mockService.EXPECT().Merge(gomock.Any(), reqBody).Return(respBody, nil)

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/merge", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.Merge(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestPullRequestHandler_ReassignReviewer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockPullRequestService(ctrl)
	handler := NewPullRequestHandler(mockService)

	reqBody := &dto.PullRequestReassignRequest{
		PullRequestID: "pr1",
		OldReviewerID: "u2",
	}
	respBody := &dto.PullRequestReassignResponse{
		ReplacedBy: "u3",
	}

	mockService.EXPECT().
		ReassignReviewer(gomock.Any(), reqBody).
		Return(respBody, nil)

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/reassign", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.ReassignReviewer(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}
