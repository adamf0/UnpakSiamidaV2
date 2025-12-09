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
// import "encoding/json"
// import "fmt"

type CreateTemplateRenstraCommandHandler struct {
	Repo                	domaintemplaterenstra.ITemplateRenstraRepository
	IndikatorRenstraRepo    domainindikatorrenstra.IIndikatorRenstraRepository
	FakultasUnitRepo    	domainfakultasunit.IFakultasUnitRepository
}

func (h *CreateTemplateRenstraCommandHandler) Handle(
	ctx context.Context,
	cmd CreateTemplateRenstraCommand,
) (string, error) {

	uuidIndikator, err := uuid.Parse(cmd.Indikator)
	if err != nil {
		return "", domaintemplaterenstra.IndikatorNotFound()
	}

	uuidFakultasUnit, err := uuid.Parse(cmd.FakultasUnit)
	if err != nil {
		return "", domaintemplaterenstra.FakultasUnitNotFound()
	}

	var (
		indikatorDefault     *domainindikatorrenstra.IndikatorRenstraDefault
		fakultasunitDefault  *domainfakultasunit.FakultasUnit
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

	if err := g.Wait(); err != nil {
		return "", err
	}

	// b, _ := json.MarshalIndent(fakultasunitDefault, "", "  ")
	// fmt.Println("DEBUG fakultasunitDefault:", string(b))
	
	// c, _ := json.MarshalIndent(indikatorDefault, "", "  ")
	// fmt.Println("DEBUG indikatorDefault:", string(c))

	result := domaintemplaterenstra.NewTemplateRenstra(
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

	templateRenstra := result.Value

	// --------------------------
	// SAVE REPOSITORY
	// --------------------------
	if err := h.Repo.Create(ctx, templateRenstra); err != nil {
		return "", err
	}

	return templateRenstra.UUID.String(), nil
}
