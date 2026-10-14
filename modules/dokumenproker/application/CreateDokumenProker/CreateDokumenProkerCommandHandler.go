package application

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	domaindokumenproker "UnpakSiamida/modules/dokumenproker/domain"
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	domainmataprogram "UnpakSiamida/modules/mataprogram/domain"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type CreateDokumenProkerCommandHandler struct {
	Repo             domaindokumenproker.IDokumenProkerRepository
	RepoMataProgram  domainmataprogram.IMataProgramRepository
	RepoFakultasUnit domainfakultasunit.IFakultasUnitRepository
}

func (h *CreateDokumenProkerCommandHandler) Handle(
	ctx context.Context,
	cmd CreateDokumenProkerCommand,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var fakultasID uint
	uuidFakultas, errFakultasUuid := uuid.Parse(strings.TrimSpace(cmd.FakultasUuid))
	if errFakultasUuid != nil {
		parsedID, errConv := strconv.ParseUint(strings.TrimSpace(cmd.FakultasUuid), 10, 64)
		if errConv != nil || parsedID == 0 {
			return "", domaindokumenproker.InvalidFakultas()
		}
		fakultasID = uint(parsedID)
	}

	uuidMataProgram, err := uuid.Parse(cmd.MataProgramUuid)
	if err != nil {
		return "", domaindokumenproker.InvalidMataProgram()
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
					return domaindokumenproker.NotFoundFakultas()
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
				return domaindokumenproker.NotFoundFakultas()
			}
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	result := domaindokumenproker.NewDokumenProker(
		mataprogram.Id,
		fakultasID,
		cmd.JenisDokumen,
		cmd.File,
		cmd.Status,
		cmd.Catatan,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	createDokumenProker := result.Value
	if err := h.Repo.Create(ctx, createDokumenProker); err != nil {
		return "", err
	}

	return result.Value.UUID.String(), nil
}
