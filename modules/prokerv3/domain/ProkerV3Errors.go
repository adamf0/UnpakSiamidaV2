package domain

import (
	commondomain "UnpakSiamida/common/domain"
)

func ProkerNotFound() commondomain.Error {
	return commondomain.NotFoundError(
		"ProkerV3.NotFound",
		"Program kerja tidak ditemukan",
	)
}

func NamaProgramRequired() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.NamaProgramRequired",
		"Nama program kerja wajib diisi",
	)
}

func TimelineNotFound() commondomain.Error {
	return commondomain.NotFoundError(
		"ProkerV3.TimelineNotFound",
		"Timeline aktivitas tidak ditemukan",
	)
}

func FileNotFound() commondomain.Error {
	return commondomain.NotFoundError(
		"ProkerV3.FileNotFound",
		"File timeline tidak ditemukan",
	)
}
