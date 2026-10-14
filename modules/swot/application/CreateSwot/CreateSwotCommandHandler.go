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

type CreateSwotCommandHandler struct {
	Repo             domainSwot.ISwotRepository
	RepoFakultasUnit domainfakultasunit.IFakultasUnitRepository
	DB               *gorm.DB
}

func (h *CreateSwotCommandHandler) Handle(
	ctx context.Context,
	cmd CreateSwotCommand,
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

	// Hitung skor otomatis jika bobot & rating ada namun skor belum diisi
	var skorVal *float64 = cmd.Skor
	if skorVal == nil && cmd.Bobot != nil && cmd.Rating != nil {
		s := (*cmd.Bobot) * float64(*cmd.Rating)
		skorVal = &s
	}

	result := domainSwot.NewSwot(
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

	item := result.Value
	if err := h.Repo.Create(ctx, item); err != nil {
		return "", err
	}

	return item.UUID.String(), nil
}
