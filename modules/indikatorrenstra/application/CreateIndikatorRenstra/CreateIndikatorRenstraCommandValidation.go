package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	helper "UnpakSiamida/common/helper"
)

func CreateIndikatorRenstraCommandValidation(cmd CreateIndikatorRenstraCommand) error {
	return validation.ValidateStruct(&cmd,

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
