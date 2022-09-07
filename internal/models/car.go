package models

type GetCarDashboardResponse struct {
	Total     int    `json:"total"`
	Growth    string `json:"growth"`
	Booked    int    `json:"booked"`
	Available int    `json:"available"`
}

type CreateCarRequest struct {
	ID          string            `json:"id" swaggerignore:"true"`
	Make        string            `json:"make" example:"Toyota"`
	Model       string            `json:"model" example:"Camry"`
	Year        int               `json:"year" example:"2014"`
	Color       string            `json:"color" example:"Black"`
	PlateNumber string            `json:"plate_number" example:"AV82062"`
	ImageID     string            `json:"image_id,omitempty"`
	Status      GetStatusResponse `json:"status"`
}

type GetCarResponse struct {
	ID          string            `json:"id" swaggerignore:"true"`
	Make        string            `json:"make" example:"Toyota"`
	Model       string            `json:"model" example:"Camry"`
	Year        int               `json:"year" example:"2014"`
	Color       string            `json:"color" example:"Black"`
	PlateNumber string            `json:"plate_number" example:"AV82062"`
	ImageID     string            `json:"image_id,omitempty"`
	Status      GetStatusResponse `json:"status"`
}

type GetAllCarsRequest struct {
	PageRequest
	Search string `json:"search" example:"Nissan"`
}

type GetAllCarsResponse struct {
	Count int              `json:"count" example:"10"`
	Cars  []GetCarResponse `json:"cars"`
}

type CarPickupDropOffRequest struct {
	ID           string `json:"id" swaggerignore:"true"`
	CarID        string `json:"car_id" swaggerignore:"true"`
	FullName     string `json:"full_name" binding:"required"`
	PhoneNumber  string `json:"phone_number" binding:"required"`
	CompanyID    string `json:"company_id" binding:"required"`
	KeyGivenBy   string `json:"key_given_by" binding:"required"`
	KeyGivenByID string `json:"key_given_by_id" swaggerignore:"true"`
	Location     string `json:"location" binding:"required"`
	Odometer     string `json:"odometer" binding:"required"`
	PickupDate   string `json:"pickup_date" binding:"required"`
	DropOffDate  string `json:"drop_off_date" binding:"required"`
	Type         string `json:"type" binding:"required,oneof=pickup dropoff"`
}

type GetCarPickupDropOffHistoryRequest struct {
	CarID string `json:"car_id" swaggerignore:"true"`
	PageRequest
}

type GetCarPickupDropOffHistoryResponse struct {
	Count   int                           `json:"count"`
	History []GetCarPickupDropOffResponse `json:"history"`
}

type GetCarPickupDropOffResponse struct {
	ID           string             `json:"id"`
	FullName     string             `json:"full_name"`
	PhoneNumber  string             `json:"phone_number"`
	Company      GetCompanyResponse `json:"company"`
	KeyGivenBy   string             `json:"key_given_by"`
	KeyGivenByID string             `json:"key_given_by_id"`
	Location     string             `json:"location"`
	Odometer     string             `json:"odometer"`
	PickupDate   string             `json:"pickup_date"`
	DropOffDate  string             `json:"drop_off_date"`
	CreatedAt    string             `json:"created_at,omitempty"`
	Type         string             `json:"type"`
}
