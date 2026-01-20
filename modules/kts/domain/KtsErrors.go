package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("Kts.EmptyData", "data is not found")
}

func ExistData() domain.Error {
	return domain.NotFoundError("Kts.ExistData", "data is exist")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("Kts.InvalidUuid", "uuid is invalid")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("Kts.InvalidData", "data is invalid")
}

func InvalidTahun() domain.Error {
	return domain.NotFoundError("Kts.InvalidTahun", "tahun is invalid")
}

func InvalidAuditor() domain.Error {
	return domain.NotFoundError("Kts.InvalidAuditor", "auditor is invalid")
}

func InvalidAuditee() domain.Error {
	return domain.NotFoundError("Kts.InvalidAuditee", "auditee is invalid")
}

func InvalidStatusAcc() domain.Error {
	return domain.NotFoundError("Kts.InvalidStatusAcc", "auditor is invalid")
}

func InvalidTanggal() domain.Error {
	return domain.NotFoundError("Kts.InvalidTanggal", "tanggal is invalid")
}

func InvalidStep() domain.Error {
	return domain.NotFoundError("Kts.InvalidStep", "tanggal is invalid")
}

func RequiredNomorLaporan() domain.Error {
	return domain.NotFoundError("Kts.RequiredNomorLaporan", "nomor laporan is required")
}

func RequiredKeteranganTolak() domain.Error {
	return domain.NotFoundError("Kts.RequiredKeteranganTolak", "keterangan tolak is required")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("Kts.NotFound", fmt.Sprintf("Kts with identifier %s not found", id))
}

func NotFoundUser() domain.Error {
	return domain.NotFoundError("Kts.NotFoundUser", "user not found")
}
