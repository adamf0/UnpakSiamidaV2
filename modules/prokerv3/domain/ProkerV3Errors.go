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

func DuplicateProkerUnitTahun() commondomain.Error {
	return commondomain.ConflictError(
		"ProkerV3.DuplicateProkerUnitTahun",
		"Program kerja tersebut sudah pernah diajukan untuk fakultas/unit dan tahun yang sama",
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

func AktivitasKosong() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.AktivitasKosong",
		"Program kerja belum memiliki aktivitas. Silakan tambahkan minimal 1 aktivitas sebelum mengajukan.",
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

func DeleteReasonRequired() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.DeleteReasonRequired",
		"Alasan penghapusan program kerja wajib diisi",
	)
}

func AlreadyDeleted() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.AlreadyDeleted",
		"Program kerja sudah berada di tab Buang",
	)
}

func InvalidDocUuid() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.InvalidDocUuid",
		"Invalid document UUID",
	)
}

func DokumenEmpty() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.DokumenEmpty",
		"Dokumen data tidak boleh kosong",
	)
}

func DokumenNotFound() commondomain.Error {
	return commondomain.NotFoundError(
		"ProkerV3.DokumenNotFound",
		"Dokumen tidak ditemukan",
	)
}

func DokumenProkerNotFound() commondomain.Error {
	return commondomain.NotFoundError(
		"ProkerV3.DokumenProkerNotFound",
		"Data program kerja dokumen tidak ditemukan",
	)
}

func VerifikasiDokumenForbidden() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.VerifikasiDokumenForbidden",
		"Hanya level admin atau fakultas yang dapat memverifikasi dokumen",
	)
}

func VerifikasiDokumenOutsideFaculty() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.VerifikasiDokumenOutsideFaculty",
		"Akses ditolak: prodi bukan di bawah fakultas penugasan anda",
	)
}

func DokumenAlreadyVerifiedByAdmin() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.DokumenAlreadyVerifiedByAdmin",
		"Dokumen sudah diverifikasi penuh oleh admin",
	)
}

func DokumenProdiRequiresFacultyFirst() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.DokumenProdiRequiresFacultyFirst",
		"Dokumen prodi harus diverifikasi oleh fakultas terlebih dahulu sebelum verifikasi admin",
	)
}

func DokumenProdiRejectRequiresFacultyFirst() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.DokumenProdiRejectRequiresFacultyFirst",
		"Dokumen prodi harus diverifikasi oleh fakultas terlebih dahulu sebelum ditolak admin",
	)
}

func DokumenNonProdiAdminOnly() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.DokumenNonProdiAdminOnly",
		"Dokumen unit selain prodi hanya dapat diverifikasi oleh admin",
	)
}

func JadwalNotFound() commondomain.Error {
	return commondomain.NotFoundError(
		"ProkerV3.JadwalNotFound",
		"Jadwal proker tidak ditemukan",
	)
}

func TahunRequired() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.TahunRequired",
		"Tahun wajib diisi",
	)
}

func TanggalMulaiRequired() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.TanggalMulaiRequired",
		"Tanggal mulai wajib diisi",
	)
}

func TanggalAkhirRequired() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.TanggalAkhirRequired",
		"Tanggal akhir wajib diisi",
	)
}

func TanggalAkhirInvalid() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.TanggalAkhirInvalid",
		"Tanggal akhir tidak boleh lebih awal dari tanggal mulai",
	)
}

func PaguIdRequired() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.PaguIdRequired",
		"ID wajib diisi",
	)
}

func AdminOnlyProkerBatal() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.AdminOnlyProkerBatal",
		"Hanya Admin yang dapat mengubah tanggal aktif Proker Batal",
	)
}

func TanggalAktifTargetRequired() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.TanggalAktifTargetRequired",
		"Tanggal aktif harus diisi untuk target fakultas/unit yang dipilih",
	)
}

func DateFormatInvalid() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.DateFormatInvalid",
		"Format tanggal aktif Proker Batal harus YYYY-MM-DD",
	)
}

func RepoNotInitialized() commondomain.Error {
	return commondomain.FailureError(
		"ProkerV3.RepoNotInitialized",
		"repository IProkerV3Repository is not initialized",
	)
}
