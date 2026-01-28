package application

import (
	"context"
	"errors"

	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	domainjadwalproker "UnpakSiamida/modules/jadwalproker/domain"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type UpdateJadwalProkerCommandHandler struct {
	Repo             domainjadwalproker.IJadwalProkerRepository
	RepoFakultasUnit domainfakultasunit.IFakultasUnitRepository
}

func (h *UpdateJadwalProkerCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateJadwalProkerCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	jadwalprokerUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainjadwalproker.InvalidUuid()
	}

	uuidFakultas, err := uuid.Parse(cmd.FakultasUuid)
	if err != nil {
		return "", domainjadwalproker.InvalidFakultas()
	}

	// -------------------------
	// GET EXISTING jadwalproker
	// -------------------------
	var (
		existingJadwalProker *domainjadwalproker.JadwalProker
		fakultas             *domainfakultasunit.FakultasUnit
	)

	g, ctxg := errgroup.WithContext(ctx)

	// query JadwalProker
	g.Go(func() error {
		var err error
		existingJadwalProker, err = h.Repo.GetByUuid(ctxg, jadwalprokerUUID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domainjadwalproker.NotFound(cmd.Uuid)
			}
			return err
		}
		return nil
	})

	// query FakultasUnit
	g.Go(func() error {
		var err error
		fakultas, err = h.RepoFakultasUnit.GetDefaultByUuid(ctxg, uuidFakultas)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domainjadwalproker.NotFoundFakultas()
			}
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return "", err
	}
	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domainjadwalproker.UpdateJadwalProker(
		existingJadwalProker,
		jadwalprokerUUID,
		fakultas.ID,
		cmd.TanggalTutupEntry,
		cmd.TanggalTutupDokumen,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedJadwalProker := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedJadwalProker); err != nil {
		return "", err
	}

	return updatedJadwalProker.UUID.String(), nil
}
