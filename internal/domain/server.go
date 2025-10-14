package domain

type SuccessResponse struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Message string         `json:"message,omitempty"`
	Errors  map[string]any `json:"errors"`
}
