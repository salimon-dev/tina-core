package auth

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

type registerPayloadSchema struct {
	InvitationCode string `json:"invitation_code" validate:"required,lte=16"`
	Username       string `json:"username" validate:"required"`
	Password       string `json:"password" validate:"required,gte=5"`
}

func RegisterHandler(ctx echo.Context) error {
	payload := new(registerPayloadSchema)
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

	invitation, err := db.FindInvitation("code = ? AND usage_remaining > 0 AND expires_at > ? AND status = ?", payload.InvitationCode, time.Now(), types.InvitationStatusActive)

	if err != nil {
		fmt.Println(err.Error())
		return helpers.InternalError(ctx)
	}

	if invitation == nil {
		return ctx.JSON(http.StatusBadRequest, helpers.MakeSingleValidationError("invitation_code", "invitation_code is invalid"))
	}

	passwordHash := md5.Sum([]byte(payload.Password))
	password := hex.EncodeToString(passwordHash[:])

	user, err := db.FindUser("username = ?", payload.Username)
	if err != nil {
		fmt.Println(err)
		return helpers.InternalError(ctx)
	}

	if user != nil {
		// user is registered already
		return ctx.JSON(http.StatusBadRequest, helpers.MakeSingleValidationError("action", "a user with same username already exists"))
	}
	user = &types.User{
		Id:           uuid.New(),
		Username:     payload.Username,
		Password:     password,
		InvitationId: invitation.Id,
		Credit:       0,
		SecretKey:    helpers.GenerateRandomString(64),
		Role:         types.UserRoleMember,
		Status:       types.UserStatusActive,
		RegisteredAt: time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = db.InsertUser(user)
	if err != nil {
		fmt.Println(err)
		return helpers.InternalError(ctx)
	}

	invitation.UsageRemaining -= 1
	result := db.InvitationsModel().Where("id = ?", invitation.Id).Updates(invitation)

	if result.Error != nil {
		fmt.Println(result.Error)
		return helpers.InternalError(ctx)
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
