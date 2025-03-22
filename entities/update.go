package entities

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
	"github.com/lib/pq"
)

type updateSchema struct {
	Name        string                 `json:"name" validate:"required,lte=32"`
	Description string                 `json:"description" validate:"required,lte=32"`
	BaseUrl     string                 `json:"base_url" validate:"required,lte=256,url"`
	Tags        pq.StringArray         `json:"tags" validate:"required"`
	Credit      *int32                 `json:"credit" validate:"required"`
	Status      types.EntityStatus     `json:"status" validate:"required"`
	Permission  types.EntityPermission `json:"permission" validate:"required"`
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

	record, err := db.FindEntity("id = ?", id)

	if err != nil {
		fmt.Println(err)
		return helpers.InternalError(ctx)
	}

	if record == nil {
		return ctx.String(http.StatusNotFound, "not found")
	}

	if record.Name != payload.Name {
		existingEntity, err := db.FindEntity("name = ?", payload.Name)
		if err != nil {
			fmt.Println(err)
			return helpers.InternalError(ctx)
		}
		if existingEntity != nil {
			return ctx.JSON(http.StatusBadRequest, helpers.MakeSingleValidationError("name", "entity with same name already exists in network"))
		}
	}

	record.Name = payload.Name
	record.Description = payload.Description
	record.Credit = *payload.Credit
	record.BaseUrl = payload.BaseUrl
	record.Tags = payload.Tags
	record.Status = payload.Status
	record.Permission = payload.Permission
	record.UpdatedAt = time.Now()

	result := db.EntityModel().Where("id = ?", id).Save(record)
	if result.Error != nil {
		fmt.Println(result.Error)
		return helpers.InternalError(ctx)
	}
	return ctx.JSON(http.StatusOK, record)
}
