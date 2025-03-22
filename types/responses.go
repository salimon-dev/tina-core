package types

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Data         User   `json:"data"`
}

type CollectionResponse[T interface{}] struct {
	Data     []T   `json:"data"`
	PageSize int   `json:"page_size"`
	Page     int   `json:"page"`
	Total    int64 `json:"total"`
}
