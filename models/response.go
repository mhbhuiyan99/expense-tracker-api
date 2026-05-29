package models

// SuccessResponse: standard shape for all successfull API responses
type SuccessResponse struct {
	Status string `json:"status"`
	Data interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse: standard shape for all error API responses
type ErrorResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
	Code string `json:"code,omitempty"`
}
