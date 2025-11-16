package team

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/service"
	"service-order-avito/internal/http/server/handlers/team/mocks"
	"testing"
)

func TestTeamHandler_AddTeam(t *testing.T) {
	tests := []struct {
		name           string
		reqBody        interface{}
		mockReturnResp *dto.AddTeamResponse
		mockReturnErr  error
		expectedCode   int
	}{
		{
			name: "success",
			reqBody: &dto.TeamAddRequest{
				TeamName: "team1",
			},
			mockReturnResp: &dto.AddTeamResponse{
				TeamName: "team1",
			},
			mockReturnErr: nil,
			expectedCode:  http.StatusCreated,
		},
		{
			name: "team already exists",
			reqBody: &dto.TeamAddRequest{
				TeamName: "team2",
			},
			mockReturnResp: nil,
			mockReturnErr:  service.ErrTeamAlreadyExists,
			expectedCode:   http.StatusBadRequest,
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockTeamService(ctrl)
			handler := NewTeamHandler(mockService)

			var bodyBytes []byte
			if str, ok := tt.reqBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.reqBody)
			}

			if tt.name != "invalid json" {
				mockService.EXPECT().
					AddTeam(gomock.Any(), gomock.Any()).
					Return(tt.mockReturnResp, tt.mockReturnErr)
			}

			req := httptest.NewRequest(http.MethodPost, "/add_team", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			handler.AddTeam(w, req)

			assert.Equal(t, tt.expectedCode, w.Result().StatusCode)
		})
	}
}

func TestTeamHandler_GetTeam(t *testing.T) {
	tests := []struct {
		name           string
		reqBody        interface{}
		mockReturnResp *dto.GetTeamResponse
		mockReturnErr  error
		expectedCode   int
	}{
		{
			name: "success",
			reqBody: &dto.GetTeamRequest{
				TeamName: "team1",
			},
			mockReturnResp: &dto.GetTeamResponse{
				TeamName: "team1",
				Members: []dto.TeamMemberResponse{
					{UserID: "u1", Username: "Alice", IsActive: true},
				},
			},
			mockReturnErr: nil,
			expectedCode:  http.StatusOK,
		},
		{
			name: "team not found",
			reqBody: &dto.GetTeamRequest{
				TeamName: "team2",
			},
			mockReturnResp: nil,
			mockReturnErr:  service.ErrTeamNotFound,
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockTeamService(ctrl)
			handler := NewTeamHandler(mockService)

			var bodyBytes []byte
			if str, ok := tt.reqBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.reqBody)
			}

			if tt.name != "invalid json" {
				mockService.EXPECT().
					GetTeam(gomock.Any(), gomock.Any()).
					Return(tt.mockReturnResp, tt.mockReturnErr)
			}

			req := httptest.NewRequest(http.MethodPost, "/get_team", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			handler.GetTeam(w, req)

			assert.Equal(t, tt.expectedCode, w.Result().StatusCode)
		})
	}
}
