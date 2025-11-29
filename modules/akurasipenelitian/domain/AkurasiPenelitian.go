package domain

import (
	"time"
	"github.com/google/uuid"
	"UnpakSiamida/common/domain"
)

type AkurasiPenelitian struct {
	domain.Entity
	ID                   *int
	UUID                 uuid.UUID
	Nama                 string
	Skor                 int
}

func NewAkurasiPenelitian(nama string, skor int) domain.ResultValue[*AkurasiPenelitian] {
	if skor < 0 {
		return domain.FailureValue[*AkurasiPenelitian](InvalidSkor())
	}
	ap := &AkurasiPenelitian{
		UUID: uuid.New(),
		Nama: nama,
		Skor: skor,
	}
	ap.Raise(AkurasiPenelitianCreatedEvent{
		EventID:    uuid.New(),
		OccurredOn: time.Now().UTC(),
	})
	return domain.SuccessValue(ap)
}
