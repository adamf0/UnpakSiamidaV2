package application

import (
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	domainRisk "UnpakSiamida/modules/riskassessment/domain"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateRiskAssessmentCommandHandler struct {
	Repo             domainRisk.IRiskAssessmentRepository
	RepoFakultasUnit domainfakultasunit.IFakultasUnitRepository
	DB               *gorm.DB
}

func (h *UpdateRiskAssessmentCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateRiskAssessmentCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsedUuid, err := uuid.Parse(cmd.UUID)
	if err != nil {
		return "", domainRisk.NotFound(cmd.UUID)
	}

	prev, err := h.Repo.GetByUuid(ctx, parsedUuid)
	if err != nil || prev == nil {
		return "", domainRisk.NotFound(cmd.UUID)
	}

	var fuID *int = cmd.IdFakultasUnit
	if fuID == nil {
		fuID = prev.IdFakultasUnit
	}
	if strings.TrimSpace(cmd.FakultasUuid) != "" {
		if parsed, err := uuid.Parse(strings.TrimSpace(cmd.FakultasUuid)); err == nil {
			target, err := h.RepoFakultasUnit.GetDefaultByUuid(ctx, parsed)
			if err == nil && target != nil {
				idInt := int(target.ID)
				fuID = &idInt
			}
		} else if id, err := strconv.Atoi(strings.TrimSpace(cmd.FakultasUuid)); err == nil && id > 0 {
			fuID = &id
		}
	}

	var tahunID *int = cmd.IdMasterTahun
	if tahunID == nil {
		tahunID = prev.IdMasterTahun
	}
	var tahunStr string = strings.TrimSpace(cmd.Tahun)
	if tahunStr == "" {
		tahunStr = prev.Tahun
	}

	if strings.TrimSpace(cmd.TahunUuid) != "" {
		if parsed, err := uuid.Parse(strings.TrimSpace(cmd.TahunUuid)); err == nil {
			var mt struct {
				ID    int
				Tahun string
			}
			if err := h.DB.WithContext(ctx).Table("master_tahun").Select("id, tahun").Where("uuid = ?", parsed).Take(&mt).Error; err == nil {
				tahunID = &mt.ID
				if strings.TrimSpace(cmd.Tahun) == "" {
					tahunStr = mt.Tahun
				}
			}
		} else if id, err := strconv.Atoi(strings.TrimSpace(cmd.TahunUuid)); err == nil && id > 0 {
			tahunID = &id
		}
	}

	idMataProgram := cmd.IdMataProgram
	if idMataProgram <= 0 && strings.TrimSpace(cmd.MataProgramUuid) != "" {
		if parsed, err := uuid.Parse(strings.TrimSpace(cmd.MataProgramUuid)); err == nil {
			var mp struct {
				ID int
			}
			if err := h.DB.WithContext(ctx).Table("mata_program").Select("id").Where("uuid = ?", parsed).Take(&mp).Error; err == nil {
				idMataProgram = mp.ID
			}
		} else if id, err := strconv.Atoi(strings.TrimSpace(cmd.MataProgramUuid)); err == nil && id > 0 {
			idMataProgram = id
		}
	}
	if idMataProgram <= 0 {
		idMataProgram = prev.IdMataProgram
	}

	isu := strings.TrimSpace(cmd.IsuInternalEksternal)
	if isu == "" {
		isu = strings.TrimSpace(cmd.Penyebab)
	}
	if isu == "" {
		isu = prev.IsuInternalEksternal
	}

	deskripsi := strings.TrimSpace(cmd.DeskripsiRisiko)
	if deskripsi == "" {
		deskripsi = strings.TrimSpace(cmd.PernyataanRisiko)
	}
	if deskripsi == "" {
		deskripsi = prev.DeskripsiRisiko
	}

	kemungkinanAwal := cmd.KemungkinanAwal
	if kemungkinanAwal <= 0 && cmd.Kemungkinan != nil {
		kemungkinanAwal = *cmd.Kemungkinan
	}
	if kemungkinanAwal <= 0 {
		kemungkinanAwal = prev.KemungkinanAwal
	}
	if kemungkinanAwal <= 0 {
		kemungkinanAwal = 3
	}

	dampakAwal := cmd.DampakAwal
	if dampakAwal <= 0 && cmd.Keparahan != nil {
		dampakAwal = *cmd.Keparahan
	}
	if dampakAwal <= 0 {
		dampakAwal = prev.DampakAwal
	}
	if dampakAwal <= 0 {
		dampakAwal = 3
	}

	mitigasi := strings.TrimSpace(cmd.TindakanMitigasi)
	if mitigasi == "" {
		mitigasi = strings.TrimSpace(cmd.RencanaMitigasi)
	}
	if mitigasi == "" {
		mitigasi = prev.TindakanMitigasi
	}

	tingkatAwal := cmd.TingkatRisikoAwal
	if tingkatAwal == "" {
		tingkatAwal = prev.TingkatRisikoAwal
	}

	peluang := cmd.Peluang
	if peluang == "" && prev.Peluang != "" {
		peluang = prev.Peluang
	}

	residuKenyataan := cmd.KenyataanResidu
	if residuKenyataan == "" && prev.KenyataanResidu != "" {
		residuKenyataan = prev.KenyataanResidu
	}

	kemungkinanResidu := cmd.KemungkinanResidu
	if kemungkinanResidu == nil {
		kemungkinanResidu = prev.KemungkinanResidu
	}

	dampakResidu := cmd.DampakResidu
	if dampakResidu == nil {
		dampakResidu = prev.DampakResidu
	}

	tingkatResidu := cmd.TingkatRisikoResidu
	if tingkatResidu == nil {
		tingkatResidu = prev.TingkatRisikoResidu
	}

	rtl := cmd.RencanaTindakLanjut
	if rtl == "" && prev.RencanaTindakLanjut != "" {
		rtl = prev.RencanaTindakLanjut
	}

	status := cmd.Status
	if status == "" {
		status = prev.Status
	}

	result := domainRisk.UpdateRiskAssessment(
		prev,
		parsedUuid,
		fuID,
		tahunID,
		tahunStr,
		idMataProgram,
		isu,
		deskripsi,
		kemungkinanAwal,
		dampakAwal,
		tingkatAwal,
		mitigasi,
		peluang,
		residuKenyataan,
		kemungkinanResidu,
		dampakResidu,
		tingkatResidu,
		rtl,
		status,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	item := result.Value
	if err := h.Repo.Update(ctx, item); err != nil {
		return "", err
	}

	return item.UUID.String(), nil
}
