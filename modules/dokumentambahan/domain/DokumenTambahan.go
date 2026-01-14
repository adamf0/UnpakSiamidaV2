package domain

import (
	"strings"
	"time"

	common "UnpakSiamida/common/domain"
	domainrenstra "UnpakSiamida/modules/renstra/domain"

	"github.com/google/uuid"
)

type DokumenTambahan struct {
	common.Entity

	ID                      uint      `gorm:"primaryKey;autoIncrement"`
	UUID                    uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	Renstra                 uint      `gorm:"column:id_renstra;"`
	TemplateDokumenTambahan uint      `gorm:"column:id_template_dokumen_tambahan;"`
	Link                    *string   `gorm:"column:file"`
	CapaianAuditor          *string   `gorm:"column:capaian_auditor;"`
	CatatanAuditor          *string   `gorm:"column:catatan_auditor;"`
}

func (DokumenTambahan) TableName() string {
	return "dokumen_tambahan"
}

func UpdateDokumenTambahan(
	prev *DokumenTambahan,
	renstra *domainrenstra.Renstra,
	Uuid uuid.UUID,
	UuidRenstra uuid.UUID,
	Tahun string,
	Mode string,
	Granted string,
	Link *string,
	CapaianAuditor *string,
	CatatanAuditor *string,
) common.ResultValue[*DokumenTambahan] {

	if renstra == nil {
		return common.FailureValue[*DokumenTambahan](InvalidRenstra())
	}

	if prev == nil {
		return common.FailureValue[*DokumenTambahan](EmptyData())
	}

	if prev.UUID != Uuid ||
		renstra.UUID != UuidRenstra ||
		renstra.Tahun != Tahun {
		return common.FailureValue[*DokumenTambahan](InvalidData())
	}

	if !contains([]string{"auditee", "auditor1", "auditor2"}, Mode) {
		return common.FailureValue[*DokumenTambahan](RejectAction())
	}

	if !IsGrantedAccess(Tahun, Mode, Granted) {
		return common.FailureValue[*DokumenTambahan](NotGranted())
	}

	switch Mode {
	case "auditee":
		prev.Link = Link

	case "auditor2":
		prev.CapaianAuditor = CapaianAuditor
		prev.CatatanAuditor = CatatanAuditor
	}

	prev.Raise(DokumenTambahanUpdatedEvent{ //[pr] ketika tersave maka buat kts
		EventID:             uuid.New(),
		OccurredOn:          time.Now().UTC(),
		DokumenTambahanUUID: prev.UUID,
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
