package entity

import (
	"fmt"
	"net/http"
	"salimon/nexus/db"
	"salimon/nexus/helpers"
	"salimon/nexus/types"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetUserHandler(ctx echo.Context) error {
	userId := ctx.Param("userId")
	user, err := db.FindUser("id = ?", userId)
	if err != nil {
		fmt.Println(err)
		return helpers.InternalError(ctx)
	}
	if user == nil {
		return ctx.String(http.StatusNotFound, "user not found")
	}

	type userResponse struct {
		ID           uuid.UUID        `json:"id"`
		Username     string           `json:"username"`
		Role         types.UserRole   `json:"role"`
		Status       types.UserStatus `json:"status"`
		RegisteredAt time.Time        `json:"registered_at"`
	}

	resp := &userResponse{
		ID:           user.Id,
		Username:     user.Username,
		Role:         user.Role,
		Status:       user.Status,
		RegisteredAt: user.RegisteredAt,
	}
	return ctx.JSON(http.StatusOK, resp)
}
