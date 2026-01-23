package dto

type TokenResponse struct {
	AccessToken  string `json:"access_token"  example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

func (dto *TokenResponse) Fill(accessToken, refreshToken string) {
	dto.AccessToken = accessToken
	dto.RefreshToken = refreshToken
}
