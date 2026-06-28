package handler

import (
	"net/http"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/models"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestPostEstate(t *testing.T) {
	t.Run("success: creates estate and returns UUID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockEstateRepository(ctrl)
		server := NewServer(mockRepo)

		estateID := uuid.New()
		mockRepo.EXPECT().
			CreateEstate(gomock.Any(), 10, 20).
			Return(&models.Estate{ID: estateID, Width: 10, Length: 20}, nil)

		ctx, rec := newEchoContext(t, http.MethodPost, "/estate", map[string]int{
			"width": 10, "length": 20,
		})

		err := server.PostEstate(ctx)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, rec.Code)

		var resp generated.CreateEstateResponse
		decodeJSON(t, rec, &resp)
		require.Equal(t, estateID, resp.Id)
	})

	t.Run("error: invalid request body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		server := NewServer(repository.NewMockEstateRepository(ctrl))
		ctx, _ := newEchoContext(t, http.MethodPost, "/estate", "{not-json")

		err := server.PostEstate(ctx)
		require.Error(t, err)
		httpError, isHttpError := err.(*echo.HTTPError)
		require.True(t, isHttpError, "expected *echo.HTTPError, got %T", err)
		require.Equal(t, http.StatusBadRequest, httpError.Code)
	})
}
