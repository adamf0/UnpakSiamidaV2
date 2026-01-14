package domain

type Account struct {
	ID           string       `json:"-"`
	UUID         string       `json:"UUID"`
	NidnUsername string       `json:"Username"`
	Password     string       `json:"-"`
	Level        string       `json:"Level"`
	Name         string       `json:"Name"`
	Email        string       `json:"Email"`
	FakultasUnit string       `json:"FakultasUnit"`
	ExtraRole    []ExtraRole  `gorm:"-"; json:"ExtraRole,omitempty"`
}

func (Account) TableName() string {
	return "users"
}