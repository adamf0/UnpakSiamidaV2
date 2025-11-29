package domain

type ValidationError struct {
	Errors []Error
}

func NewValidationError(errors []Error) *ValidationError {
	return &ValidationError{Errors: errors}
}

func FromResults(results []Result) *ValidationError {
	var errs []Error
	for _, r := range results {
		if !r.IsSuccess {
			errs = append(errs, r.Error)
		}
	}
	return NewValidationError(errs)
}
