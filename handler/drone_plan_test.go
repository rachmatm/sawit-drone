package handler

import (
	"net/http"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/models"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetEstateEstateIdDronePlan(t *testing.T) {
	estateID := uuid.New()

	makeEstate := func() *models.Estate {
		return &models.Estate{ID: estateID, Width: 5, Length: 1}
	}

	t.Run("success: no max-distance returns total distance", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockEstateRepository(ctrl)
		server := NewServer(mockRepo)

		mockRepo.EXPECT().GetEstate(gomock.Any(), estateID).Return(makeEstate(), nil)
		mockRepo.EXPECT().GetTreesByEstate(gomock.Any(), estateID).Return(nil, nil)

		ctx, rec := newEchoContext(t, http.MethodGet, "/estate/"+estateID.String()+"/drone-plan", nil)

		err := server.GetEstateEstateIdDronePlan(ctx, estateID,
			generated.GetEstateEstateIdDronePlanParams{MaxDistance: nil})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)

		var resp generated.DronePlanResponse
		decodeJSON(t, rec, &resp)
		require.Nil(t, resp.Rest)
		require.Greater(t, resp.Distance, 0)
	})

	t.Run("success: with max-distance returns rest point", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockEstateRepository(ctrl)
		server := NewServer(mockRepo)

		mockRepo.EXPECT().GetEstate(gomock.Any(), estateID).Return(makeEstate(), nil)
		mockRepo.EXPECT().GetTreesByEstate(gomock.Any(), estateID).Return(nil, nil)

		maxDist := 15
		ctx, rec := newEchoContext(t, http.MethodGet,
			"/estate/"+estateID.String()+"/drone-plan?max-distance=15", nil)

		err := server.GetEstateEstateIdDronePlan(ctx, estateID,
			generated.GetEstateEstateIdDronePlanParams{MaxDistance: &maxDist})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)

		var resp generated.DronePlanResponse
		decodeJSON(t, rec, &resp)
		require.NotNil(t, resp.Rest)
	})

	t.Run("error: max-distance <= 0 → 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		server := NewServer(repository.NewMockEstateRepository(ctrl))

		ctx, rec := newEchoContext(t, http.MethodGet,
			"/estate/"+estateID.String()+"/drone-plan?max-distance=0", nil)

		maxDist := 0
		err := server.GetEstateEstateIdDronePlan(ctx, estateID,
			generated.GetEstateEstateIdDronePlanParams{MaxDistance: &maxDist})
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("error: estate not found → 404", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := repository.NewMockEstateRepository(ctrl)
		server := NewServer(mockRepo)

		mockRepo.EXPECT().GetEstate(gomock.Any(), estateID).
			Return(nil, repository.ErrEstateNotFound)

		ctx, rec := newEchoContext(t, http.MethodGet, "/estate/"+estateID.String()+"/drone-plan", nil)

		err := server.GetEstateEstateIdDronePlan(ctx, estateID,
			generated.GetEstateEstateIdDronePlanParams{})
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}
