package main

import (
	"net/http"
	"salimon/tina-core/middlewares"
	"salimon/tina-core/types"

	"github.com/labstack/echo/v4"
)

func InteractHandler(ctx echo.Context) error {
	payload := new(types.InteractSchema)
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

	// user := ctx.Get("user").(*types.User)

	data := payload.Data

	response := make([]types.Message, len(data)+1)
	for i := 0; i < len(data); i++ {
		message := types.Message{From: data[i].From, Body: data[i].Body}
		response[i] = message
	}
	response[len(response)-1] = types.Message{From: "tina", Body: "Hello!"}
	return ctx.JSON(http.StatusOK, response)
}
