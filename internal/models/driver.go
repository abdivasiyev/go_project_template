package models

type GetDriverDashboardResponse struct {
	Total     int    `json:"total"`
	Growth    string `json:"growth"`
	Pending   int    `json:"pending"`
	InProcess int    `json:"in_process"`
	Active    int    `json:"active"`
}

type GetAllDriversRequest struct {
	PageRequest
	Search string `json:"search" form:"search"`
}

type GetAllDriversResponse struct {
	Count   int                 `json:"count"`
	Drivers []GetDriverResponse `json:"drivers"`
}

type CreateDriverRequest struct {
	ID        string `json:"id" swaggerignore:"true"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required"`
	ImageID   string `json:"image_id" format:"uuid4,omitempty"`
}

type GetDriverResponse struct {
	ID                   string                        `json:"id"`
	FirstName            string                        `json:"first_name"`
	LastName             string                        `json:"last_name"`
	Email                string                        `json:"email"`
	Phone                string                        `json:"phone"`
	Image                GetFileResponse               `json:"image"`
	OnBoardDate          string                        `json:"on_board_date"`
	CreatedAt            string                        `json:"created_at"`
	UpdatedAt            string                        `json:"updated_at"`
	DepartmentStatistics []DepartmentStatisticResponse `json:"department_statistics"`
	AssignedTruck        *AssignedTruckResponse        `json:"assigned_truck,omitempty"`
	AssignedTrailer      *AssignedTrailerResponse      `json:"assigned_trailer,omitempty"`
	Status               GetStatusResponse             `json:"status"`
}
