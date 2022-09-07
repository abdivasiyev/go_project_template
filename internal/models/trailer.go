package models

type (
	TrailerInspectionType int
	TrailerDeviceType     int
)

type GetTrailerDashboardResponse struct {
	Total    int    `json:"total"`
	Growth   string `json:"growth"`
	Vacant   int    `json:"vacant"`
	Assigned int    `json:"assigned"`
	Full     int    `json:"full"`
}

type GetTrailerPickUpDropOffHistoryRequest struct {
	PageRequest
	Search string `json:"search" form:"search"`
}

type GetTrailerPickUpDropOffHistoryResponse struct {
	Count   int                           `json:"count"`
	History []TrailerPickUpDropOffHistory `json:"history"`
}

type TrailerPickUpDropOffHistory struct {
	ID             string `json:"id"`
	TrailerID      string `json:"trailer_id"`
	TrailerNumber  string `json:"trailer_number"`
	DriverID       string `json:"driver_id"`
	DriverName     string `json:"driver_name"`
	InspectionType string `json:"inspection_type"`
	CreatedAt      string `json:"created_at"`
}

type GetTrailerInspectionRequest struct {
	TrailerID string `json:"trailer_id" uri:"trailer_id"`
}

type GetTrailerInspectionResponse struct {
	ID                               string               `json:"id,omitempty"`
	Trailer                          GetTrailerResponse   `json:"trailer"`
	IsEmpty                          bool                 `json:"is_empty"`
	Location                         string               `json:"location"`
	FuelLevelImages                  []string             `json:"fuel_level_images"`
	LeftSideImages                   []string             `json:"left_side_images"`
	FrontSideImages                  []string             `json:"front_side_images"`
	RightSideImages                  []string             `json:"right_side_images"`
	BackSideImages                   []string             `json:"back_side_images"`
	InSideImages                     []string             `json:"in_side_images"`
	TireImages                       []string             `json:"tire_images"`
	DamageImages                     []string             `json:"damage_images"`
	NeedRepairDevices                []string             `json:"need_repair_devices"`
	DriverSignatureID                string               `json:"driver_signature_id"`
	CompanyRepresentativeSignatureID string               `json:"company_representative_signature_id"`
	Comments                         []GetCommentResponse `json:"comments"`
	InspectionType                   string               `json:"inspection_type"`
}

type CreateTrailerInspectionRequest struct {
	ID                               string               `json:"id" swaggerignore:"true"`
	TrailerID                        string               `json:"trailer_id" binding:"required,uuid4"`
	DriverID                         string               `json:"driver_id" binding:"required,uuid4"`
	InspectionType                   string               `json:"inspection_type" binding:"required,oneof=pickup drop-off"`
	IsEmpty                          bool                 `json:"is_empty"`
	Location                         string               `json:"location" binding:"required"`
	FuelLevelImages                  []string             `json:"fuel_level_images" binding:"required"`
	LeftSideImages                   []string             `json:"left_side_images" binding:"required"`
	FrontSideImages                  []string             `json:"front_side_images" binding:"required"`
	RightSideImages                  []string             `json:"right_side_images" binding:"required"`
	BackSideImages                   []string             `json:"back_side_images" binding:"required"`
	InSideImages                     []string             `json:"in_side_images" binding:"required"`
	TireImages                       []string             `json:"tire_images" binding:"required"`
	DamageImages                     []string             `json:"damage_images" binding:"required"`
	NeedRepairDevices                []string             `json:"need_repair_devices" binding:"required"`
	DriverSignatureID                string               `json:"driver_signature_id" binding:"required,uuid4"`
	CompanyRepresentativeSignatureID string               `json:"company_representative_signature_id" binding:"required,uuid4"`
	Comments                         []GetCommentResponse `json:"comments"`
}

type AssignedTrailerResponse struct {
	ID     string `json:"id"`
	Number string `json:"number"`
}

type CreateTrailerRequest struct {
	ID          string `json:"id" swaggerignore:"true"`
	Make        string `json:"make" binding:"required" example:"Volvo"`
	Number      string `json:"number" binding:"required" example:"#120"`
	YearMade    int    `json:"year_made" binding:"required" example:"2020"`
	TrailerType string `json:"trailer_type" binding:"required" example:"Dry"`
	PlateNumber string `json:"plate_number" binding:"required" example:"3920-3920"`
}

type GetTrailerResponse struct {
	ID             string            `json:"id" example:"973cb235-bdc7-4ffc-94f8-bf4eaf23b778"`
	Make           string            `json:"make" example:"Volvo"`
	Number         string            `json:"number" example:"#120"`
	YearMade       int               `json:"year_made" example:"2020"`
	TrailerType    string            `json:"trailer_type" example:"Dry"`
	PlateNumber    string            `json:"plate_number"`
	Status         GetStatusResponse `json:"status"`
	CreatedAt      string            `json:"created_at"`
	UpdatedAt      string            `json:"updated_at"`
	AssignedDriver struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Since string `json:"since"`
	} `json:"assigned_driver"`
}

type GetAllTrailersRequest struct {
	PageRequest
	Search   string `json:"search" form:"search" example:"Nissan"`
	StatusID string `json:"status_id" form:"status_id"`
}

type GetAllTrailersResponse struct {
	Count    int                  `json:"count" example:"10"`
	Trailers []GetTrailerResponse `json:"trailers"`
}
