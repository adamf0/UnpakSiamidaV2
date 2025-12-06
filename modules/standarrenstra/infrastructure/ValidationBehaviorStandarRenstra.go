package infrastructure

import (
	"context"
	"strings"

	mediatr "github.com/mehdihadeli/go-mediatr"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	createStandarRenstra "UnpakSiamida/modules/standarrenstra/application/CreateStandarRenstra"
	updateStandarRenstra "UnpakSiamida/modules/standarrenstra/application/UpdateStandarRenstra"
	deleteStandarRenstra "UnpakSiamida/modules/standarrenstra/application/DeleteStandarRenstra"

	commoninfra "UnpakSiamida/common/infrastructure"
	domain "UnpakSiamida/common/domain"
)

type ValidationBehaviorStandarRenstra struct{}

func NewValidationBehaviorStandarRenstra() *ValidationBehaviorStandarRenstra {
	return &ValidationBehaviorStandarRenstra{}
}

func (b *ValidationBehaviorStandarRenstra) Handle(
	ctx context.Context,
	request interface{},
	next mediatr.RequestHandlerFunc,
) (interface{}, error) {

	// -------------------------
	//  1. VALIDASI (OZZO)
	// -------------------------
	switch cmd := request.(type) {

	case createStandarRenstra.CreateStandarRenstraCommand:

		if err := createStandarRenstra.CreateStandarRenstraCommandValidation(cmd); err != nil {

			// jika error berasal dari ozzo validation (validation.Errors -> map[string]error)
			if ve, ok := err.(validation.Errors); ok {
				msgs := make(map[string]string)
				for field, ferr := range ve {
					// ubah nama field jadi lower-case (sesuai keinginan)
					key := strings.ToLower(field)
					msgs[key] = ferr.Error()
				}
				return nil, commoninfra.NewResponseError(
					"StandarRenstra.Validation",
					msgs,
				)
			}

			// fallback: bukan validation.Errors, kirim string message biasa
			return nil, commoninfra.NewResponseError(
				"StandarRenstra.Validation",
				err.Error(),
			)
		}

	case updateStandarRenstra.UpdateStandarRenstraCommand:

		if err := updateStandarRenstra.UpdateStandarRenstraCommandValidation(cmd); err != nil {

			// jika error berasal dari ozzo validation (validation.Errors -> map[string]error)
			if ve, ok := err.(validation.Errors); ok {
				msgs := make(map[string]string)
				for field, ferr := range ve {
					// ubah nama field jadi lower-case (sesuai keinginan)
					key := strings.ToLower(field)
					msgs[key] = ferr.Error()
				}
				return nil, commoninfra.NewResponseError(
					"StandarRenstra.Validation",
					msgs,
				)
			}

			// fallback: bukan validation.Errors, kirim string message biasa
			return nil, commoninfra.NewResponseError(
				"StandarRenstra.Validation",
				err.Error(),
			)
		}

	case deleteStandarRenstra.DeleteStandarRenstraCommand:

		if err := deleteStandarRenstra.DeleteStandarRenstraCommandValidation(cmd); err != nil {

			// jika error berasal dari ozzo validation (validation.Errors -> map[string]error)
			if ve, ok := err.(validation.Errors); ok {
				msgs := make(map[string]string)
				for field, ferr := range ve {
					// ubah nama field jadi lower-case (sesuai keinginan)
					key := strings.ToLower(field)
					msgs[key] = ferr.Error()
				}
				return nil, commoninfra.NewResponseError(
					"StandarRenstra.Validation",
					msgs,
				)
			}

			// fallback: bukan validation.Errors, kirim string message biasa
			return nil, commoninfra.NewResponseError(
				"StandarRenstra.Validation",
				err.Error(),
			)
		}
	}

	// -------------------------
	//  2. LANJUT KE HANDLER
	// -------------------------
	result, err := next(ctx)

	// -------------------------
	//  3. ERROR DOMAIN (domain.Error)
	// -------------------------
	if derr, ok := err.(domain.Error); ok {
		return nil, commoninfra.NewResponseError(
			derr.Code,
			derr.Description,
		)
	}

	// -------------------------
	//  4. ERROR TAK TERDUGA
	// -------------------------
	if err != nil {
		return nil, commoninfra.NewInternalError(err)
	}

	return result, nil
}
