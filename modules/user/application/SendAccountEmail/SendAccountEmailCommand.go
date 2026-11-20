package application

type SendAccountEmailCommand struct {
	UserUUIDs []string `json:"user_uuids"`
}
