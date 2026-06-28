package handler

import (
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s *Server) PostEstateEstateIdTree(ctx echo.Context, estateId openapi_types.UUID) error {

	var req generated.CreateTreeRequest

	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request",
		)
	}

	tree, err := s.repo.CreateTree(
		ctx.Request().Context(),
		uuid.UUID(estateId),
		req.X,
		req.Y,
		req.Height,
	)

	if err != nil {

		switch err {

		case repository.ErrEstateNotFound:
			return ctx.JSON(
				http.StatusNotFound,
				generated.ErrorResponse{
					Message: "estate not found",
				},
			)

		case repository.ErrInvalidCoordinate:
			return ctx.JSON(
				http.StatusBadRequest,
				generated.ErrorResponse{
					Message: "invalid coordinate",
				},
			)

		case repository.ErrTreeAlreadyExists:
			return ctx.JSON(
				http.StatusConflict,
				generated.ErrorResponse{
					Message: "tree already exists",
				},
			)
		}

		return err
	}

	return ctx.JSON(
		http.StatusCreated,
		generated.CreateTreeResponse{
			Id: tree.ID,
		},
	)
}
