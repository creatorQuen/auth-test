package handlers

import (
	"auth-test/dto"
	"auth-test/infrastructure/services"
	"auth-test/pkg/logging"
	"github.com/labstack/echo/v4"
	"net/http"
)

type authUserHandler struct {
	authUserService services.AuthUserService
	log             *logging.Logger
}

func NewAuthUserHandler(authUserService services.AuthUserService, log *logging.Logger) *authUserHandler {
	return &authUserHandler{
		authUserService: authUserService,
		log:             log,
	}
}

func (a *authUserHandler) Register(ctx echo.Context) error {
	var req dto.AuthUserRegisterReq
	err := ctx.Bind(&req)
	if err != nil {
		a.log.Error("Request to bind ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = ctx.Validate(req)
	if err != nil {
		a.log.Error("Request to validate ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, err := a.authUserService.CreateAuthUser(ctx.Request().Context(), req)
	if err != nil {
		a.log.Error("authUserService.CreateAuthUser: ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusCreated, id)
}

func (a *authUserHandler) Login(ctx echo.Context) error {
	var req dto.AuthUserLoginReq
	err := ctx.Bind(&req)
	if err != nil {
		a.log.Error("Request to bind ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = ctx.Validate(req)
	if err != nil {
		a.log.Error("Request to validate ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	token, err := a.authUserService.GenerateToken(ctx.Request().Context(), req)
	if err != nil {
		a.log.Error("authUserService.GenerateToken: ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, token)
}
