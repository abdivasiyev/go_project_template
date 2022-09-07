package models

type CreatePermissionRequest struct {
	Alias string `json:"alias"`
	Name  string `json:"name"`
}

type UpdatePermissionRequest struct {
	ID    string `json:"id"`
	Alias string `json:"alias"`
	Name  string `json:"name"`
}

type GetPermissionResponse struct {
	ID              string `json:"id"`
	GroupID         string `json:"-"`
	Alias           string `json:"alias"`
	Name            string `json:"name"`
	Path            string `json:"path"`
	Method          string `json:"method"`
	QueryParam      string `json:"query_param,omitempty"`
	QueryParamValue string `json:"query_param_value,omitempty"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type GetAllPermissionRequest struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type GetAllPermissionResponse struct {
	Count       int                     `json:"count"`
	Permissions []GetPermissionResponse `json:"permissions"`
}
