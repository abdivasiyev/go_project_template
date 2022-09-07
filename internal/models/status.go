package models

type GetStatusResponse struct {
	ID       string `json:"id" example:"973cb235-bdc7-4ffc-94f8-bf4eaf23b778"`
	Alias    string `json:"alias" example:"booked"`
	Name     string `json:"name" example:"Booked"`
	Sequence int    `json:"sequence" example:"1"`
	Color    string `json:"color" example:"#EB5757"`
}
