package domain

type Account struct {
	ID           uint         `gorm:"column:id;primaryKey" json:"-"`
	UUID         string       `json:"UUID"`
	NidnUsername string       `json:"Username"`
	Password     string       `json:"-"`
	Level        string       `json:"Level"`
	Name         string       `json:"Name"`
	Email            string       `json:"Email"`
	FakultasUnit     string       `json:"FakultasUnit"`
	NamaFakultasUnit string       `gorm:"column:nama_fakultas_unit" json:"NamaFakultasUnit,omitempty"`
	ExtraRole        []ExtraRole  `gorm:"-" json:"ExtraRole,omitempty"`
}

func (Account) TableName() string {
	return "users"
}