package application

type UpdateUserCommand struct {
    Uuid     	 string
	Name     	 string
	Username     string
	Password     *string
	Email        string
	Level 		 string
	FakultasUnit *string
}
