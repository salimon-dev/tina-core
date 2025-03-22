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

type createSchema struct {
	Name        string                 `json:"name" validate:"required,lte=32"`
	Description string                 `json:"description" validate:"required,lte=32"`
	BaseUrl     string                 `json:"base_url" validate:"required,lte=256,url"`
	Tags        pq.StringArray         `json:"tags" validate:"required"`
	Credit      *int32                 `json:"credit" validate:"required"`
	Status      types.EntityStatus     `json:"status" validate:"required"`
	Permission  types.EntityPermission `json:"permission" validate:"required"`
}

func CreateHandler(ctx echo.Context) error {
	payload := new(createSchema)
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

	existingEntity, err := db.FindEntity("name = ?", payload.Name)
	if err != nil {
		fmt.Println(err)
		return helpers.InternalError(ctx)
	}
	if existingEntity != nil {
		return ctx.JSON(http.StatusBadRequest, helpers.MakeSingleValidationError("name", "entity with same name already exists in network"))
	}

	secretKey := helpers.GenerateRandomString(64)

	record := types.Entity{
		Id:          uuid.New(),
		Name:        payload.Name,
		Description: payload.Description,
		BaseUrl:     payload.BaseUrl,
		Credit:      *payload.Credit,
		Status:      payload.Status,
		Tags:        payload.Tags,
		Permission:  payload.Permission,
		SecretKey:   secretKey,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result := db.EntityModel().Create(record)
	if result.Error != nil {
		fmt.Println(result.Error)
		return helpers.InternalError(ctx)
	}
	return ctx.JSON(http.StatusOK, record)
}
