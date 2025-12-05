package infrastructure

import (
	"context"
	"strings"

	mediatr "github.com/mehdihadeli/go-mediatr"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	createUser "UnpakSiamida/modules/user/application/CreateUser"
	updateUser "UnpakSiamida/modules/user/application/UpdateUser"
	deleteUser "UnpakSiamida/modules/user/application/DeleteUser"

	commoninfra "UnpakSiamida/common/infrastructure"
	domain "UnpakSiamida/common/domain"
)

type ValidationBehavior struct{}

func NewValidationBehavior() *ValidationBehavior {
	return &ValidationBehavior{}
}

func (b *ValidationBehavior) Handle(
	ctx context.Context,
	request interface{},
	next mediatr.RequestHandlerFunc,
) (interface{}, error) {

	// -------------------------
	//  1. VALIDASI (OZZO)
	// -------------------------
	switch cmd := request.(type) {

	case createUser.CreateUserCommand:

		if err := createUser.CreateUserCommandValidation(cmd); err != nil {

			// jika error berasal dari ozzo validation (validation.Errors -> map[string]error)
			if ve, ok := err.(validation.Errors); ok {
				msgs := make(map[string]string)
				for field, ferr := range ve {
					// ubah nama field jadi lower-case (sesuai keinginan)
					key := strings.ToLower(field)
					msgs[key] = ferr.Error()
				}
				return nil, commoninfra.NewResponseError(
					"User.Validation",
					msgs,
				)
			}

			// fallback: bukan validation.Errors, kirim string message biasa
			return nil, commoninfra.NewResponseError(
				"User.Validation",
				err.Error(),
			)
		}

	case updateUser.UpdateUserCommand:

		if err := updateUser.UpdateUserCommandValidation(cmd); err != nil {

			// jika error berasal dari ozzo validation (validation.Errors -> map[string]error)
			if ve, ok := err.(validation.Errors); ok {
				msgs := make(map[string]string)
				for field, ferr := range ve {
					// ubah nama field jadi lower-case (sesuai keinginan)
					key := strings.ToLower(field)
					msgs[key] = ferr.Error()
				}
				return nil, commoninfra.NewResponseError(
					"User.Validation",
					msgs,
				)
			}

			// fallback: bukan validation.Errors, kirim string message biasa
			return nil, commoninfra.NewResponseError(
				"User.Validation",
				err.Error(),
			)
		}

	case deleteUser.DeleteUserCommand:

		if err := deleteUser.DeleteUserCommandValidation(cmd); err != nil {

			// jika error berasal dari ozzo validation (validation.Errors -> map[string]error)
			if ve, ok := err.(validation.Errors); ok {
				msgs := make(map[string]string)
				for field, ferr := range ve {
					// ubah nama field jadi lower-case (sesuai keinginan)
					key := strings.ToLower(field)
					msgs[key] = ferr.Error()
				}
				return nil, commoninfra.NewResponseError(
					"User.Validation",
					msgs,
				)
			}

			// fallback: bukan validation.Errors, kirim string message biasa
			return nil, commoninfra.NewResponseError(
				"User.Validation",
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
