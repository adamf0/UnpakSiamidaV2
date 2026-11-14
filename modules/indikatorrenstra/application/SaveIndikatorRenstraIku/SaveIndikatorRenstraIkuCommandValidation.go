package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func SaveIndikatorRenstraIkuCommandValidation(cmd SaveIndikatorRenstraIkuCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Identifier,
			validation.Required.Error("Identifier cannot be blank"),
		),
	)
}
