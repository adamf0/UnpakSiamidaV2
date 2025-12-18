package domain

import (
	"UnpakSiamida/common/domain"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("PreviewTemplate.EmptyData", "data is not found")
}

func NotFoundTreeIndikator() domain.Error {
	return domain.NotFoundError("PreviewTemplate.NotFoundTreeIndikator", "Tree Indikator not found")
}

func NotFound() domain.Error {
	return domain.NotFoundError("PreviewTemplate.NotFound", "Preview Template not found")
}
