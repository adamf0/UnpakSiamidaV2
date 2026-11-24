package domain

import commondomain "UnpakSiamida/common/domain"

func UuidRequired() commondomain.Error {
	return commondomain.FailureError(
		"MataProgramV2.UuidRequired",
		"UUID wajib diisi",
	)
}

func MataProgramRequired() commondomain.Error {
	return commondomain.FailureError(
		"MataProgramV2.MataProgramRequired",
		"Nama Program Kerja (mata_program) wajib diisi",
	)
}

func TargetsRequired() commondomain.Error {
	return commondomain.FailureError(
		"MataProgramV2.TargetsRequired",
		"Minimal satu Template IKU sasaran wajib dipilih",
	)
}

func DuplicateMataProgram() commondomain.Error {
	return commondomain.ConflictError(
		"MataProgramV2.DuplicateMataProgram",
		"Mata program dengan nama dan tahun yang sama sudah ada",
	)
}

func InvalidTemplates() commondomain.Error {
	return commondomain.FailureError(
		"MataProgramV2.InvalidTemplates",
		"Template IKU sasaran yang dipilih tidak valid atau tidak sesuai dengan tahun mata program",
	)
}

func NotFound() commondomain.Error {
	return commondomain.NotFoundError(
		"MataProgramV2.NotFound",
		"Data Mata Program V2 tidak ditemukan",
	)
}

func NotFoundOrDeleted() commondomain.Error {
	return commondomain.NotFoundError(
		"MataProgramV2.NotFound",
		"Data Mata Program V2 tidak ditemukan atau sudah dihapus",
	)
}

func RepoNotInitialized() commondomain.Error {
	return commondomain.FailureError(
		"MataProgramV2.RepoNotInitialized",
		"repository IMataProgramV2Repository is not initialized",
	)
}
