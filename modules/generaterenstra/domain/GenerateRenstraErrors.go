package domain

import (
	"UnpakSiamida/common/domain"
	"fmt"
)

func EmptyData() domain.Error {
	return domain.NotFoundError("GenerateRenstra.EmptyData", "data is not found")
}

func InvalidUuid() domain.Error {
	return domain.NotFoundError("GenerateRenstra.InvalidUuid", "uuid is invalid")
}

func InvalidData() domain.Error {
	return domain.NotFoundError("GenerateRenstra.InvalidData", "data is invalid")
}

func NotFound(id string) domain.Error {
	return domain.NotFoundError("GenerateRenstra.NotFound", fmt.Sprintf("GenerateRenstra with identifier %s not found", id) )
}
func NotFoundFakultasUnit(id string) domain.Error {
	return domain.NotFoundError("GenerateRenstra.NotFoundFakultasUnit", fmt.Sprintf("Fakultas unit with identifier %s not found", id) )
}
func NotFoundRenstra(id string) domain.Error {
	return domain.NotFoundError("GenerateRenstra.NotFoundRenstra", fmt.Sprintf("renstra with identifier %s not found", id) )
}
func NotFoundTemplate(tahun string, fakultas string) domain.Error {
	return domain.NotFoundError("GenerateRenstra.NotFoundTemplate", fmt.Sprintf("template not found with identifier tahun %s & fakultas %s", tahun, fakultas) )
}
func NotFoundAudit(tahun string, fakultas string) domain.Error {
	return domain.NotFoundError("GenerateRenstra.NotFoundAudit", fmt.Sprintf("previous audit not found with identifier tahun %s & fakultas %s", tahun, fakultas) )
}

func InvalidFakultasUnit() domain.Error {
	return domain.NotFoundError("GenerateRenstra.InvalidFakultasUnit", "fakultas unit is invalid")
}

func InvalidTahunRenstra(
	templateUuidExisting string,
	indikatorExisting string,
	tahunTemplateExisting string,
	tahun string,
	operation string,
) domain.Error {

	action := "rollback"
	if operation == "insert" {
		action = "generated"
	}

	nextSugestion := "Please manually delete the problematic audit questions, then you can regenerate them."
	if operation == "insert" {
		nextSugestion = "Please check the year of the indicator question again, is it correct?"
	}

	return domain.NotFoundError(
		"GenerateRenstra.InvalidTahunRenstra",
		fmt.Sprintf(
			"template %s with indicator '%s (%s)' cannot be %s because the indicator year does not match the audit year (%s). %s",
			templateUuidExisting,
			indikatorExisting,
			tahunTemplateExisting,
			action,
			tahun,
			nextSugestion,
		),
	)
}

func InvalidTahunDokumenTambahan(
	templateUuidExisting string,
	jenisFileExisting string,
	tahunTemplateExisting string,
	tahun string,
	operation string,
) domain.Error {

	action := "rollback"
	if operation == "insert" {
		action = "generated"
	}

	nextSugestion := "Please manually delete the problematic audit questions, then you can regenerate them."
	if operation == "insert" {
		nextSugestion = "Please check the year of the jenis file question again, is it correct?"
	}

	return domain.NotFoundError(
		"GenerateRenstra.InvalidTahunDokumenTambahan",
		fmt.Sprintf(
			"template %s with jenis file '%s (%s)' cannot be %s because the jenis file year does not match the audit year (%s). %s",
			templateUuidExisting,
			jenisFileExisting,
			tahunTemplateExisting,
			action,
			tahun,
			nextSugestion,
		),
	)
}

func InvalidTemplate() domain.Error {
	return domain.NotFoundError("GenerateRenstra.InvalidTemplate", "template is invalid")
}

func InvalidRenstra() domain.Error {
	return domain.NotFoundError("GenerateRenstra.InvalidRenstra", "renstra is invalid")
}

func InvalidType(tipe string) domain.Error {
	return domain.NotFoundError("GenerateRenstra.InvalidType", fmt.Sprintf(
			"type %s is invalid",
			tipe,
	))
}

func InvalidTugas() domain.Error {
	return domain.NotFoundError("GenerateRenstra.InvalidTugas", "tugas is invalid")
}

func InvalidParsing(target string) domain.Error {
	return domain.NotFoundError("GenerateRenstra.IvalidParsing", fmt.Sprintf("failed parsing %s to UUID", target) )
}