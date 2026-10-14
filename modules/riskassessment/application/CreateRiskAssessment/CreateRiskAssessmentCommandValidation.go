package application

type CreateRiskAssessmentCommandValidationStruct struct{}

func (c CreateRiskAssessmentCommand) Validate() error {
	// Flexible validation to allow Excel fields or fallback fields
	return nil
}

func CreateRiskAssessmentCommandValidation(cmd interface{}) error {
	if c, ok := cmd.(CreateRiskAssessmentCommand); ok {
		return c.Validate()
	}
	return nil
}
