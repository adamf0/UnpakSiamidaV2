package domain

import (
	"time"

	common "UnpakSiamida/common/domain"

	"github.com/google/uuid"
)

type StandarRenstra struct {
	common.Entity
	ID           uint       `gorm:"primaryKey;autoIncrement"`
	UUID         uuid.UUID  `gorm:"type:char(36);uniqueIndex"`
	Nama         string     `gorm:"size:255;not null"`
}
func (StandarRenstra) TableName() string {
	return "master_standar_renstra"
}


// === CREATE ===
func NewStandarRenstra(nama string) common.ResultValue[*StandarRenstra] {

	standarrenstra := &StandarRenstra{
		UUID:         uuid.New(),
		Nama:         nama,
	}

	standarrenstra.Raise(StandarRenstraCreatedEvent{
		EventID:    uuid.New(),
		OccurredOn: time.Now().UTC(),
		StandarRenstraUUID:   standarrenstra.UUID,
	})

	return common.SuccessValue(standarrenstra)
}

// === UPDATE ===
func UpdateStandarRenstra(
	prev *StandarRenstra, 
	uid uuid.UUID,
	nama string,
) common.ResultValue[*StandarRenstra] {

	if prev == nil {
		return common.FailureValue[*StandarRenstra](EmptyData())
	}

	if prev.UUID != uid {
		return common.FailureValue[*StandarRenstra](InvalidData())
	}

	prev.Nama = nama

	prev.Raise(StandarRenstraUpdatedEvent{
		EventID:   	uuid.New(),
		OccurredOn: time.Now().UTC(),
		StandarRenstraUUID:   prev.UUID,
	})

	return common.SuccessValue(prev)
}
