package models

type PageRequest struct {
	Page  int `json:"page" form:"page" default:"1" example:"1"`
	Limit int `json:"limit" form:"limit" default:"5" example:"5"`
}
