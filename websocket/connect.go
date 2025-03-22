package websocket

import (
	"fmt"
	"salimon/nexus/db"
	"salimon/nexus/types"
)

func handleConnect(entityName string, ctx *types.WsContext) {
	if entityName == "" {
		sendInvalidPayload(ctx.Conn)
		return
	}
	entity, err := db.FindEntity("name = ?", entityName)
	if err != nil {
		fmt.Println(err)
		sendErrorPayload("internal error", ctx.Conn)
		return
	}

	ctx.Entity = entity

	result := true
	payload := types.WsResponse{
		Action: "CONNECT",
		Result: &result,
	}
	sendPayload(payload, ctx.Conn)
}
