package models

type (
	InspectionType int
	DeviceType     int
)

type GetTruckDashboardResponse struct {
	Total    int    `json:"total"`
	Growth   string `json:"growth"`
	Vacant   int    `json:"vacant"`
	Assigned int    `json:"assigned"`
	Full     int    `json:"full"`
}

type GetTruckPickUpDropOffHistoryRequest struct {
	PageRequest
	Search string `json:"search" form:"search"`
}

type GetTruckPickUpDropOffHistoryResponse struct {
	Count   int                         `json:"count"`
	History []TruckPickUpDropOffHistory `json:"history"`
}

type TruckPickUpDropOffHistory struct {
	ID             string `json:"id"`
	TruckID        string `json:"truck_id"`
	TruckNumber    string `json:"truck_number"`
	DriverID       string `json:"driver_id"`
	DriverName     string `json:"driver_name"`
	InspectionType string `json:"inspection_type"`
	CreatedAt      string `json:"created_at"`
}

type GetTruckInspectionRequest struct {
	TruckID string `json:"truck_id" uri:"truck_id"`
}

type GetTruckInspectionResponse struct {
	ID                               string               `json:"id,omitempty"`
	Truck                            GetTruckResponse     `json:"truck"`
	Location                         string               `json:"location"`
	OdometerImages                   []string             `json:"odometer_images"`
	FuelLevelImages                  []string             `json:"fuel_level_images"`
	DriverSideImages                 []string             `json:"driver_side_images"`
	FrontSideImages                  []string             `json:"front_side_images"`
	PassengerSideImages              []string             `json:"passenger_side_images"`
	BackSideImages                   []string             `json:"back_side_images"`
	TireImages                       []string             `json:"tire_images"`
	DamageImages                     []string             `json:"damage_images"`
	IncabDevices                     []string             `json:"incab_devices"`
	ExternalDisplayed                []string             `json:"external_displayed"`
	DriverSignatureID                string               `json:"driver_signature_id"`
	CompanyRepresentativeSignatureID string               `json:"company_representative_signature_id"`
	Comments                         []GetCommentResponse `json:"comments"`
	InspectionType                   string               `json:"inspection_type"`
}

type CreateTruckInspectionRequest struct {
	ID                               string               `json:"id" swaggerignore:"true"`
	TruckID                          string               `json:"truck_id" binding:"required,uuid4"`
	DriverID                         string               `json:"driver_id" binding:"required,uuid4"`
	InspectionType                   string               `json:"inspection_type" binding:"required,oneof=pickup drop-off"`
	Location                         string               `json:"location" binding:"required"`
	OdometerImages                   []string             `json:"odometer_images" binding:"required"`
	FuelLevelImages                  []string             `json:"fuel_level_images" binding:"required"`
	DriverSideImages                 []string             `json:"driver_side_images" binding:"required"`
	FrontSideImages                  []string             `json:"front_side_images" binding:"required"`
	PassengerSideImages              []string             `json:"passenger_side_images" binding:"required"`
	BackSideImages                   []string             `json:"back_side_images" binding:"required"`
	TireImages                       []string             `json:"tire_images" binding:"required"`
	DamageImages                     []string             `json:"damage_images" binding:"required"`
	IncabDevices                     []string             `json:"incab_devices" binding:"required"`
	ExternalDisplayed                []string             `json:"external_displayed" binding:"required"`
	DriverSignatureID                string               `json:"driver_signature_id" binding:"required,uuid4"`
	CompanyRepresentativeSignatureID string               `json:"company_representative_signature_id" binding:"required,uuid4"`
	Comments                         []GetCommentResponse `json:"comments"`
}

type GetCommentResponse struct {
	CreatedBy string `json:"created_by"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}

type AssignedTruckResponse struct {
	ID     string `json:"id"`
	Number string `json:"number"`
}

type CreateTruckRequest struct {
	ID          string  `json:"id" swaggerignore:"true"`
	Make        string  `json:"make" binding:"required" example:"Volvo"`
	Model       string  `json:"model" binding:"required"`
	Number      string  `json:"number" binding:"required" example:"#120"`
	YearMade    int     `json:"year_made" binding:"required" example:"2020"`
	Milage      float64 `json:"milage" binding:"required" example:"350000"`
	PlateNumber string  `json:"plate_number" binding:"required" example:"3920-3920"`
}

type GetTruckResponse struct {
	ID             string            `json:"id" example:"973cb235-bdc7-4ffc-94f8-bf4eaf23b778"`
	Make           string            `json:"make" example:"Volvo"`
	Number         string            `json:"number" example:"#120"`
	YearMade       int               `json:"year_made" example:"2020"`
	Milage         float64           `json:"milage" example:"350000"`
	PlateNumber    string            `json:"plate_number"`
	Model          string            `json:"model"`
	Status         GetStatusResponse `json:"status"`
	CreatedAt      string            `json:"created_at"`
	UpdatedAt      string            `json:"updated_at"`
	AssignedDriver struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Since string `json:"since"`
	} `json:"assigned_driver"`
}

type GetAllTrucksRequest struct {
	PageRequest
	Search   string `json:"search" form:"search" example:"Nissan"`
	StatusID string `json:"status_id" form:"status_id"`
}

type GetAllTrucksResponse struct {
	Count  int                `json:"count" example:"10"`
	Trucks []GetTruckResponse `json:"trucks"`
}
