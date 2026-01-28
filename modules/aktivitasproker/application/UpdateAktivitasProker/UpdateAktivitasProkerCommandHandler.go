package application

import (
	"context"
	"errors"

	domainaktivitasproker "UnpakSiamida/modules/aktivitasproker/domain"
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	domainmataprogram "UnpakSiamida/modules/mataprogram/domain"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type UpdateAktivitasProkerCommandHandler struct {
	Repo             domainaktivitasproker.IAktivitasProkerRepository
	RepoMataProgram  domainmataprogram.IMataProgramRepository
	RepoFakultasUnit domainfakultasunit.IFakultasUnitRepository
}

func (h *UpdateAktivitasProkerCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateAktivitasProkerCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	uuidAktivitasProker, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domainaktivitasproker.InvalidUuid()
	}

	uuidFakultas, err := uuid.Parse(cmd.FakultasUuid)
	if err != nil {
		return "", domainaktivitasproker.InvalidFakultas()
	}

	uuidMataProgram, err := uuid.Parse(cmd.MataProgramUuid)
	if err != nil {
		return "", domainaktivitasproker.InvalidMataProgram()
	}

	// -------------------------
	// GET EXISTING Aktivitasproker
	// -------------------------
	var (
		existingAktivitasProker *domainaktivitasproker.AktivitasProker
		fakultas                *domainfakultasunit.FakultasUnit
		mataprogram             *domainmataprogram.MataProgramDefault
	)

	g, ctxg := errgroup.WithContext(ctx)

	// query AktivitasProker
	g.Go(func() error {
		var err error
		existingAktivitasProker, err = h.Repo.GetByUuid(ctxg, uuidAktivitasProker)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domainaktivitasproker.NotFound(cmd.Uuid)
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
				return domainaktivitasproker.NotFoundFakultas()
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
				return domainaktivitasproker.NotFoundFakultas()
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
	result := domainaktivitasproker.UpdateAktivitasProker(
		existingAktivitasProker,
		uuidAktivitasProker,
		mataprogram.Id,
		fakultas.ID,
		cmd.Aktivitas,
		cmd.PIC,
		cmd.TanggalRKAwal,
		cmd.TanggalRKAkhir,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedAktivitasProker := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedAktivitasProker); err != nil {
		return "", err
	}

	return updatedAktivitasProker.UUID.String(), nil
}
