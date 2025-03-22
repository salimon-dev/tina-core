package websocket

import (
	"encoding/json"
	"fmt"
	"salimon/nexus/types"

	"github.com/gorilla/websocket"
)

func sendPayload[T any](data T, conn *websocket.Conn) {
	responseStr, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = conn.WriteMessage(websocket.TextMessage, responseStr)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func sendInvalidPayload(conn *websocket.Conn) {
	sendErrorPayload("invalid payload", conn)
}

func sendErrorPayload(message string, conn *websocket.Conn) {
	resp := types.WsResponse{
		Action:  "ERROR",
		Message: message,
	}
	sendPayload(resp, conn)
}
