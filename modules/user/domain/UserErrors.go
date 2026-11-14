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

func AbacUuidRequired() domain.Error {
	return domain.FailureError("User.AbacUuidRequired", "UUID penugasan wajib disertakan")
}

func AbacInvalidUuid() domain.Error {
	return domain.FailureError("User.AbacInvalidUuid", "Format UUID penugasan tidak valid")
}

func AbacNotFound() domain.Error {
	return domain.NotFoundError("User.AbacNotFound", "Data penugasan tidak ditemukan")
}

func AbacTahunRequired() domain.Error {
	return domain.FailureError("User.AbacTahunRequired", "Tahun penugasan wajib diisi")
}

func AbacUserUuidInvalid() domain.Error {
	return domain.FailureError("User.AbacUserUuidInvalid", "UUID User tidak valid")
}

func AbacUserNotFound() domain.Error {
	return domain.NotFoundError("User.AbacUserNotFound", "User yang dipilih tidak ditemukan dalam sistem")
}

func AbacUserRequired() domain.Error {
	return domain.FailureError("User.AbacUserRequired", "Pengguna wajib dipilih")
}

func AbacAdminNotAllowed() domain.Error {
	return domain.FailureError("User.AbacAdminNotAllowed", "Role admin tidak dapat ditugaskan melalui ABAC, gunakan level akun pengguna")
}

func AbacFakultasUnitNotFound() domain.Error {
	return domain.NotFoundError("User.AbacFakultasUnitNotFound", "Fakultas/Unit/Prodi yang dipilih tidak ditemukan")
}

func AbacTargetRequired() domain.Error {
	return domain.FailureError("User.AbacTargetRequired", "Target Fakultas/Unit/Prodi wajib dipilih untuk role ini")
}

func AbacTargetUpmsInvalid() domain.Error {
	return domain.FailureError("User.AbacTargetUpmsInvalid", "Target untuk role UPMS harus berupa Program Studi (Prodi)")
}

func AbacTargetUpmfInvalid() domain.Error {
	return domain.FailureError("User.AbacTargetUpmfInvalid", "Target untuk role Fakultas / UPMF harus berupa Fakultas")
}

func AbacTargetUpmlInvalid() domain.Error {
	return domain.FailureError("User.AbacTargetUpmlInvalid", "Target untuk role UPML harus berupa Fakultas atau Unit (bukan Program Studi)")
}

func AbacDuplicateAssignment() domain.Error {
	return domain.ConflictError("User.AbacDuplicateAssignment", "Penugasan untuk user, tahun, dan target unit ini sudah ada")
}
