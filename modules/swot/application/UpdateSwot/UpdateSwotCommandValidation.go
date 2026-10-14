package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (c UpdateSwotCommand) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.UUID, validation.Required.Error("UUID SWOT wajib diisi")),
		validation.Field(&c.Kategori, validation.Required.Error("Kategori SWOT wajib diisi")),
		validation.Field(&c.Deskripsi, validation.Required.Error("Deskripsi SWOT wajib diisi")),
	)
}

func UpdateSwotCommandValidation(cmd interface{}) error {
	if c, ok := cmd.(UpdateSwotCommand); ok {
		return c.Validate()
	}
	return nil
}
