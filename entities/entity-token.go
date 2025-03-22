package entities

import (
	"fmt"
	"net/http"

	"salimon/nexus/db"
	"salimon/nexus/helpers"
	"salimon/nexus/middlewares"

	"github.com/labstack/echo/v4"
)

type entityTokenSchema struct {
	Entity string `json:"entity" validate:"required"`
}

type entityTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func EntityTokenHandler(ctx echo.Context) error {
	payload := new(entityTokenSchema)
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

	entity, err := db.FindEntity("name = ?", payload.Entity)
	if err != nil {
		fmt.Println(err)
		return helpers.InternalError(ctx)
	}
	if entity == nil {
		return echo.NewHTTPError(http.StatusNotFound, "entity not found")
	}

	accessToken, refreshToken, err := helpers.GenerateJWTFromEntity(entity.Id, entity.SecretKey)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	response := entityTokenResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}

	return ctx.JSON(http.StatusOK, response)
}
