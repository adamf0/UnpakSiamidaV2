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
		"Minimal satu Fakultas / Prodi / Unit (target) wajib dipilih",
	)
}

func InvalidTargets() commondomain.Error {
	return commondomain.FailureError(
		"MataProgramV2.InvalidTargets",
		"Gagal menyimpan data mata program v2. Pastikan target unit valid.",
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
