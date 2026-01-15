package domain

import (
	"strings"
	"time"

	common "UnpakSiamida/common/domain"
	domainrenstra "UnpakSiamida/modules/renstra/domain"
	event "UnpakSiamida/modules/renstranilai/event"

	"github.com/google/uuid"
)

type RenstraNilai struct {
	common.Entity

	ID              uint      `gorm:"primaryKey;autoIncrement"`
	UUID            uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	Renstra         uint      `gorm:"column:id_renstra;"`
	TemplateRenstra uint      `gorm:""`
	Tugas           string    `gorm:""`
	Capaian         *string   `gorm:""`
	Catatan         *string   `gorm:""`
	LinkBukti       *string   `gorm:"column:link_bukti;"`
	CapaianAuditor  *string   `gorm:"column:capaian_auditor;"`
	CatatanAuditor  *string   `gorm:"column:catatan_auditor;"`
}

func (RenstraNilai) TableName() string {
	return "renstra_nilai"
}

func UpdateRenstraNilai(
	prev *RenstraNilai,
	renstra *domainrenstra.Renstra,
	Uuid uuid.UUID,
	UuidRenstra uuid.UUID,
	Tahun string,
	Mode string,
	Granted string,
	Capaian *string,
	Catatan *string,
	LinkBukti *string,
	CapaianAuditor *string,
	CatatanAuditor *string,
) common.ResultValue[*RenstraNilai] {

	if renstra == nil {
		return common.FailureValue[*RenstraNilai](InvalidRenstra())
	}

	if prev == nil {
		return common.FailureValue[*RenstraNilai](EmptyData())
	}

	if prev.UUID != Uuid ||
		renstra.UUID != UuidRenstra ||
		renstra.Tahun != Tahun {
		return common.FailureValue[*RenstraNilai](InvalidData())
	}

	if !contains([]string{"auditee", "auditor1", "auditor2"}, Mode) {
		return common.FailureValue[*RenstraNilai](RejectAction())
	}

	if !IsGrantedAccess(Tahun, Mode, Granted) {
		return common.FailureValue[*RenstraNilai](NotGranted())
	}

	switch Mode {
	case "auditee":
		prev.Capaian = Capaian
		prev.Catatan = Catatan
		prev.LinkBukti = LinkBukti

	case "auditor1":
		prev.CapaianAuditor = CapaianAuditor
		prev.CatatanAuditor = CatatanAuditor
	}

	prev.Raise(event.RenstraNilaiUpdatedEvent{ //[pr] ketika tersave maka buat kts
		EventID:          uuid.New(),
		OccurredOn:       time.Now().UTC(),
		RenstraNilaiUUID: prev.UUID,
	})

	return common.SuccessValue(prev)
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func IsGrantedAccess(tahun, mode, granted string) bool {
	requiredKey := tahun + "#" + mode

	grantedList := strings.Split(granted, ",")

	for _, g := range grantedList {
		if strings.TrimSpace(g) == requiredKey {
			return true
		}
	}

	return false
}
