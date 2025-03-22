package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"salimon/nexus/types"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(c echo.Context) error {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return err
	}
	defer conn.Close()

	ctx := types.WsContext{
		Conn: conn,
	}

	for {
		// Read message from WebSocket
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		messageStr := string(data)
		if messageStr == "ping" {
			err = conn.WriteMessage(messageType, []byte("pong"))
			if err != nil {
				fmt.Println("Error writing message:", err)
				break
			}
			continue
		}

		event, err := praseEvent(data)

		if err != nil {
			fmt.Println("Error parsing message")
			break
		}

		if event.Action == "AUTH" {
			handleAuth(event.AccessToken, &ctx)
			continue
		}
		if event.Action == "CONNECT" {
			handleConnect(event.Entity, &ctx)
		}
		if event.Action == "MESSAGE" {
			handleMessage(event.Body, &ctx)
			continue
		}
		err = conn.WriteMessage(messageType, data)
		if err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}

	return nil
}

func praseEvent(data []byte) (types.WsEvent, error) {
	var event types.WsEvent
	err := json.Unmarshal(data, &event)
	return event, err
}
