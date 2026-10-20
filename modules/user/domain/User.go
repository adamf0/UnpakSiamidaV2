package domain

import (
	"time"

	common "UnpakSiamida/common/domain"
	helper "UnpakSiamida/common/helper"
	event "UnpakSiamida/modules/user/event"

	"github.com/google/uuid"
)

type User struct {
	common.Entity    `json:"-"`
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"ID"`
	UUID             uuid.UUID `gorm:"type:char(36);uniqueIndex" json:"UUID"`
	Username         string    `gorm:"column:nidn_username;size:100;not null" json:"Username"`
	Password         string    `gorm:"type:longtext;not null" json:"-"`
	Name             string    `gorm:"size:255;not null" json:"Name"`
	Email            string    `gorm:"size:255;" json:"Email"`
	Level            string    `gorm:"size:255;" json:"Level"`
	FakultasUnit     *int      `gorm:"column:fakultas_unit;" json:"FakultasUnit"`
	NamaFakultasUnit *string   `gorm:"column:nama_fak_prod_unit;->" json:"NamaFakultasUnit,omitempty"`
	UuidFakultasUnit *string   `gorm:"column:uuid_fakultas_unit;->" json:"UuidFakultasUnit,omitempty"`
	IdUserMerge      *uint             `gorm:"column:id_user_merge;" json:"IdUserMerge"`
	NamaUserMerge    *string           `gorm:"column:nama_user_merge;->" json:"NamaUserMerge,omitempty"`
	UserAbac         []UserAbacSummary `gorm:"-" json:"UserAbac"`
}

type UserAbacSummary struct {
	ID               uint   `json:"ID"`
	UUID             string `json:"UUID"`
	IdUser           uint   `json:"IdUser"`
	Tahun            string `json:"Tahun"`
	IdFakultasUnit   uint   `json:"IdFakultasUnit"`
	NamaFakultasUnit string `json:"NamaFakultasUnit"`
	Role             string `json:"Role"`
	ComputedRole     string `json:"ComputedRole"`
	Status           string `json:"Status"`
}

func (User) TableName() string {
	return "users"
}

// === CREATE ===
func NewUser(username string, password string, name string, email string, level string, fakultasunit *int, target *string, tipe *string) common.ResultValue[*User] {

	if !helper.IsValidUnpakEmail(email) {
		return common.FailureValue[*User](InvalidEmail())
	}

	if fakultasunit != nil && *fakultasunit <= 0 {
		return common.FailureValue[*User](InvalidFakultasUnit())
	}

	user := &User{
		UUID:         uuid.New(),
		Username:     username,
		Password:     password,
		Name:         name,
		Email:        email,
		Level:        level,
		FakultasUnit: fakultasunit,
	}

	user.Raise(event.UserCreatedEvent{
		EventID:      uuid.New(),
		OccurredOn:   time.Now().UTC(),
		UserUUID:     user.UUID,
		Username:     username,
		Password:     password,
		Name:         name,
		Email:        email,
		Level:        level,
		FakultasUnit: target,
		Tipe:         tipe,
	})

	return common.SuccessValue(user)
}

// === UPDATE ===
func UpdateUser(
	prev *User,
	uid uuid.UUID,
	username string,
	password *string,
	name string,
	email string,
	level string,
	fakultasunit *int,
	target *string,
	tipe *string,
) common.ResultValue[*User] {

	if prev == nil {
		return common.FailureValue[*User](EmptyData())
	}

	if prev.UUID != uid {
		return common.FailureValue[*User](InvalidData())
	}

	if fakultasunit != nil && *fakultasunit <= 0 {
		return common.FailureValue[*User](InvalidFakultasUnit())
	}

	if !helper.IsValidUnpakEmail(email) {
		return common.FailureValue[*User](InvalidEmail())
	}

	// update opsional
	updatedPassword := prev.Password
	if password != nil {
		updatedPassword = *password
		prev.Password = *password
	}

	prev.Username = username
	prev.Name = name
	prev.Email = email
	prev.Level = level
	prev.FakultasUnit = fakultasunit

	prev.Raise(event.UserUpdatedEvent{
		EventID:      uuid.New(),
		OccurredOn:   time.Now().UTC(),
		UserUUID:     prev.UUID,
		Username:     username,
		Password:     updatedPassword,
		Name:         name,
		Email:        email,
		Level:        level,
		FakultasUnit: target,
		Tipe:         tipe,
	})

	return common.SuccessValue(prev)
}
