package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("User.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("User.InvalidUuid", "uuid is invalid")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("User.InvalidData", "data is invalid")
}

func InvalidEmail() domain.Error {
	return domain.NotFoundError("User.InvalidEmail", "email tidak valid atau tidak diperbolehkan")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("User.NotFound", fmt.Sprintf("User with identifier %s not found", id))
}

func InvalidParsing(target string) domain.Error {
	return domain.NotFoundError("User.IvalidParsing", fmt.Sprintf("failed parsing %s to UUID", target))
}

func NotFoundFakultasUnit(id string) domain.Error {
	return domain.NotFoundError("User.NotFoundFakultasUnit", fmt.Sprintf("Fakultas unit with identifier %s not found", id))
}
