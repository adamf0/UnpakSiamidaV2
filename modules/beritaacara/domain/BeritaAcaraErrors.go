package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("BeritaAcara.EmptyData", "data is not found")
}

func InvalidTanggal() domain.Error {
	return domain.NotFoundError("BeritaAcara.InvalidTanggal", "tanggal is invalid")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("BeritaAcara.InvalidUuid", "uuid is invalid")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("BeritaAcara.InvalidData", "data is invalid")
}

func DuplicateAssigment() domain.Error {
	return domain.NotFoundError("BeritaAcara.DuplicateAssigment", "auditee, auditee 1, and auditor 2 must not have the same target")
}

func InvalidFakultasUnit() domain.Error {
	return domain.NotFoundError("BeritaAcara.InvalidFakultasUnit", "fakultas is invalid")
}

func InvalidAuditee() domain.Error {
	return domain.NotFoundError("BeritaAcara.InvalidAuditee", "auditee is invalid")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("BeritaAcara.NotFound", fmt.Sprintf("BeritaAcara with identifier %s not found", id))
}

func NotFoundFakultas() domain.Error {
	return domain.NotFoundError("BeritaAcara.NotFoundFakultas", "fakultas not found")
}

func NotFoundAuditee() domain.Error {
	return domain.NotFoundError("BeritaAcara.NotFoundAuditee", "auditee not found")
}

func NotFoundAuditor() domain.Error {
	return domain.NotFoundError("BeritaAcara.NotFoundAuditor", "auditor not found")
}

func NotGranted() domain.Error {
	return domain.NotFoundError("BeritaAcara.NotGranted", "Download berita acara is rejected because not granted")
}
func NotPushDownload() domain.Error {
	return domain.NotFoundError("BeritaAcara.NotPushDownload", "Download berita acara is not available")
}
func NoPermission() domain.Error {
	return domain.NotFoundError("BeritaAcara.NoPermission", "Something wrong to access resource")
}
func NoResource() domain.Error {
	return domain.NotFoundError("BeritaAcara.NoResource", "resource/assets not found to build pdf")
}
func GeneratePDF(message string) domain.Error {
	return domain.NotFoundError("BeritaAcara.GeneratePDF", fmt.Sprintf("error: %s", message))
}
