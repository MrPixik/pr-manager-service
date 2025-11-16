package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"service-order-avito/internal/domain/dto"
	"service-order-avito/internal/domain/errors/service"
	"service-order-avito/internal/http/server/handlers/user/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_SetIsActive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockService)

	t.Run("success", func(t *testing.T) {
		reqBody := &dto.SetIsActiveRequest{
			UserID:   "u1",
			IsActive: true,
		}
		respBody := &dto.SetIsActiveResponse{
			User: dto.UserResponse{
				UserID:   "u1",
				Username: "John",
				TeamName: "teamA",
				IsActive: true,
			},
		}

		mockService.EXPECT().
			SetIsActive(gomock.Any(), reqBody).
			Return(respBody, nil)

		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/set-active", bytes.NewReader(bodyBytes))
		w := httptest.NewRecorder()

		handler.SetIsActive(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var resp dto.SetIsActiveResponse
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.Equal(t, respBody.User.UserID, resp.User.UserID)
		assert.Equal(t, respBody.User.IsActive, resp.User.IsActive)
	})

	t.Run("user not found", func(t *testing.T) {
		reqBody := &dto.SetIsActiveRequest{
			UserID:   "unknown",
			IsActive: true,
		}

		mockService.EXPECT().
			SetIsActive(gomock.Any(), reqBody).
			Return(nil, service.ErrUserNotFound)

		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/set-active", bytes.NewReader(bodyBytes))
		w := httptest.NewRecorder()

		handler.SetIsActive(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/set-active", bytes.NewReader([]byte("{invalid-json")))
		w := httptest.NewRecorder()

		handler.SetIsActive(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("internal error", func(t *testing.T) {
		reqBody := &dto.SetIsActiveRequest{
			UserID:   "u1",
			IsActive: true,
		}

		mockService.EXPECT().
			SetIsActive(gomock.Any(), reqBody).
			Return(nil, errors.New("some error"))

		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/set-active", bytes.NewReader(bodyBytes))
		w := httptest.NewRecorder()

		handler.SetIsActive(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})
}
