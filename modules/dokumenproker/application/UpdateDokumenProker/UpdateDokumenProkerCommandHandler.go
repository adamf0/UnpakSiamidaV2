package application

import (
	"context"
	"errors"

	domaindokumenproker "UnpakSiamida/modules/dokumenproker/domain"
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	domainmataprogram "UnpakSiamida/modules/mataprogram/domain"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type UpdateDokumenProkerCommandHandler struct {
	Repo             domaindokumenproker.IDokumenProkerRepository
	RepoMataProgram  domainmataprogram.IMataProgramRepository
	RepoFakultasUnit domainfakultasunit.IFakultasUnitRepository
}

func (h *UpdateDokumenProkerCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateDokumenProkerCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	uuidDokumenProker, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domaindokumenproker.InvalidUuid()
	}

	uuidFakultas, err := uuid.Parse(cmd.FakultasUuid)
	if err != nil {
		return "", domaindokumenproker.InvalidFakultas()
	}

	uuidMataProgram, err := uuid.Parse(cmd.MataProgramUuid)
	if err != nil {
		return "", domaindokumenproker.InvalidMataProgram()
	}

	var (
		existingDokumenProker *domaindokumenproker.DokumenProker
		fakultas              *domainfakultasunit.FakultasUnit
		mataprogram           *domainmataprogram.MataProgramDefault
	)

	g, ctxg := errgroup.WithContext(ctx)

	// query DokumenProker
	g.Go(func() error {
		var err error
		existingDokumenProker, err = h.Repo.GetByUuid(ctxg, uuidDokumenProker)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domaindokumenproker.NotFound(cmd.Uuid)
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
				return domaindokumenproker.NotFoundFakultas()
			}
			return err
		}
		return nil
	})

	g.Go(func() error {
		var err error
		mataprogram, err = h.RepoMataProgram.GetDefaultByUuid(ctxg, uuidMataProgram)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domaindokumenproker.NotFoundFakultas()
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
	result := domaindokumenproker.UpdateDokumenProker(
		existingDokumenProker,
		uuidDokumenProker,
		mataprogram.Id,
		fakultas.ID,
		cmd.JenisDokumen,
		cmd.File,
		cmd.Status,
		cmd.Catatan,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedDokumenProker := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedDokumenProker); err != nil {
		return "", err
	}

	return updatedDokumenProker.UUID.String(), nil
}
