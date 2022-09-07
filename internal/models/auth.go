package models

type ResetPasswordRequest struct {
	Email     string `json:"email" binding:"required,email"`
	ResetCode string `json:"reset_code,omitempty"`
	Password  string `json:"password,omitempty"`
}

type ResetPasswordCache struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	ResetCode string `json:"reset_code"`
	Verified  bool   `json:"verified"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required,min=3,max=30"`
}

type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type AuthenticationResponse struct {
	AccessToken  string                  `json:"access_token"`
	RefreshToken string                  `json:"refresh_token"`
	Permissions  []GetPermissionResponse `json:"permissions"`
}

type HasAccessRequest struct {
	StepID   string `json:"step_id" binding:"required"`
	DriverID string `json:"driver_id" binding:"required"`
}

type HasAccessResponse struct {
	CanView   bool `json:"can_view"`
	CanEdit   bool `json:"can_edit"`
	CanCreate bool `json:"can_create"`
	CanDelete bool `json:"can_delete"`
}
