package websocket

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"salimon/nexus/types"
	"strings"
)

func handleMessage(body string, ctx *types.WsContext) {
	if body == "" {
		sendInvalidPayload(ctx.Conn)
		return
	}
	if ctx.Entity == nil {
		sendErrorPayload("no entity connected", ctx.Conn)
		return
	}

	data := fmt.Sprintf(`{"data": "%s"}`, body)
	resp, err := http.Post(ctx.Entity.BaseUrl+"/interact", "application/json", strings.NewReader(data))

	if err != nil {
		fmt.Println(err)
		sendErrorPayload("internal error", ctx.Conn)
		return
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	tokenStart := types.WsResponse{
		Action: "TOKEN_START",
	}
	tokenEnd := types.WsResponse{
		Action: "TOKEN_END",
	}

	sendPayload(tokenStart, ctx.Conn)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			sendErrorPayload("internal error", ctx.Conn)
			return
		}
		token := strings.TrimSpace(line)
		payload := types.WsResponse{
			Action: "TOKEN",
			Token:  token,
		}
		fmt.Println(token)
		sendPayload(payload, ctx.Conn)
	}
	sendPayload(tokenEnd, ctx.Conn)
}
