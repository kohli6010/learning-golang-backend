package apprequest

// UpsertUserPersonalDetailRequest ...
type UpsertUserPersonalDetailRequest struct {
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Avatar   string `json:"avatar"`
}