package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	helper "UnpakSiamida/common/helper"
)

func CreateTemplateRenstraCommandValidation(cmd CreateTemplateRenstraCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Tahun,
			validation.Required.Error("Tahun cannot be blank"),
		),
		validation.Field(&cmd.Indikator,
			validation.Required.Error("Indikator cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
		),
		validation.Field(&cmd.IsPertanyaan,
			validation.Required.Error("IsPertanyaan cannot be blank"),
			validation.In("1", "0").Error("IsPertanyaan invalid value"),
		),
		validation.Field(&cmd.FakultasUnit,
			validation.Required.Error("FakultasUnit cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
		),
		validation.Field(&cmd.Kategori,
			validation.Required.Error("Kategori cannot be blank"),
		),
		validation.Field(&cmd.Klasifikasi,
			validation.Required.Error("Klasifikasi cannot be blank"),
			validation.In("minor", "major").Error("Klasifikasi invalid value"),
		),
		// validation.Field(&cmd.Target,
		// 	validation.When(cmd.Target != nil, validation.Nil), // skip TargetMin/Max jika Target ada
		// ),
		// validation.Field(&cmd.TargetMin,
		// 	validation.When(cmd.Target == nil, validation.Required.Error("TargetMin cannot be blank")),
		// ),
		// validation.Field(&cmd.TargetMax,
		// 	validation.When(cmd.Target == nil, validation.Required.Error("TargetMin cannot be blank")),
		// ),
		validation.Field(&cmd.Tugas,
			validation.Required.Error("Tugas cannot be blank"),
			validation.In("auditor1", "auditor2").Error("Tugas invalid value"),
		),
	)
}