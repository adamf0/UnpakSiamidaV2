package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("TemplateRenstra.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("TemplateRenstra.InvalidUuid", "uuid is invalid")
}

func IndikatorNotFound() domain.Error {
	return domain.NotFoundError("TemplateRenstra.IndikatorNotFound", "indikator not found")
}

func FakultasUnitNotFound() domain.Error {
	return domain.NotFoundError("TemplateRenstra.FakultasUnitNotFound", "fakultas unit not found")
}

func InvalidValueTarget() domain.Error {
	return domain.NotFoundError("TemplateRenstra.InvalidValueTarget", "invalid target combination, either provide Target only or provide both TargetMin and TargetMax")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("TemplateRenstra.InvalidData", "data is invalid")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("TemplateRenstra.NotFound", fmt.Sprintf("TemplateRenstra with identifier %s not found", id))
}

func DuplicateData() domain.Error {
	return domain.NotFoundError("TemplateRenstra.DuplicateData", "data not alowed duplicate")
}
