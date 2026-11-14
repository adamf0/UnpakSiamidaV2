package domain

import (
	commondomain "UnpakSiamida/common/domain"
	"fmt"
)

var (
	AktivitasNotFoundError = commondomain.Error{
		Code:        "Aktivitas.NotFound",
		Description: "Aktivitas tidak ditemukan",
		Type:        commondomain.NotFound,
	}

	TimelineNotFoundError = commondomain.Error{
		Code:        "Timeline.NotFound",
		Description: "Timeline tidak ditemukan",
		Type:        commondomain.NotFound,
	}
)

func TimelineNotFound() commondomain.Error {
	return commondomain.NotFoundError(
		"AktivitasV2.TimelineNotFound",
		"Timeline aktivitas tidak ditemukan",
	)
}

func FileNotFound() commondomain.Error {
	return commondomain.NotFoundError(
		"AktivitasV2.FileNotFound",
		"File timeline tidak ditemukan",
	)
}

func AktivitasNotFound() commondomain.Error {
	return commondomain.NotFoundError(
		"AktivitasV2.AktivitasNotFound",
		"Aktivitas tidak ditemukan",
	)
}

func NegativeAnggaranTimeline() commondomain.Error {
	return commondomain.FailureError(
		"AktivitasV2.NegativeAnggaranTimeline",
		"Nominal anggaran tahapan timeline tidak boleh negatif",
	)
}

func AnggaranTimelineExceedsRealisasi(anggaran, maxAnggaran float64) commondomain.Error {
	return commondomain.FailureError(
		"AktivitasV2.AnggaranTimelineExceedsRealisasi",
		fmt.Sprintf("Nominal anggaran tahapan (Rp %.0f) tidak boleh melebihi anggaran realisasi kegiatan (Rp %.0f)", anggaran, maxAnggaran),
	)
}

func InvalidTimelineUuid() commondomain.Error {
	return commondomain.FailureError(
		"AktivitasV2.InvalidTimelineUuid",
		"Invalid timeline UUID",
	)
}

func InvalidAktivitasUuid() commondomain.Error {
	return commondomain.FailureError(
		"AktivitasV2.InvalidAktivitasUuid",
		"Invalid aktivitas UUID",
	)
}

func InvalidMatrixUuid() commondomain.Error {
	return commondomain.FailureError(
		"AktivitasV2.InvalidMatrixUuid",
		"Invalid matrix UUID",
	)
}

func ProkerNotFound() commondomain.Error {
	return commondomain.NotFoundError(
		"AktivitasV2.ProkerNotFound",
		"Program kerja tidak ditemukan",
	)
}

func ProkerDeleted() commondomain.Error {
	return commondomain.FailureError(
		"AktivitasV2.ProkerDeleted",
		"Program kerja sudah dihapus dan tidak dapat ditambah aktivitas",
	)
}

func NamaAktivitasRequired() commondomain.Error {
	return commondomain.FailureError(
		"AktivitasV2.NamaAktivitasRequired",
		"Nama aktivitas / kegiatan wajib diisi",
	)
}

func InvalidFileUuid() commondomain.Error {
	return commondomain.FailureError(
		"AktivitasV2.InvalidFileUuid",
		"Invalid file UUID",
	)
}
