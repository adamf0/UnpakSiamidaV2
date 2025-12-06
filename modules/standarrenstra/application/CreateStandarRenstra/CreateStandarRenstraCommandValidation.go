package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func CreateStandarRenstraCommandValidation(cmd CreateStandarRenstraCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Nama,
			validation.Required.Error("Nama cannot be blank"),
		),
	)
}