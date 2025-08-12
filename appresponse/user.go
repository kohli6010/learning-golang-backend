package appresponse

// UserResponse ...
type UserResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"Name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	Avatar    string `json:"avatar"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
