package invitations

import (
	"fmt"
	"net/http"
	"time"

	"salimon/nexus/db"
	"salimon/nexus/helpers"
	"salimon/nexus/middlewares"
	"salimon/nexus/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type updateSchema struct {
	Status         types.InvitationStatus `json:"status" validate:"required"`
	UsageRemaining int16                  `json:"usage_remaining" validate:"required,gte=1"`
	ExpiresAt      time.Time              `json:"expires_at" validate:"required"`
}

func UpdateHandler(ctx echo.Context) error {
	idString := ctx.Param("id")

	id, err := uuid.Parse(idString)
	if err != nil {
		fmt.Println(err)
		return ctx.String(http.StatusNotFound, "not found")
	}
	payload := new(updateSchema)
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

	record, err := db.FindInvitation("id = ?", id)

	if err != nil {
		fmt.Println(err)
		return helpers.InternalError(ctx)
	}

	if record == nil {
		return ctx.String(http.StatusNotFound, "not found")
	}

	record.Status = payload.Status
	record.ExpiresAt = payload.ExpiresAt
	record.UsageRemaining = payload.UsageRemaining
	record.UpdatedAt = time.Now()

	result := db.InvitationsModel().Where("id = ?", id).Save(record)
	if result.Error != nil {
		fmt.Println(result.Error)
		return helpers.InternalError(ctx)
	}
	return ctx.JSON(http.StatusOK, record)
}
