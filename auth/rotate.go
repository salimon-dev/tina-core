package auth

import (
	"fmt"
	"net/http"

	"salimon/nexus/db"
	"salimon/nexus/helpers"
	"salimon/nexus/middlewares"
	"salimon/nexus/types"

	"github.com/labstack/echo/v4"
)

type rotatePayloadSchema struct {
	Token string `json:"token" validate:"required"`
}

func RotateHandler(ctx echo.Context) error {
	payload := new(rotatePayloadSchema)
	if err := ctx.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}

	// validation errors
	vError, err := middlewares.ValidatePayload(*payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}
	if vError != nil {
		return ctx.JSON(http.StatusBadRequest, vError)
	}

	claims, err := helpers.VerifyNexusJWT(payload.Token)
	if err != nil {
		fmt.Println(err)
		return helpers.UnauthorizedError(ctx)
	}
	if claims == nil {
		return helpers.UnauthorizedError(ctx)
	}

	if claims.Type != "refresh" {
		return helpers.UnauthorizedError(ctx)
	}

	user, err := db.FindUser("id = ?", claims.UserID)
	if err != nil {
		fmt.Println(err)
		return helpers.InternalError(ctx)
	}
	if user == nil {
		return helpers.UnauthorizedError(ctx)
	}

	accessToken, refreshToken, err := helpers.GenerateNexusJwt(user)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	response := types.AuthResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
		Data:         *user,
	}

	return ctx.JSON(http.StatusOK, response)
}
