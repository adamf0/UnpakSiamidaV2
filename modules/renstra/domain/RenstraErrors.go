package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("Renstra.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("Renstra.InvalidUuid", "uuid is invalid")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("Renstra.InvalidData", "data is invalid")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("Renstra.NotFound", fmt.Sprintf("Renstra with identifier %s not found", id) )
}

func InvalidFakultasUnit() domain.Error {
	return domain.NotFoundError("Renstra.InvalidFakultasUnit", "fakultas unit is invalid")
}

func DataExisting() domain.Error {
	return domain.NotFoundError("Renstra.DataExisting", "the audit data form already existed before")
}

func MissingAuditee() domain.Error {
	return domain.NotFoundError("Renstra.MissingAuditee", "auditee have not been assigned")
}

func MissingAuditor1() domain.Error {
	return domain.NotFoundError("Renstra.MissingAuditor1", "auditor1 have not been assigned")
}

func MissingAuditor2() domain.Error {
	return domain.NotFoundError("Renstra.MissingAuditor2", "auditor2 have not been assigned")
}

func InvalidParsing(target string) domain.Error {
	return domain.NotFoundError("Renstra.IvalidParsing", fmt.Sprintf("failed parsing %s to UUID", target) )
}

func DuplicateAssigment() domain.Error {
	return domain.NotFoundError("Renstra.DuplicateAssigment", "auditee, auditee 1, and auditor 2 must not have the same target")
}

func InvalidDate(target string) domain.Error {
	return domain.NotFoundError("Renstra.InvalidDate", fmt.Sprintf("%s period have wrong date format", target) )
}

func PeriodOverlapUploadDokumen() domain.Error {
	return domain.NotFoundError("Renstra.PeriodOverlapUploadDokumen", "upload period overlaps with document period")
}

func PeriodOverlapUploadLapangan() domain.Error {
	return domain.NotFoundError("Renstra.PeriodOverlapUploadLapangan", "upload period overlaps with AL period")
}

func PeriodOverlapDokumenLapangan() domain.Error {
	return domain.NotFoundError("Renstra.PeriodOverlapDokumenLapangan", "document period overlaps with AL period")
}