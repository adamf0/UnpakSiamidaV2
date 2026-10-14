package application

import (
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	domainSwot "UnpakSiamida/modules/swot/domain"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateSwotCommandHandler struct {
	Repo             domainSwot.ISwotRepository
	RepoFakultasUnit domainfakultasunit.IFakultasUnitRepository
	DB               *gorm.DB
}

func (h *UpdateSwotCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateSwotCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parsedUID, err := uuid.Parse(strings.TrimSpace(cmd.UUID))
	if err != nil {
		return "", domainSwot.InvalidData()
	}

	existing, err := h.Repo.GetByUuid(ctx, parsedUID)
	if err != nil || existing == nil {
		return "", domainSwot.EmptyData()
	}

	var fuID *int = existing.IdFakultasUnit
	if cmd.IdFakultasUnit != nil {
		fuID = cmd.IdFakultasUnit
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

	var tahunID *int = existing.IdMasterTahun
	if cmd.IdMasterTahun != nil {
		tahunID = cmd.IdMasterTahun
	}
	var tahunStr string = existing.Tahun
	if strings.TrimSpace(cmd.Tahun) != "" {
		tahunStr = strings.TrimSpace(cmd.Tahun)
	}

	if strings.TrimSpace(cmd.TahunUuid) != "" {
		if parsed, err := uuid.Parse(strings.TrimSpace(cmd.TahunUuid)); err == nil {
			var mt struct {
				ID    int
				Tahun string
			}
			if err := h.DB.WithContext(ctx).Table("master_tahun").Select("id, tahun").Where("uuid = ?", parsed).Take(&mt).Error; err == nil {
				tahunID = &mt.ID
				tahunStr = mt.Tahun
			}
		} else if id, err := strconv.Atoi(strings.TrimSpace(cmd.TahunUuid)); err == nil && id > 0 {
			tahunID = &id
		}
	}

	var skorVal *float64 = cmd.Skor
	if skorVal == nil && cmd.Bobot != nil && cmd.Rating != nil {
		s := (*cmd.Bobot) * float64(*cmd.Rating)
		skorVal = &s
	}

	result := domainSwot.UpdateSwot(
		existing,
		parsedUID,
		fuID,
		tahunID,
		tahunStr,
		cmd.Kategori,
		cmd.Faktor,
		cmd.Deskripsi,
		cmd.Strategi,
		cmd.Bobot,
		cmd.Rating,
		skorVal,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	if err := h.Repo.Update(ctx, existing); err != nil {
		return "", err
	}

	return existing.UUID.String(), nil
}
