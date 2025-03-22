package websocket

import (
	"fmt"
	"salimon/nexus/db"
	"salimon/nexus/helpers"
	"salimon/nexus/types"

	"github.com/gorilla/websocket"
)

func handleAuth(accessToken string, ctx *types.WsContext) {
	if accessToken == "" {
		sendInvalidPayload(ctx.Conn)
	}

	claims, err := helpers.VerifyNexusJWT(accessToken)

	if err != nil {
		fmt.Println(err.Error())
		sendAuthResultPayload(false, ctx.Conn)
		return
	}

	if claims.Type != "access" {
		sendAuthResultPayload(false, ctx.Conn)
		return
	}

	user, err := db.FindUser("id = ?", claims.UserID)

	if err != nil {
		fmt.Println(err.Error())
		sendAuthResultPayload(false, ctx.Conn)
		return
	}

	ctx.User = user
	sendAuthResultPayload(user != nil, ctx.Conn)
}
func sendAuthResultPayload(result bool, conn *websocket.Conn) {
	resp := types.WsResponse{
		Action: "AUTH",
		Result: &result,
	}
	sendPayload(resp, conn)
}
