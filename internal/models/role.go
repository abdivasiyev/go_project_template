package models

type CreateRoleRequest struct {
	ID          string   `json:"id" swaggerignore:"true"`
	Alias       string   `json:"alias" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions" binding:"required"`
}

type GetModulesResponse struct {
	Count   int                 `json:"count"`
	Modules []GetModuleResponse `json:"modules"`
}

type GetModuleResponse struct {
	ID     string                       `json:"id"`
	Alias  string                       `json:"alias"`
	Name   string                       `json:"name"`
	Groups []GetPermissionGroupResponse `json:"groups"`
}

type GetPermissionGroupResponse struct {
	ID          string                  `json:"id"`
	ModuleID    string                  `json:"-"`
	Name        string                  `json:"name"`
	Alias       string                  `json:"alias"`
	Permissions []GetPermissionResponse `json:"permissions"`
}

type GetRoleResponse struct {
	ID          string                  `json:"id"`
	Alias       string                  `json:"alias"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Permissions []GetPermissionResponse `json:"permissions,omitempty"`
	CreatedAt   string                  `json:"created_at"`
	UpdatedAt   string                  `json:"updated_at"`
}
type GetAllRoleRequest struct {
	PageRequest
}

type GetAllRoleResponse struct {
	Count int               `json:"count"`
	Roles []GetRoleResponse `json:"roles"`
}
