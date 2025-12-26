package domain

import (
	"time"

	common "UnpakSiamida/common/domain"

	"github.com/google/uuid"
)

type JenisFile struct {
	common.Entity

	ID              uint       `gorm:"primaryKey;autoIncrement"`
	UUID            uuid.UUID  `gorm:"type:char(36);uniqueIndex"`
	Nama       		string     `gorm:"type:longtext;not null"`
}

func (JenisFile) TableName() string {
	return "jenis_file_renstra"
}

// === CREATE ===
func NewJenisFile(nama string) common.ResultValue[*JenisFile] {

	jenisfile := &JenisFile{
		UUID:         uuid.New(),
		Nama:         nama,
	}

	jenisfile.Raise(JenisFileCreatedEvent{
		EventID:    uuid.New(),
		OccurredOn: time.Now().UTC(),
		JenisFileUUID:   jenisfile.UUID,
	})

	return common.SuccessValue(jenisfile)
}

// === UPDATE ===
func UpdateJenisFile(
	prev *JenisFile, 
	uid uuid.UUID,
	nama string,
) common.ResultValue[*JenisFile] {

	if prev == nil {
		return common.FailureValue[*JenisFile](EmptyData())
	}

	if prev.UUID != uid {
		return common.FailureValue[*JenisFile](InvalidData())
	}

	prev.Nama = nama

	prev.Raise(JenisFileUpdatedEvent{
		EventID:   	uuid.New(),
		OccurredOn: time.Now().UTC(),
		JenisFileUUID:   prev.UUID,
	})

	return common.SuccessValue(prev)
}