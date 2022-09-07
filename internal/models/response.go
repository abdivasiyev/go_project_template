package models

type SuccessResponse struct {
	Ok bool `json:"ok"`
}

type ErrorResponse struct {
	ErrorCode    int                  `json:"error_code" example:"400"`
	ErrorMessage string               `json:"error_message" example:"Bad Request"`
	Validations  []ValidationResponse `json:"validations,omitempty"`
}

type ValidationResponse struct {
	Field string `json:"field" example:"login"`
	Error string `json:"error" example:"Login can not be empty"`
}
