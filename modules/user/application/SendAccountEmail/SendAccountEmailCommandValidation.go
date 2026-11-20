package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func SendAccountEmailCommandValidation(cmd SendAccountEmailCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.UserUUIDs,
			validation.Required.Error("Setidaknya satu user_uuid harus dipilih"),
		),
	)
}
