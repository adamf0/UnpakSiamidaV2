package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func SendEmailRenstraCommandValidation(cmd SendEmailRenstraCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Uuids,
			validation.Required.Error("Pilih minimal 1 jadwal audit untuk dikirim email kode akses."),
		),
	)
}
