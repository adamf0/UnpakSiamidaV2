package application

import (
	"context"
	"strconv"
	"strings"
	"time"

	domainmataprogram "UnpakSiamida/modules/mataprogram/domain"
	domaintahunproker "UnpakSiamida/modules/tahunproker/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateMataProgramCommandHandler struct {
	Repo            domainmataprogram.IMataProgramRepository
	RepoTahunProker domaintahunproker.ITahunProkerRepository
	DB              *gorm.DB
}

func (h *CreateMataProgramCommandHandler) Handle(
	ctx context.Context,
	cmd CreateMataProgramCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tahunprokerUUID, err := uuid.Parse(cmd.TahunUuid)
	if err != nil {
		return "", domainmataprogram.InvalidParseTahun()
	}

	var tahunprokerId uint
	tahunproker, err := h.RepoTahunProker.GetByUuid(ctx, tahunprokerUUID)
	if err != nil || tahunproker == nil {
		tahunprokerId = 0
	} else {
		tahunprokerId = tahunproker.ID
	}

	var fuID uint = cmd.IdFakultasUnit
	if fuID == 0 && strings.TrimSpace(cmd.FakultasUnitUuid) != "" {
		trimmed := strings.TrimSpace(cmd.FakultasUnitUuid)
		if parsed, err := uuid.Parse(trimmed); err == nil && h.DB != nil {
			var fu struct {
				ID uint
			}
			if err := h.DB.WithContext(ctx).Table("v_fakultas_unit").Select("id").Where("uuid = ?", parsed).Take(&fu).Error; err == nil {
				fuID = fu.ID
			}
		} else if id, err := strconv.Atoi(trimmed); err == nil && id > 0 {
			fuID = uint(id)
		}
	}

	result := domainmataprogram.NewMataProgram(
		tahunprokerId,
		cmd.MataProgram,
		fuID,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	createMataProgram := result.Value
	if err := h.Repo.Create(ctx, createMataProgram); err != nil {
		return "", err
	}

	return result.Value.UUID.String(), nil
}
