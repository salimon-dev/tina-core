package e2e

import (
	"net/http"
	"salimon/nexus/middlewares"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type interactSchema struct {
	Data string `json:"data" validate:"required"`
}

func InteractHandler(ctx echo.Context) error {
	payload := new(interactSchema)
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

	tokens := strings.Fields(payload.Data)

	ctx.Response().Header().Set("Content-Type", "text/plain")
	ctx.Response().WriteHeader(http.StatusOK)

	for _, token := range tokens {
		ctx.Response().Write([]byte(token + "\n"))
		ctx.Response().Flush()
		time.Sleep(50 * time.Millisecond)
	}

	return nil
}
