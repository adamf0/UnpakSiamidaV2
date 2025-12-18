package application

import (
    "context"
	"golang.org/x/sync/errgroup"
    domainpreviewtemplate "UnpakSiamida/modules/previewtemplate/domain"
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
    "errors"
    "gorm.io/gorm"
	"github.com/google/uuid"
	"strconv"
)

type GetPreviewTemplateByTahunFakultasUnitQueryHandler struct {
    Repo domainpreviewtemplate.IPreviewTemplateRepository
	RepoFakultasUnit domainfakultasunit.IFakultasUnitRepository
}

func (h *GetPreviewTemplateByTahunFakultasUnitQueryHandler) Handle(
	ctx context.Context,
	q GetPreviewTemplateByTahunFakultasUnitQuery,
) ([]domainpreviewtemplate.PreviewTemplate, error) {
	UuidFakultasUnit, err := uuid.Parse(q.FakultasUnit)
	if err != nil {
		return nil, domainfakultasunit.NotFound(q.FakultasUnit) //[PR] harus pindah ke preview error
	}

	fakultasunit, err := h.RepoFakultasUnit.GetDefaultByUuid(ctx, UuidFakultasUnit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainfakultasunit.NotFound(q.FakultasUnit)
		}
		return nil, err
	}

	var (
		tree    []domainpreviewtemplate.IndikatorTree
		preview []domainpreviewtemplate.PreviewTemplate
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		tree, err = h.Repo.GetIndikatorTree(ctx, q.Tahun)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domainpreviewtemplate.NotFoundTreeIndikator()
			}
			return err
		}
		return nil
	})

	g.Go(func() error {
		var err error
		preview, err = h.Repo.GetByTahunFakultasUnit(ctx, q.Tahun, strconv.FormatUint(uint64(fakultasunit.ID), 10) )
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domainpreviewtemplate.NotFound()
			}
			return err
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
	tree []domainpreviewtemplate.IndikatorTree,
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