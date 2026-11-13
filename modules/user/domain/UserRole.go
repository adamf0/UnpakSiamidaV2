package domain

import (
	common "UnpakSiamida/common/domain"
	"time"

	"github.com/google/uuid"
)

type UserRole struct {
	common.Entity       `json:"-"`
	ID                  uint       `gorm:"primaryKey;autoIncrement" json:"ID"`
	UUID                uuid.UUID  `gorm:"type:char(36);uniqueIndex" json:"UUID"`
	IdUser              uint       `gorm:"column:id_user;not null" json:"IdUser"`
	UserUUID            string     `gorm:"column:user_uuid;->" json:"UserUUID,omitempty"`
	UserName            string     `gorm:"column:user_name;->" json:"UserName,omitempty"`
	UserUsername        string     `gorm:"column:user_username;->" json:"UserUsername,omitempty"`
	UserEmail           string     `gorm:"column:user_email;->" json:"UserEmail,omitempty"`
	Tahun               string     `gorm:"column:tahun;size:10;not null" json:"Tahun"`
	IdFakultasUnit      *uint      `gorm:"column:id_fakultas_unit" json:"IdFakultasUnit"`
	UuidFakultasUnit    *string    `gorm:"column:uuid_fakultas_unit;->" json:"UuidFakultasUnit,omitempty"`
	NamaFakultasUnit    string     `gorm:"column:nama_fak_prod_unit;->" json:"NamaFakultasUnit,omitempty"`
	TypeFakultasUnit    *string    `gorm:"column:type_fakultas_unit;->" json:"TypeFakultasUnit,omitempty"`
	JenjangFakultasUnit *string    `gorm:"column:jenjang_fakultas_unit;->" json:"JenjangFakultasUnit,omitempty"`
	Level               string     `gorm:"column:level;size:50;not null" json:"Level"`
	Status              string     `gorm:"column:status;size:20;not null;default:'active'" json:"Status"`
	CreatedAt           *time.Time `gorm:"column:created_at;->" json:"CreatedAt,omitempty"`
	UpdatedAt           *time.Time `gorm:"column:updated_at;->" json:"UpdatedAt,omitempty"`
}

func (UserRole) TableName() string {
	return "userv3_role"
}
