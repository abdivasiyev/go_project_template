package models

type GetAppVersionResponse struct {
	Version     string `json:"version"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ForceUpdate bool   `json:"force_update"`
}
