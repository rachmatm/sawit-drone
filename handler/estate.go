package handler

import (
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/labstack/echo/v4"
)

func (s *Server) PostEstate(ctx echo.Context) error {

	var req generated.CreateEstateRequest

	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request body",
		)
	}

	estate, err := s.repo.CreateEstate(ctx.Request().Context(), req.Width, req.Length)

	if err != nil {
		return err
	}

	resp := generated.CreateEstateResponse{Id: estate.ID}
	return ctx.JSON(http.StatusCreated, resp)
}
