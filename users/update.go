package users

import (
	"crypto/md5"
	"encoding/hex"
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
	Username  string           `json:"username" validate:"required"`
	Password  *string          `json:"password,omitempty" validate:"omitempty,gte=5"`
	Status    types.UserStatus `json:"status" validate:"required"`
	Role      types.UserRole   `json:"role" validate:"required,numeric"`
	Credit    int32            `json:"credit" validate:"required,numeric"`
	SecretKey string           `json:"secret_key" validate:"required,lte=64"`
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

	record, err := db.FindUser("id = ?", id)

	if err != nil {
		fmt.Println(err)
		return helpers.InternalError(ctx)
	}

	if record == nil {
		return ctx.String(http.StatusNotFound, "not found")
	}

	if record.Username != payload.Username {
		existingUser, err := db.FindUser("username = ?", payload.Username)
		if err != nil {
			fmt.Println(err)
			return helpers.InternalError(ctx)
		}
		if existingUser != nil {
			return ctx.JSON(http.StatusBadRequest, helpers.MakeSingleValidationError("username", "username already exists"))
		}
	}

	if payload.Password != nil {
		passwordHash := md5.Sum([]byte(*payload.Password))
		password := hex.EncodeToString(passwordHash[:])
		record.Password = password
	}

	record.Username = payload.Username
	record.Status = payload.Status
	record.Role = payload.Role
	record.Credit = payload.Credit
	record.SecretKey = payload.SecretKey
	record.UpdatedAt = time.Now()

	result := db.UsersModel().Where("id = ?", id).Save(record)
	if result.Error != nil {
		fmt.Println(result.Error)
		return helpers.InternalError(ctx)
	}

	return ctx.JSON(http.StatusOK, record)
}
