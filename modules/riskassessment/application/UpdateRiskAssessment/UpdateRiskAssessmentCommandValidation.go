package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (c UpdateRiskAssessmentCommand) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.UUID, validation.Required.Error("UUID analisis risiko wajib diisi")),
	)
}

func UpdateRiskAssessmentCommandValidation(cmd interface{}) error {
	if c, ok := cmd.(UpdateRiskAssessmentCommand); ok {
		return c.Validate()
	}
	return nil
}
