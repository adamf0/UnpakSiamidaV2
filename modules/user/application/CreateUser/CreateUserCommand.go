package application

type CreateUserCommand struct {
	Name             string
	Username         string
	Password         string
	Email            string
	Level            string
	UuidFakultasUnit *string
}
