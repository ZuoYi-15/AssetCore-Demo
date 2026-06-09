package auth

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username        string   `json:"username" binding:"required"`
	Password        string   `json:"password" binding:"required"`
	DisplayName     string   `json:"display_name"`
	Email           string   `json:"email"`
	RoleCode        string   `json:"role_code" binding:"required"`
	PermissionCodes []string `json:"permission_codes"`
}

type UpdateUserRequest struct {
	Username        string   `json:"username" binding:"required"`
	Password        string   `json:"password"`
	DisplayName     string   `json:"display_name"`
	Email           string   `json:"email"`
	Status          string   `json:"status" binding:"required"`
	RoleCode        string   `json:"role_code" binding:"required"`
	PermissionCodes []string `json:"permission_codes"`
}

type UserQuery struct {
	Keyword string
	Status  string
	Role    string
}

type UserProfile struct {
	ID          uint64   `json:"id"`
	Username    string   `json:"username"`
	DisplayName string   `json:"display_name"`
	Email       string   `json:"email"`
	Status      string   `json:"status"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  UserProfile `json:"user"`
}

type Claims struct {
	UserID      uint64   `json:"user_id"`
	Username    string   `json:"username"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	ExpiresAt   int64    `json:"exp"`
	Issuer      string   `json:"iss"`
}
