package handler

import (
	"errors"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s *Server) GetEstateEstateIdDronePlan(
	ctx echo.Context,
	estateId openapi_types.UUID,
	params generated.GetEstateEstateIdDronePlanParams,
) error {

	if params.MaxDistance != nil &&
		*params.MaxDistance <= 0 {

		return ctx.JSON(
			http.StatusBadRequest,
			generated.ErrorResponse{
				Message: "invalid max-distance",
			},
		)
	}

	estate, err := s.repo.GetEstate(
		ctx.Request().Context(),
		uuid.UUID(estateId),
	)

	if err != nil {

		if errors.Is(
			err,
			repository.ErrEstateNotFound,
		) {

			return ctx.JSON(
				http.StatusNotFound,
				generated.ErrorResponse{
					Message: "estate not found",
				},
			)
		}

		return err
	}

	trees, err := s.repo.GetTreesByEstate(
		ctx.Request().Context(),
		uuid.UUID(estateId),
	)
	if err != nil {
		return err
	}

	plan := service.CalculateDronePlan(
		estate,
		trees,
		params.MaxDistance,
	)

	resp := generated.DronePlanResponse{
		Distance: plan.Distance,
	}

	if params.MaxDistance != nil {

		resp.Rest = &generated.LandingPoint{
			X: plan.RestX,
			Y: plan.RestY,
		}
	}

	return ctx.JSON(
		http.StatusOK,
		resp,
	)
}
