package domain

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	UserID       string
}