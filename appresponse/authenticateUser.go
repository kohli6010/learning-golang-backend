package appresponse

// AuthenticateUserResponse ...
type AuthenticateUserResponse struct {
	Success        bool   `json:"success"`
	AccessToken    string `json:"access"`
	AccessExpires  string `json:"access_expires"`
	RefreshTokens  string `json:"refresh_token"`
	RefreshExpires string `json:"refresh_expires"`
}
