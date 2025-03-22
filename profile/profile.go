package profile

import (
	"net/http"
	"salimon/nexus/types"

	"github.com/labstack/echo/v4"
)

func GetHandler(ctx echo.Context) error {
	user := ctx.Get("user").(*types.User)

	return ctx.JSON(http.StatusOK, user)
}
