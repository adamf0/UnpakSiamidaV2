package application

type CreateUserCommand struct {
    Name     	 string
	Username     string
	Password     string
	Email        string
	FakultasUnit *string
}
