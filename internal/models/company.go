package models

type CreateCompanyRequest struct {
	ID   string `json:"id" swaggerignore:"true"`
	Name string `json:"name"`
}

type UpdateCompanyRequest struct {
	ID   string `json:"id" swaggerignore:"true"`
	Name string `json:"name"`
}

type GetCompanyResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type GetAllCompaniesRequest struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type GetAllCompaniesResponse struct {
	Count     int                  `json:"count"`
	Companies []GetCompanyResponse `json:"companies"`
}
