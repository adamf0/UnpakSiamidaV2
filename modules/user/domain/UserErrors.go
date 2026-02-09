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

func InvalidFakultasUnit() domain.Error {
	return domain.NotFoundError("User.InvalidFakultasUnit", "fakultas unit tidak valid")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("User.NotFound", fmt.Sprintf("user with identifier %s not found", id))
}

func InvalidParsing(target string) domain.Error {
	return domain.NotFoundError("User.InvalidParsing", fmt.Sprintf("failed parsing %s to UUID", target))
}

func NotFoundFakultasUnit(id string) domain.Error {
	return domain.NotFoundError("User.NotFoundFakultasUnit", fmt.Sprintf("fakultas unit with identifier %s not found", id))
}

func NotGranted() domain.Error {
	return domain.NotFoundError("User.NotGranted", "you can't access this file because don't have permission")
}
func NotGrantedNonUser() domain.Error {
	return domain.NotFoundError("User.NotGrantedNonUser", "you can't access this file because don't have permission non user")
}
