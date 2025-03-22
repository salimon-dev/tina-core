package rest

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type heartBeatResponse struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
	Time        int64  `json:"time"`
}

func HeartBeatHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, heartBeatResponse{
		Name:        "nexus service",
		Environment: os.Getenv("ENV"),
		Time:        time.Now().Unix(),
	})
}
