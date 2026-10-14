package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (c CreateSwotCommand) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Kategori, validation.Required.Error("Kategori SWOT wajib diisi")),
		validation.Field(&c.Deskripsi, validation.Required.Error("Deskripsi SWOT wajib diisi")),
	)
}

func CreateSwotCommandValidation(cmd interface{}) error {
	if c, ok := cmd.(CreateSwotCommand); ok {
		return c.Validate()
	}
	return nil
}
