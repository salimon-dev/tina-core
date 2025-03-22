package types

type Message struct {
	From string `json:"from" validate:"required"`
	Body string `json:"body" validate:"required"`
}

type InteractSchema struct {
	Data []Message `json:"data" validate:"required"`
}
