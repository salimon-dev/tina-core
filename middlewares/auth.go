package middlewares

import (
	"fmt"
	"salimon/tina-core/db"
	"salimon/tina-core/helpers"
	"salimon/tina-core/nexus"
	"salimon/tina-core/types"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authorization := ctx.Request().Header.Get("Authorization")

		// Check if the header is empty or doesn't start with "Bearer "
		if authorization == "" || !strings.HasPrefix(authorization, "Bearer ") {
			return helpers.UnauthorizedError(ctx)
		}

		// Extract the token part from the header
		token := strings.TrimPrefix(authorization, "Bearer ")

		claims, err := helpers.VerifyJWT(token)

		if err != nil || claims == nil {
			return helpers.UnauthorizedError(ctx)
		}

		if claims.Type != "access" {
			return helpers.UnauthorizedError(ctx)
		}

		user, err := db.FindUser("network_id = ?", claims.UserID)
		if err != nil {
			fmt.Println(err)
			return helpers.InternalError(ctx)
		}

		if user != nil {
			ctx.Set("user", user)
			return next(ctx)
		}

		fmt.Printf("user with id %s not found, fetching from nexus\n", claims.UserID)
		userData, err := nexus.FetchUserData(claims.UserID)
		if err != nil {
			fmt.Println(err)
			return helpers.UnauthorizedError(ctx)
		}
		user = &types.User{
			Id:           uuid.New(),
			Network:      "official",
			NetworkId:    userData.Id,
			Username:     userData.Username,
			Status:       userData.Status,
			Role:         userData.Role,
			RegisteredAt: userData.RegisteredAt,
			CreateAt:     time.Now(),
			UpdatedAt:    time.Now(),
		}
		err = db.InsertUser(user)
		if err != nil {
			fmt.Println(err)
			return helpers.UnauthorizedError(ctx)
		}

		ctx.Set("user", user)

		// return helpers.UnauthorizedError(ctx)
		return next(ctx)

	}
}
