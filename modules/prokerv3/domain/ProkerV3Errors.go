package domain

import (
	commondomain "UnpakSiamida/common/domain"
	"fmt"
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

func AktivitasNotFound() commondomain.Error {
	return commondomain.NotFoundError(
		"ProkerV3.AktivitasNotFound",
		"Aktivitas tidak ditemukan",
	)
}

func NegativeAnggaranTimeline() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.NegativeAnggaranTimeline",
		"Nominal anggaran tahapan timeline tidak boleh negatif",
	)
}

func AnggaranTimelineExceedsRealisasi(anggaran, maxAnggaran float64) commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.AnggaranTimelineExceedsRealisasi",
		fmt.Sprintf("Nominal anggaran tahapan (Rp %.0f) tidak boleh melebihi anggaran realisasi kegiatan (Rp %.0f)", anggaran, maxAnggaran),
	)
}

func InvalidTimelineUuid() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.InvalidTimelineUuid",
		"Invalid timeline UUID",
	)
}

func InvalidAktivitasUuid() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.InvalidAktivitasUuid",
		"Invalid aktivitas UUID",
	)
}
