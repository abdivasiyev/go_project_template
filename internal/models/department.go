package models

type CreateDepartmentRequest struct {
	ID       string           `json:"id" swaggerignore:"true"`
	Alias    string           `json:"alias" binding:"required"`
	Name     string           `json:"name" binding:"required"`
	Sequence int              `json:"sequence" binding:"required"`
	Gradient GradientResponse `json:"gradient" binding:"required"`
}

type DepartmentStatisticResponse struct {
	ID         string           `json:"id"`
	Alias      string           `json:"alias"`
	Name       string           `json:"name"`
	Sequence   int              `json:"sequence"`
	Percentage float64          `json:"percentage"`
	Gradient   GradientResponse `json:"gradient"`
}

type GetAllDepartmentsRequest struct {
	PageRequest
	Search string `json:"search" form:"search"`
}

type GetAllDepartmentsResponse struct {
	Count       int                     `json:"count"`
	Departments []GetDepartmentResponse `json:"departments"`
}

type GetDepartmentResponse struct {
	ID       string           `json:"id"`
	Alias    string           `json:"alias"`
	Name     string           `json:"name"`
	Sequence int              `json:"sequence"`
	Gradient GradientResponse `json:"gradient"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
