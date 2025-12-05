package application

import (
	"fmt"
	"strings"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	helper "UnpakSiamida/common/helper"
)

func UpdateUserCommandValidation(cmd UpdateUserCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Uuid,
			validation.Required.Error("UUID cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
		),

		validation.Field(&cmd.Username,
			validation.Required.Error("Username cannot be blank"),
		),

		// validation.Field(&cmd.Password,
		// 	validation.Required.Error("Password cannot be blank"),
		// ),

		validation.Field(&cmd.Name,
			validation.Required.Error("Name cannot be blank"),
		),

		validation.Field(&cmd.Email,
			validation.Required.Error("Email cannot be blank"),
			validation.By(helper.ValidateUnpakEmail),
		),

		validation.Field(&cmd.FakultasUnit,
			validation.By(validateFakultasUnit),
		),
	)
}

func validateFakultasUnit(value interface{}) error {
	if value == nil {
		return nil // fakultas boleh null
	}

	ptr, ok := value.(*string)
	if !ok {
		return fmt.Errorf("FakultasUnit invalid type")
	}

	if ptr == nil {
		return nil
	}

	s := strings.TrimSpace(*ptr)
	if s == "" {
		return fmt.Errorf("FakultasUnit cannot be blank")
	}

	val, err := parseInt64(s)
	if err != nil {
		return err
	}

	if val < 1 {
		return fmt.Errorf("FakultasUnit invalid value")
	}

	return nil
}

func parseInt64(s string) (int64, error) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		if numErr, ok := err.(*strconv.NumError); ok {
			switch numErr.Err {
			case strconv.ErrRange:
				return 0, fmt.Errorf("Number out of range")
			case strconv.ErrSyntax:
				return 0, fmt.Errorf("Must be a number")
			}
		}
		return 0, fmt.Errorf("Invalid number")
	}
	return val, nil
}