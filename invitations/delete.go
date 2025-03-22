package invitations

import (
	"fmt"
	"net/http"

	"salimon/nexus/db"
	"salimon/nexus/helpers"
	"salimon/nexus/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func DeleteHandler(ctx echo.Context) error {
	idString := ctx.Param("id")
	id, err := uuid.Parse(idString)

	if err != nil {
		fmt.Println(err)
		return ctx.String(http.StatusNotFound, "not found")
	}

	record, err := db.FindInvitation("id = ?", id)

	if err != nil {
		fmt.Println(err)
		return helpers.InternalError(ctx)
	}

	if record == nil {
		return ctx.String(http.StatusNotFound, "not found")
	}

	db.InvitationsModel().Where("id = ?", id).Delete(&types.Invitation{})
	return ctx.JSON(http.StatusOK, "record deleted")
}
