package invitations

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

	type Record struct {
		types.Invitation
		Username string `json:"username"`
	}

	var records []Record
	results := db.InvitationsModel().Select("invitations.*, users.username").Joins("JOIN users ON users.id = invitations.created_by").Offset(offset).Limit(limit).Find(&records)

	if results.Error != nil {
		fmt.Println(results.Error)
		return helpers.InternalError(ctx)
	}

	var count int64
	results = db.InvitationsModel().Select("*").Count(&count)

	if results.Error != nil {
		fmt.Println(results.Error)
		return helpers.InternalError(ctx)
	}

	data := types.CollectionResponse[Record]{
		Data:     records,
		Total:    count,
		PageSize: pageSize,
		Page:     page,
	}

	return ctx.JSON(http.StatusOK, data)
}
