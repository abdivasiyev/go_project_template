package models

type GradientResponse struct {
	StartColor string `json:"start_color" binding:"required" example:"#101010"`
	EndColor   string `json:"end_color" binding:"required" example:"#202020"`
}
