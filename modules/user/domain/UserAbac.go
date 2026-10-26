package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserAbac struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement" json:"ID"`
	UUID               uuid.UUID `gorm:"type:char(36);uniqueIndex" json:"UUID"`
	IdUser             uint      `gorm:"column:id_user;not null" json:"IdUser"`
	Tahun              string    `gorm:"column:tahun;size:10;not null" json:"Tahun"`
	IdFakultasUnit     *uint     `gorm:"column:id_fakultas_unit" json:"IdFakultasUnit"`
	Role               string    `gorm:"column:role;size:50" json:"Role"`
	Status             string    `gorm:"column:status;size:20;default:'active'" json:"Status"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"CreatedAt"`
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"UpdatedAt"`

	// Joined helper fields for presentation
	UserName           string `gorm:"column:user_name;->" json:"UserName,omitempty"`
	UserUsername       string `gorm:"column:user_username;->" json:"UserUsername,omitempty"`
	UserEmail          string `gorm:"column:user_email;->" json:"UserEmail,omitempty"`
	UserUuid           string `gorm:"column:user_uuid;->" json:"UserUuid,omitempty"`
	UserLevel          string `gorm:"column:user_level;->" json:"UserLevel,omitempty"`
	NamaFakultasUnit   string `gorm:"column:nama_fakultas_unit;->" json:"NamaFakultasUnit,omitempty"`
	JenjangFakultasUnit string `gorm:"column:jenjang_fakultas_unit;->" json:"JenjangFakultasUnit,omitempty"`
	UuidFakultasUnit   string `gorm:"column:uuid_fakultas_unit;->" json:"UuidFakultasUnit,omitempty"`
	TypeFakultasUnit   string `gorm:"column:type_fakultas_unit;->" json:"TypeFakultasUnit,omitempty"`
	FakultasUnitParent string `gorm:"column:fakultas_unit_parent;->" json:"FakultasUnitParent,omitempty"`
	ComputedRole       string `gorm:"column:computed_role;->" json:"ComputedRole,omitempty"`
}

func (UserAbac) TableName() string {
	return "user_abac"
}
