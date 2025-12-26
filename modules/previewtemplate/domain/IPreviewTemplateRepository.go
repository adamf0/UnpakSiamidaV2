package domain

import (
	"context"
)

type IPreviewTemplateRepository interface {
	GetByTahunFakultasUnit(ctx context.Context, tahun string, fakultasUnit string) ([]PreviewTemplate, error)
}
