package application

import (
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
	domainpreviewtemplate "UnpakSiamida/modules/previewtemplate/domain"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type GetPreviewTemplateByTahunFakultasUnitQueryHandler struct {
	Repo             domainpreviewtemplate.IPreviewTemplateRepository
	RepoIndikator    domainindikatorrenstra.IIndikatorRenstraRepository
	RepoFakultasUnit domainfakultasunit.IFakultasUnitRepository
}

func (h *GetPreviewTemplateByTahunFakultasUnitQueryHandler) Handle(
	ctx context.Context,
	q GetPreviewTemplateByTahunFakultasUnitQuery,
) ([]domainpreviewtemplate.PreviewTemplate, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	UuidFakultasUnit, err := uuid.Parse(q.FakultasUnit)
	if err != nil {
		return nil, domainpreviewtemplate.NotFoundFakultasUnit(q.FakultasUnit)
	}

	fakultasunit, err := h.RepoFakultasUnit.GetDefaultByUuid(ctx, UuidFakultasUnit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainpreviewtemplate.NotFoundFakultasUnit(q.FakultasUnit)
		}
		return nil, err
	}

	var (
		tree    []domainindikatorrenstra.IndikatorTree
		preview []domainpreviewtemplate.PreviewTemplate
	)

	g, ctxg := errgroup.WithContext(ctx)

	g.Go(func() error {
		// var err error
		tree, err = h.RepoIndikator.GetIndikatorTree(ctxg, q.Tahun)
		// if err != nil {
		// 	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 		return domainpreviewtemplate.NotFoundTreeIndikator()
		// 	}
		// 	return err
		// }
		return err
	})

	g.Go(func() error {
		// var err error
		if q.Tipe == "renstra" {
			preview, err = h.Repo.GetByTahunFakultasUnit(ctxg, q.Tahun, strconv.FormatUint(uint64(fakultasunit.ID), 10))
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return domainpreviewtemplate.NotFound()
				}
				return err
			}
		} else {
			preview, err = h.Repo.GetByTahunTag(ctxg, q.Tahun, fmt.Sprintf("%s#all", fakultasunit.Type))
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return domainpreviewtemplate.NotFound()
				}
				return err
			}
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return mapPointing(preview, tree), nil
}

func mapPointing(
	preview []domainpreviewtemplate.PreviewTemplate,
	tree []domainindikatorrenstra.IndikatorTree,
) []domainpreviewtemplate.PreviewTemplate {

	pointMap := make(map[int]string)

	for _, t := range tree {
		pointMap[t.IndikatorId] = t.Pointing
	}

	for i := range preview {
		if p, ok := pointMap[preview[i].IndikatorId]; ok {
			preview[i].Pointing = p
		}
	}

	return preview
}
