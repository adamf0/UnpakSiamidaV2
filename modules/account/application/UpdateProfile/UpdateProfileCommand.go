package application

type UpdateProfileCommand struct {
	SID             string
	Username        string
	Name            string
	Email           string
	CurrentPassword string
	NewPassword     string
	ConfirmPassword string
}
