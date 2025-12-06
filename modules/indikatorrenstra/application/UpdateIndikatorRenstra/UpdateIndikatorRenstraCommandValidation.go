package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	helper "UnpakSiamida/common/helper"
)

func UpdateIndikatorRenstraCommandValidation(cmd UpdateIndikatorRenstraCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Uuid,
			validation.Required.Error("UUID cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
		),

		validation.Field(&cmd.StandarRenstra,
			validation.Required.Error("Standar Renstra wajib diisi"),
		),

		validation.Field(&cmd.Indikator,
			validation.Required.Error("Indikator wajib diisi"),
		),

		validation.Field(&cmd.Parent,
			validation.By(helper.ValidateParent),
		),

		validation.Field(&cmd.Tahun,
			validation.Required.Error("Tahun wajib diisi"),
		),

		validation.Field(&cmd.TipeTarget,
			validation.Required.Error("Tipe Target wajib diisi"),
		),
	)
}