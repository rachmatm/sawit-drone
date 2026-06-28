package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/models"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// GetEstateEstateIdStats
// ---------------------------------------------------------------------------

func TestGetEstateEstateIdStats(t *testing.T) {
	estateID := uuid.New()

	makeEstate := func() *models.Estate {
		return &models.Estate{ID: estateID, Width: 5, Length: 5}
	}

	t.Run("success: estate with odd number of trees returns correct stats", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockEstateRepository(ctrl)
		server := NewServer(mockRepo)

		// trees must be sorted by height ASC (as the real query does)
		trees := []models.Tree{
			{Height: 10},
			{Height: 20},
			{Height: 30},
		}
		mockRepo.EXPECT().GetEstate(gomock.Any(), estateID).Return(makeEstate(), nil)
		mockRepo.EXPECT().GetTreesByEstate(gomock.Any(), estateID).Return(trees, nil)

		ctx, rec := newEchoContext(t, http.MethodGet, "/estate/"+estateID.String()+"/stats", nil)

		err := server.GetEstateEstateIdStats(ctx, estateID)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)

		var resp generated.EstateStatsResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.Equal(t, 3, resp.Count)
		require.Equal(t, 10, resp.Min)
		require.Equal(t, 30, resp.Max)
		require.Equal(t, 20, resp.Median) // middle of [10, 20, 30]
	})

	t.Run("success: estate with even number of trees returns floor-average median", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockEstateRepository(ctrl)
		server := NewServer(mockRepo)

		trees := []models.Tree{
			{Height: 10},
			{Height: 20},
			{Height: 30},
			{Height: 40},
		}
		mockRepo.EXPECT().GetEstate(gomock.Any(), estateID).Return(makeEstate(), nil)
		mockRepo.EXPECT().GetTreesByEstate(gomock.Any(), estateID).Return(trees, nil)

		ctx, rec := newEchoContext(t, http.MethodGet, "/estate/"+estateID.String()+"/stats", nil)

		err := server.GetEstateEstateIdStats(ctx, estateID)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)

		var resp generated.EstateStatsResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.Equal(t, 4, resp.Count)
		require.Equal(t, 10, resp.Min)
		require.Equal(t, 40, resp.Max)
		require.Equal(t, 25, resp.Median) // (20+30)/2
	})

	t.Run("success: estate with a single tree — min/max/median are all equal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockEstateRepository(ctrl)
		server := NewServer(mockRepo)

		trees := []models.Tree{
			{Height: 15},
		}
		mockRepo.EXPECT().GetEstate(gomock.Any(), estateID).Return(makeEstate(), nil)
		mockRepo.EXPECT().GetTreesByEstate(gomock.Any(), estateID).Return(trees, nil)

		ctx, rec := newEchoContext(t, http.MethodGet, "/estate/"+estateID.String()+"/stats", nil)

		err := server.GetEstateEstateIdStats(ctx, estateID)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)

		var resp generated.EstateStatsResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.Equal(t, 1, resp.Count)
		require.Equal(t, 15, resp.Min)
		require.Equal(t, 15, resp.Max)
		require.Equal(t, 15, resp.Median)
	})

	t.Run("success: estate with no trees returns all zeros", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockEstateRepository(ctrl)
		server := NewServer(mockRepo)

		mockRepo.EXPECT().GetEstate(gomock.Any(), estateID).Return(makeEstate(), nil)
		mockRepo.EXPECT().GetTreesByEstate(gomock.Any(), estateID).Return(nil, nil)

		ctx, rec := newEchoContext(t, http.MethodGet, "/estate/"+estateID.String()+"/stats", nil)

		err := server.GetEstateEstateIdStats(ctx, estateID)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)

		var resp generated.EstateStatsResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.Equal(t, 0, resp.Count)
		require.Equal(t, 0, resp.Min)
		require.Equal(t, 0, resp.Max)
		require.Equal(t, 0, resp.Median)
	})

	t.Run("not found: GetEstate returns nil estate → 404", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockEstateRepository(ctrl)
		server := NewServer(mockRepo)

		// nil estate with no error is the "not found" signal in this codebase
		mockRepo.EXPECT().GetEstate(gomock.Any(), estateID).Return(nil, nil)
		// GetTreesByEstate must NOT be called when estate is not found
		mockRepo.EXPECT().GetTreesByEstate(gomock.Any(), gomock.Any()).Times(0)

		ctx, rec := newEchoContext(t, http.MethodGet, "/estate/"+estateID.String()+"/stats", nil)

		err := server.GetEstateEstateIdStats(ctx, estateID)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, rec.Code)

		var resp generated.ErrorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.NotEmpty(t, resp.Message)
	})

	t.Run("error: GetEstate returns a repository error → propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockEstateRepository(ctrl)
		server := NewServer(mockRepo)

		dbErr := errors.New("connection refused")
		mockRepo.EXPECT().GetEstate(gomock.Any(), estateID).Return(nil, dbErr)
		mockRepo.EXPECT().GetTreesByEstate(gomock.Any(), gomock.Any()).Times(0)

		ctx, rec := newEchoContext(t, http.MethodGet, "/estate/"+estateID.String()+"/stats", nil)
		_ = rec

		err := server.GetEstateEstateIdStats(ctx, estateID)
		require.ErrorIs(t, err, dbErr)
	})

	t.Run("error: GetTreesByEstate returns a repository error → propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockEstateRepository(ctrl)
		server := NewServer(mockRepo)

		dbErr := errors.New("query timeout")
		mockRepo.EXPECT().GetEstate(gomock.Any(), estateID).Return(makeEstate(), nil)
		mockRepo.EXPECT().GetTreesByEstate(gomock.Any(), estateID).Return(nil, dbErr)

		ctx, rec := newEchoContext(t, http.MethodGet, "/estate/"+estateID.String()+"/stats", nil)
		_ = rec

		err := server.GetEstateEstateIdStats(ctx, estateID)
		require.ErrorIs(t, err, dbErr)
	})
}
