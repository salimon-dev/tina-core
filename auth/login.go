package auth

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"

	"salimon/nexus/db"
	"salimon/nexus/helpers"
	"salimon/nexus/middlewares"
	"salimon/nexus/types"

	"github.com/labstack/echo/v4"
)

type loginSchema struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,gte=5"`
}

func LoginHandler(ctx echo.Context) error {
	payload := new(loginSchema)
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

	passwordHash := md5.Sum([]byte(payload.Password))
	password := hex.EncodeToString(passwordHash[:])

	// fetch user based on email of verfication
	user, err := db.FindUser("username = ? AND password = ?", payload.Username, password)
	if err != nil {
		fmt.Println(err.Error())
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
