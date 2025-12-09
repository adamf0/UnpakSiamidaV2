package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	helper "UnpakSiamida/common/helper"
)

func UpdateTemplateDokumenTambahanCommandValidation(cmd UpdateTemplateDokumenTambahanCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Uuid,
			validation.Required.Error("UUID cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
		),
		validation.Field(&cmd.Tahun,
			validation.Required.Error("Tahun cannot be blank"),
		),
		validation.Field(&cmd.JenisFile,
			validation.Required.Error("JenisFile cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
		),
		validation.Field(&cmd.Pertanyaan,
			validation.Required.Error("Pertanyaan cannot be blank"),
		),
		validation.Field(&cmd.Klasifikasi,
			validation.Required.Error("Klasifikasi cannot be blank"),
			validation.In("minor", "major").Error("Klasifikasi invalid value"),
		),
		validation.Field(&cmd.Kategori, //fakultas_prodi_unit
			validation.Required.Error("Kategori cannot be blank"),
		),
		validation.Field(&cmd.Tugas,
			validation.Required.Error("Tugas cannot be blank"),
			validation.In("auditor1", "auditor2").Error("Tugas invalid value"),
		),
	)
}