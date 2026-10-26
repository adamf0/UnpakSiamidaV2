package application

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	domainmataprogram "UnpakSiamida/modules/mataprogram/domain"
	domaintahunproker "UnpakSiamida/modules/tahunproker/domain"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type UpdateMataProgramCommandHandler struct {
	Repo            domainmataprogram.IMataProgramRepository
	RepoTahunProker domaintahunproker.ITahunProkerRepository
	DB              *gorm.DB
}

func (h *UpdateMataProgramCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateMataProgramCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	mataprogramUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainmataprogram.InvalidUuid()
	}
	tahunprokerUUID, err := uuid.Parse(cmd.TahunUuid)
	if err != nil {
		return "", domainmataprogram.InvalidParseTahun()
	}

	// -------------------------
	// GET EXISTING mataprogram
	// -------------------------
	var (
		existingMataProgram *domainmataprogram.MataProgram
		tahunprokerId       uint
	)

	g, ctxg := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		existingMataProgram, err = h.Repo.GetByUuid(ctxg, mataprogramUUID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domainmataprogram.NotFound(cmd.Uuid)
			}
			return err
		}
		return nil
	})

	g.Go(func() error {
		tahunproker, err := h.RepoTahunProker.GetByUuid(ctxg, tahunprokerUUID)
		if err != nil {
			tahunprokerId = 0
		} else {
			tahunprokerId = tahunproker.ID
		}
		return nil
	})

	// tunggu semua selesai
	if err := g.Wait(); err != nil {
		return "", err
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
	if fuID == 0 && existingMataProgram != nil {
		fuID = existingMataProgram.FakultasUnitId
	}

	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domainmataprogram.UpdateMataProgram(
		existingMataProgram,
		mataprogramUUID,
		tahunprokerId,
		cmd.MataProgram,
		fuID,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedMataProgram := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedMataProgram); err != nil {
		return "", err
	}

	return updatedMataProgram.UUID.String(), nil
}
