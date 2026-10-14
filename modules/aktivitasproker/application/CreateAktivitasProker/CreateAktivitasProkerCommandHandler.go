package application

import (
	"context"
	"errors"
	"strconv"
	"strings"

	domainaktivitasproker "UnpakSiamida/modules/aktivitasproker/domain"
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	domainmataprogram "UnpakSiamida/modules/mataprogram/domain"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type CreateAktivitasProkerCommandHandler struct {
	Repo             domainaktivitasproker.IAktivitasProkerRepository
	RepoMataProgram  domainmataprogram.IMataProgramRepository
	RepoFakultasUnit domainfakultasunit.IFakultasUnitRepository
}

func (h *CreateAktivitasProkerCommandHandler) Handle(
	ctx context.Context,
	cmd CreateAktivitasProkerCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var fakultasID uint
	uuidFakultas, errFakultasUuid := uuid.Parse(strings.TrimSpace(cmd.FakultasUuid))
	if errFakultasUuid != nil {
		parsedID, errConv := strconv.ParseUint(strings.TrimSpace(cmd.FakultasUuid), 10, 64)
		if errConv != nil || parsedID == 0 {
			return "", domainaktivitasproker.InvalidFakultas()
		}
		fakultasID = uint(parsedID)
	}

	uuidMataProgram, err := uuid.Parse(cmd.MataProgramUuid)
	if err != nil {
		return "", domainaktivitasproker.InvalidMataProgram()
	}

	var (
		fakultas    *domainfakultasunit.FakultasUnit
		mataprogram *domainmataprogram.MataProgramDefault
	)

	g, ctxg := errgroup.WithContext(ctx)

	// query FakultasUnit if UUID was provided
	if errFakultasUuid == nil {
		g.Go(func() error {
			var err error
			fakultas, err = h.RepoFakultasUnit.GetDefaultByUuid(ctxg, uuidFakultas)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return domainaktivitasproker.NotFoundFakultas()
				}
				return err
			}
			fakultasID = fakultas.ID
			return nil
		})
	}

	g.Go(func() error {
		var err error
		mataprogram, err = h.RepoMataProgram.GetDefaultByUuid(ctxg, uuidMataProgram)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domainaktivitasproker.NotFoundMataProgram()
			}
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	result := domainaktivitasproker.NewAktivitasProker(
		mataprogram.Id,
		fakultasID,
		cmd.Aktivitas,
		cmd.PIC,
		cmd.TanggalRKAwal,
		cmd.TanggalRKAkhir,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	createAktivitasProker := result.Value
	if err := h.Repo.Create(ctx, createAktivitasProker); err != nil {
		return "", err
	}

	return result.Value.UUID.String(), nil
}
