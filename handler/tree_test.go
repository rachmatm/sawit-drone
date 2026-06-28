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

func TestPostEstateEstateIdTree(t *testing.T) {
	estateID := uuid.New()
	treeID := uuid.New()

	t.Run("success: creates tree and returns UUID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockEstateRepository(ctrl)
		server := NewServer(mockRepo)

		mockRepo.EXPECT().
			CreateTree(gomock.Any(), estateID, 3, 4, 15).
			Return(&models.Tree{ID: treeID, EstateID: estateID, X: 3, Y: 4, Height: 15}, nil)

		ctx, rec := newEchoContext(t, http.MethodPost, "/estate/"+estateID.String()+"/tree",
			map[string]int{"x": 3, "y": 4, "height": 15})

		err := server.PostEstateEstateIdTree(ctx, estateID)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, rec.Code)

		var resp generated.CreateTreeResponse
		decodeJSON(t, rec, &resp)
		require.Equal(t, treeID, resp.Id)
	})

	t.Run("error: invalid request body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockEstateRepository(ctrl)
		server := NewServer(mockRepo)

		ctx, _ := newEchoContext(t, http.MethodPost, "/estate/"+estateID.String()+"/tree", "{bad json")

		err := server.PostEstateEstateIdTree(ctx, estateID)
		require.Error(t, err)
		he, _ := err.(*echo.HTTPError)
		require.Equal(t, http.StatusBadRequest, he.Code)
	})

	t.Run("error: estate not found → 404", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockEstateRepository(ctrl)
		server := NewServer(mockRepo)

		mockRepo.EXPECT().
			CreateTree(gomock.Any(), estateID, 1, 1, 5).
			Return(nil, repository.ErrEstateNotFound)

		ctx, rec := newEchoContext(t, http.MethodPost, "/estate/"+estateID.String()+"/tree",
			map[string]int{"x": 1, "y": 1, "height": 5})

		err := server.PostEstateEstateIdTree(ctx, estateID)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("error: invalid coordinate → 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockEstateRepository(ctrl)
		server := NewServer(mockRepo)

		mockRepo.EXPECT().
			CreateTree(gomock.Any(), estateID, 0, 0, 5).
			Return(nil, repository.ErrInvalidCoordinate)

		ctx, rec := newEchoContext(t, http.MethodPost, "/estate/"+estateID.String()+"/tree",
			map[string]int{"x": 0, "y": 0, "height": 5})

		err := server.PostEstateEstateIdTree(ctx, estateID)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("error: duplicate tree → 409", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockEstateRepository(ctrl)
		server := NewServer(mockRepo)

		mockRepo.EXPECT().
			CreateTree(gomock.Any(), estateID, 2, 2, 10).
			Return(nil, repository.ErrTreeAlreadyExists)

		ctx, rec := newEchoContext(t, http.MethodPost, "/estate/"+estateID.String()+"/tree",
			map[string]int{"x": 2, "y": 2, "height": 10})

		err := server.PostEstateEstateIdTree(ctx, estateID)
		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, rec.Code)
	})
}
