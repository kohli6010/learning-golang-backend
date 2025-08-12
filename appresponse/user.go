package appresponse

// UserResponse ...
type UserResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"Name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
