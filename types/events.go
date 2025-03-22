package types

import "github.com/gorilla/websocket"

type WsMessage struct {
	From string `json:"from"`
	Body string `json:"body"`
}

type WsEvent struct {
	Action      string `json:"action,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	Body        string `json:"body,omitempty"`
	Entity      string `json:"entity,omitempty"`
}

type WsResponse struct {
	Action  string `json:"action,omitempty"`
	Message string `json:"message,omitempty"`
	Result  *bool  `json:"result,omitempty"`
	Token   string `json:"token,omitempty"`
}

type WsContext struct {
	Conn   *websocket.Conn
	User   *User
	Entity *Entity
}
