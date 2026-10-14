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

	existingDokumenProker, err := h.Repo.GetByUuid(ctx, uuidDokumenProker)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domaindokumenproker.NotFound(cmd.Uuid)
		}
		return "", err
	}

	fakultasID := existingDokumenProker.FakultasUnit
	mataProgramID := existingDokumenProker.MataProgram
	jenisDokumen := existingDokumenProker.JenisDokumen
	fileDokumen := existingDokumenProker.File
	status := existingDokumenProker.Status
	catatan := existingDokumenProker.Catatan

	g, ctxg := errgroup.WithContext(ctx)

	// Optional FakultasUnit update
	if strings.TrimSpace(cmd.FakultasUuid) != "" {
		uuidFakultas, errFakultasUuid := uuid.Parse(strings.TrimSpace(cmd.FakultasUuid))
		if errFakultasUuid != nil {
			parsedID, errConv := strconv.ParseUint(strings.TrimSpace(cmd.FakultasUuid), 10, 64)
			if errConv == nil && parsedID > 0 {
				fakultasID = uint(parsedID)
			}
		} else {
			g.Go(func() error {
				fak, errFak := h.RepoFakultasUnit.GetDefaultByUuid(ctxg, uuidFakultas)
				if errFak != nil {
					if errors.Is(errFak, gorm.ErrRecordNotFound) {
						return domaindokumenproker.NotFoundFakultas()
					}
					return errFak
				}
				fakultasID = fak.ID
				return nil
			})
		}
	}

	// Optional MataProgram update
	if strings.TrimSpace(cmd.MataProgramUuid) != "" {
		uuidMataProgram, errMp := uuid.Parse(strings.TrimSpace(cmd.MataProgramUuid))
		if errMp == nil {
			g.Go(func() error {
				mp, errMpFetch := h.RepoMataProgram.GetDefaultByUuid(ctxg, uuidMataProgram)
				if errMpFetch != nil {
					if errors.Is(errMpFetch, gorm.ErrRecordNotFound) {
						return domaindokumenproker.NotFoundMataProgram()
					}
					return errMpFetch
				}
				mataProgramID = mp.Id
				return nil
			})
		}
	}

	if err := g.Wait(); err != nil {
		return "", err
	}

	if strings.TrimSpace(cmd.JenisDokumen) != "" {
		jenisDokumen = cmd.JenisDokumen
	}
	if strings.TrimSpace(cmd.File) != "" {
		fileDokumen = cmd.File
	}
	if strings.TrimSpace(cmd.Status) != "" {
		status = cmd.Status
	}
	if cmd.Catatan != nil {
		catatan = cmd.Catatan
	}

	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domaindokumenproker.UpdateDokumenProker(
		existingDokumenProker,
		uuidDokumenProker,
		mataProgramID,
		fakultasID,
		jenisDokumen,
		fileDokumen,
		status,
		catatan,
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
