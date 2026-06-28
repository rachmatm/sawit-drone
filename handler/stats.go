package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/service"
)

func (s *Server) GetEstateEstateIdStats(ctx echo.Context, estateId openapi_types.UUID) error {

	estate, err := s.repo.GetEstate(
		ctx.Request().Context(),
		uuid.UUID(estateId),
	)
	if err != nil {
		return err
	}

	if estate == nil {
		return ctx.JSON(
			http.StatusNotFound,
			generated.ErrorResponse{
				Message: "estate not found",
			},
		)
	}

	trees, err := s.repo.GetTreesByEstate(
		ctx.Request().Context(),
		uuid.UUID(estateId),
	)
	if err != nil {
		return err
	}

	stats := service.CalculateStats(trees)

	return ctx.JSON(http.StatusOK, generated.EstateStatsResponse{
		Count:  stats.Count,
		Min:    stats.Min,
		Max:    stats.Max,
		Median: stats.Median,
	})
}
