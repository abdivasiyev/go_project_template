package models

type CreateUserRequest struct {
	ID        string `json:"id" swaggerignore:"true"`
	CompanyID string `json:"company_id"`
	RoleID    string `json:"role_id" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

type UpdateUserRequest struct {
	ID          string `json:"id" swaggerignore:"true"`
	CompanyID   string `json:"company_id"`
	RoleID      string `json:"role_id" binding:"required"`
	Username    string `json:"username" binding:"required"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type UpdateProfileRequest struct {
	ID          string `json:"id" swaggerignore:"true"`
	Username    string `json:"username" binding:"required"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	ImageID     string `json:"image_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type GetUserResponse struct {
	ID           string             `json:"id,omitempty"`
	Company      GetCompanyResponse `json:"company,omitempty"`
	Role         GetRoleResponse    `json:"role"`
	Username     string             `json:"username,omitempty"`
	FirstName    string             `json:"first_name,omitempty"`
	LastName     string             `json:"last_name,omitempty"`
	ImageID      string             `json:"image_id"`
	PasswordHash string             `json:"-" swaggerignore:"true"`
	CreatedAt    string             `json:"created_at,omitempty"`
	UpdatedAt    string             `json:"updated_at,omitempty"`
}

type GetAllUsersRequest struct {
	PageRequest
	Search string `json:"search" form:"search"`
}

type GetAllUsersResponse struct {
	Count int               `json:"count"`
	Users []GetUserResponse `json:"users"`
}
