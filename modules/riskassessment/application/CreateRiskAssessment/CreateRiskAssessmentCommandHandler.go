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

type CreateRiskAssessmentCommandHandler struct {
	Repo             domainRisk.IRiskAssessmentRepository
	RepoFakultasUnit domainfakultasunit.IFakultasUnitRepository
	DB               *gorm.DB
}

func (h *CreateRiskAssessmentCommandHandler) Handle(
	ctx context.Context,
	cmd CreateRiskAssessmentCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var fuID *int = cmd.IdFakultasUnit
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
	var tahunStr string = strings.TrimSpace(cmd.Tahun)

	if strings.TrimSpace(cmd.TahunUuid) != "" {
		if parsed, err := uuid.Parse(strings.TrimSpace(cmd.TahunUuid)); err == nil {
			var mt struct {
				ID    int
				Tahun string
			}
			if err := h.DB.WithContext(ctx).Table("master_tahun").Select("id, tahun").Where("uuid = ?", parsed).Take(&mt).Error; err == nil {
				tahunID = &mt.ID
				if tahunStr == "" {
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

	// Fallbacks from legacy fields
	isu := strings.TrimSpace(cmd.IsuInternalEksternal)
	if isu == "" {
		isu = strings.TrimSpace(cmd.Penyebab)
	}

	deskripsi := strings.TrimSpace(cmd.DeskripsiRisiko)
	if deskripsi == "" {
		deskripsi = strings.TrimSpace(cmd.PernyataanRisiko)
	}

	kemungkinanAwal := cmd.KemungkinanAwal
	if kemungkinanAwal <= 0 && cmd.Kemungkinan != nil {
		kemungkinanAwal = *cmd.Kemungkinan
	}
	if kemungkinanAwal <= 0 {
		kemungkinanAwal = 3
	}

	dampakAwal := cmd.DampakAwal
	if dampakAwal <= 0 && cmd.Keparahan != nil {
		dampakAwal = *cmd.Keparahan
	}
	if dampakAwal <= 0 {
		dampakAwal = 3
	}

	mitigasi := strings.TrimSpace(cmd.TindakanMitigasi)
	if mitigasi == "" {
		mitigasi = strings.TrimSpace(cmd.RencanaMitigasi)
	}

	result := domainRisk.NewRiskAssessment(
		fuID,
		tahunID,
		tahunStr,
		idMataProgram,
		isu,
		deskripsi,
		kemungkinanAwal,
		dampakAwal,
		cmd.TingkatRisikoAwal,
		mitigasi,
		cmd.Peluang,
		cmd.KenyataanResidu,
		cmd.KemungkinanResidu,
		cmd.DampakResidu,
		cmd.TingkatRisikoResidu,
		cmd.RencanaTindakLanjut,
		cmd.Status,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	item := result.Value
	if err := h.Repo.Create(ctx, item); err != nil {
		return "", err
	}

	return item.UUID.String(), nil
}
