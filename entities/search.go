package entities

import (
	"fmt"
	"net/http"
	"salimon/nexus/db"
	"salimon/nexus/helpers"
	"salimon/nexus/types"
	"strconv"

	"github.com/labstack/echo/v4"
)

func SearchHandler(ctx echo.Context) error {
	pageStr := ctx.QueryParam("page")
	pageSizeStr := ctx.QueryParam("page_size")

	page := 1
	pageSize := 10

	if pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil {
			page = v
		}
	}
	if pageSizeStr != "" {
		if v, err := strconv.Atoi(pageSizeStr); err == nil {
			pageSize = v
		}
	}

	offset := (page - 1) * pageSize
	limit := pageSize

	var records []types.Entity
	results := db.EntityModel().Select("*").Offset(offset).Limit(limit).Find(&records)

	if results.Error != nil {
		fmt.Println(results.Error)
		return helpers.InternalError(ctx)
	}

	var count int64
	results = db.EntityModel().Select("*").Count(&count)

	if results.Error != nil {
		fmt.Println(results.Error)
		return helpers.InternalError(ctx)
	}

	user := ctx.Get("user").(*types.User)

	if user.Role >= types.UserRoleMember {
		// data for public
		type PublicEntity = struct {
			Name        string   `json:"name"`
			BaseUrl     string   `json:"base_url"`
			Description string   `json:"description"`
			Tags        []string `json:"tags"`
		}

		publicRecords := make([]PublicEntity, len(records))
		for i := 0; i < len(records); i++ {
			publicRecords[i] = PublicEntity{
				Name:        records[i].Name,
				BaseUrl:     records[i].BaseUrl,
				Description: records[i].Description,
				Tags:        records[i].Tags,
			}
		}

		data := types.CollectionResponse[PublicEntity]{
			Data:     publicRecords,
			Total:    count,
			PageSize: pageSize,
			Page:     page,
		}
		return ctx.JSON(http.StatusOK, data)

	} else {
		// full data for admin
		data := types.CollectionResponse[types.Entity]{
			Data:     records,
			Total:    count,
			PageSize: pageSize,
			Page:     page,
		}
		return ctx.JSON(http.StatusOK, data)
	}
}
