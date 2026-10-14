package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (c DeleteSwotCommand) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.UUID, validation.Required.Error("UUID SWOT wajib diisi")),
	)
}

func DeleteSwotCommandValidation(cmd interface{}) error {
	if c, ok := cmd.(DeleteSwotCommand); ok {
		return c.Validate()
	}
	return nil
}
