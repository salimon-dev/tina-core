package middlewares

import (
	"net/http"
	"salimon/nexus/types"

	"github.com/labstack/echo/v4"
)

func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*types.User)
		if user.Role > types.UserRoleAdmin {
			return ctx.String(http.StatusForbidden, "permission denied")
		}
		return next(ctx)
	}
}
