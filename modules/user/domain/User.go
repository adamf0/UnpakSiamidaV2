package domain

import (
	"time"

	common "UnpakSiamida/common/domain"
	helper "UnpakSiamida/common/helper"
	event "UnpakSiamida/modules/user/event"

	"github.com/google/uuid"
)

type User struct {
	common.Entity
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	UUID         uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	Username     string    `gorm:"column:nidn_username;size:100;not null"`
	Password     string    `gorm:"type:longtext;not null"`
	Name         string    `gorm:"size:255;not null"`
	Email        string    `gorm:"size:255;"`
	Level        string    `gorm:"size:255;"`
	FakultasUnit *int      `gorm:"column:fakultas_unit;"`
}

func (User) TableName() string {
	return "users"
}

// === CREATE ===
func NewUser(username string, password string, name string, email string, level string, fakultasunit *int, target *string, tipe *string) common.ResultValue[*User] {

	if !helper.IsValidUnpakEmail(email) {
		return common.FailureValue[*User](InvalidEmail())
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

	// check email
	if !helper.IsValidUnpakEmail(email) {
		return common.FailureValue[*User](InvalidEmail())
	}

	// update opsional
	updatedPassword := prev.Password
	if password != nil {
		updatedPassword = *password
		prev.Password = *password
	}

	prev.Name = name
	prev.Email = email
	prev.Level = level

	if fakultasunit != nil {
		prev.FakultasUnit = fakultasunit
	}

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
