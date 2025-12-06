package infrastructure

import (
	"context"
	"strings"

	mediatr "github.com/mehdihadeli/go-mediatr"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	createIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/CreateIndikatorRenstra"
	updateIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/UpdateIndikatorRenstra"
	deleteIndikatorRenstra "UnpakSiamida/modules/indikatorrenstra/application/DeleteIndikatorRenstra"

	commoninfra "UnpakSiamida/common/infrastructure"
	domain "UnpakSiamida/common/domain"
)

type ValidationBehaviorIndikatorRenstra struct{}

func NewValidationBehaviorIndikatorRenstra() *ValidationBehaviorIndikatorRenstra {
	return &ValidationBehaviorIndikatorRenstra{}
}

func (b *ValidationBehaviorIndikatorRenstra) Handle(
	ctx context.Context,
	request interface{},
	next mediatr.RequestHandlerFunc,
) (interface{}, error) {

	// -------------------------
	//  1. VALIDASI (OZZO)
	// -------------------------
	switch cmd := request.(type) {

	case createIndikatorRenstra.CreateIndikatorRenstraCommand:

		if err := createIndikatorRenstra.CreateIndikatorRenstraCommandValidation(cmd); err != nil {

			// jika error berasal dari ozzo validation (validation.Errors -> map[string]error)
			if ve, ok := err.(validation.Errors); ok {
				msgs := make(map[string]string)
				for field, ferr := range ve {
					// ubah nama field jadi lower-case (sesuai keinginan)
					key := strings.ToLower(field)
					msgs[key] = ferr.Error()
				}
				return nil, commoninfra.NewResponseError(
					"IndikatorRenstra.Validation",
					msgs,
				)
			}

			// fallback: bukan validation.Errors, kirim string message biasa
			return nil, commoninfra.NewResponseError(
				"IndikatorRenstra.Validation",
				err.Error(),
			)
		}

	case updateIndikatorRenstra.UpdateIndikatorRenstraCommand:

		if err := updateIndikatorRenstra.UpdateIndikatorRenstraCommandValidation(cmd); err != nil {

			// jika error berasal dari ozzo validation (validation.Errors -> map[string]error)
			if ve, ok := err.(validation.Errors); ok {
				msgs := make(map[string]string)
				for field, ferr := range ve {
					// ubah nama field jadi lower-case (sesuai keinginan)
					key := strings.ToLower(field)
					msgs[key] = ferr.Error()
				}
				return nil, commoninfra.NewResponseError(
					"IndikatorRenstra.Validation",
					msgs,
				)
			}

			// fallback: bukan validation.Errors, kirim string message biasa
			return nil, commoninfra.NewResponseError(
				"IndikatorRenstra.Validation",
				err.Error(),
			)
		}

	case deleteIndikatorRenstra.DeleteIndikatorRenstraCommand:

		if err := deleteIndikatorRenstra.DeleteIndikatorRenstraCommandValidation(cmd); err != nil {

			// jika error berasal dari ozzo validation (validation.Errors -> map[string]error)
			if ve, ok := err.(validation.Errors); ok {
				msgs := make(map[string]string)
				for field, ferr := range ve {
					// ubah nama field jadi lower-case (sesuai keinginan)
					key := strings.ToLower(field)
					msgs[key] = ferr.Error()
				}
				return nil, commoninfra.NewResponseError(
					"IndikatorRenstra.Validation",
					msgs,
				)
			}

			// fallback: bukan validation.Errors, kirim string message biasa
			return nil, commoninfra.NewResponseError(
				"IndikatorRenstra.Validation",
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
