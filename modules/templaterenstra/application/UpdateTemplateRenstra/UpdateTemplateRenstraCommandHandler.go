package application

import (
	"context"
	"golang.org/x/sync/errgroup"
	"github.com/google/uuid"
	
	domaintemplaterenstra "UnpakSiamida/modules/templaterenstra/domain"
	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	"errors"
    "gorm.io/gorm"
)

type UpdateTemplateRenstraCommandHandler struct {
	Repo                	domaintemplaterenstra.ITemplateRenstraRepository
	IndikatorRenstraRepo    domainindikatorrenstra.IIndikatorRenstraRepository
	FakultasUnitRepo    	domainfakultasunit.IFakultasUnitRepository
}

func (h *UpdateTemplateRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd UpdateTemplateRenstraCommand,
) (string, error) {

	// -------------------------
	// VALIDATE UUID
	// -------------------------
	templaterenstraUUID, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		return "", domaintemplaterenstra.InvalidUuid()
	}
	uuidIndikator, err := uuid.Parse(cmd.Indikator)
	if err != nil {
		return "", domaintemplaterenstra.IndikatorNotFound()
	}
	uuidFakultasUnit, err := uuid.Parse(cmd.FakultasUnit)
	if err != nil {
		return "", domaintemplaterenstra.FakultasUnitNotFound()
	}

	// -------------------------
	// GET EXISTING templaterenstra
	// -------------------------
	var (
		indikatorDefault     		 *domainindikatorrenstra.IndikatorRenstraDefault
		fakultasunitDefault  		 *domainfakultasunit.FakultasUnit
		existingTemplateRenstra 	 *domaintemplaterenstra.TemplateRenstra
	)

	g, gctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		r, err := h.IndikatorRenstraRepo.GetDefaultByUuid(gctx, uuidIndikator)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domaintemplaterenstra.IndikatorNotFound()
			}
			return err;
		}
		indikatorDefault = r
		return nil
	})

	g.Go(func() error {
		r, err := h.FakultasUnitRepo.GetDefaultByUuid(gctx, uuidFakultasUnit)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domaintemplaterenstra.FakultasUnitNotFound()
			}
			return err
		}
		fakultasunitDefault = r
		return nil
	})

	g.Go(func() error {
		r, err := h.Repo.GetByUuid(ctx, templaterenstraUUID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domaintemplaterenstra.NotFound(cmd.Uuid)
			}
			return err
		}
		existingTemplateRenstra = r
		return nil
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	// -------------------------
	// AGGREGATE ROOT LOGIC
	// -------------------------
	result := domaintemplaterenstra.UpdateTemplateRenstra(
		existingTemplateRenstra,
		templaterenstraUUID,
		cmd.Tahun,
		indikatorDefault.Id,
		cmd.IsPertanyaan=="1",
		fakultasunitDefault.ID,
		cmd.Kategori,
		cmd.Klasifikasi,
		cmd.Satuan,
		cmd.Target,
		cmd.TargetMin,
		cmd.TargetMax,
		cmd.Tugas,
	)

	if !result.IsSuccess {
		return "", result.Error
	}

	updatedTemplateRenstra := result.Value

	// -------------------------
	// SAVE TO REPOSITORY
	// -------------------------
	if err := h.Repo.Update(ctx, updatedTemplateRenstra); err != nil {
		return "", err
	}

	return updatedTemplateRenstra.UUID.String(), nil
}
