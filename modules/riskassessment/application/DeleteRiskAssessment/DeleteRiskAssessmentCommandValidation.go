package application

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (c DeleteRiskAssessmentCommand) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.UUID, validation.Required.Error("UUID analisis risiko wajib diisi")),
	)
}

func DeleteRiskAssessmentCommandValidation(cmd interface{}) error {
	if c, ok := cmd.(DeleteRiskAssessmentCommand); ok {
		return c.Validate()
	}
	return nil
}
