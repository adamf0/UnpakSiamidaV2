package application

import (
	"fmt"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func CreateUserCommandValidation(cmd CreateUserCommand) error {
	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.Nama, validation.Required.Error("cannot be blank")),
		validation.Field(&cmd.Skor, validation.By(func(value interface{}) error {
			if cmd.Nama == "admin" { // skip validation
				return nil
			}
			return validateSkor(value)
		})),
	)
}

// validateSkor adalah wrapper utama
func validateSkor(value interface{}) error {
	s, err := toString(value)
	if err != nil {
		return err
	}

	if err := checkEmpty(s); err != nil {
		return err
	}

	skor, err := parseInt64(s)
	if err != nil {
		return err
	}

	return checkMin(skor)
}

// konversi interface{} -> string
func toString(value interface{}) (string, error) {
	s, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("invalid type")
	}
	return strings.TrimSpace(s), nil
}

// cek kosong
func checkEmpty(s string) error {
	if s == "" {
		return fmt.Errorf("cannot be blank")
	}
	return nil
}

// parse string -> int64 dengan error handling
func parseInt64(s string) (int64, error) {
	skor, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		if numErr, ok := err.(*strconv.NumError); ok {
			switch numErr.Err {
			case strconv.ErrRange:
				return 0, fmt.Errorf("number out of range")
			case strconv.ErrSyntax:
				return 0, fmt.Errorf("must be a number")
			default:
				return 0, fmt.Errorf("invalid number")
			}
		}
		return 0, fmt.Errorf("invalid number")
	}
	return skor, nil
}

// cek minimum >= 0
func checkMin(skor int64) error {
	if skor < 0 {
		return fmt.Errorf("must be no less than 0")
	}
	return nil
}
