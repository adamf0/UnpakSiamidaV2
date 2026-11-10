package application

import (
	"fmt"
	"strings"

	helper "UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func UpdateProfileCommandValidation(cmd UpdateProfileCommand) error {
	hasNewPassword := strings.TrimSpace(cmd.NewPassword) != ""

	return validation.ValidateStruct(&cmd,
		validation.Field(&cmd.SID,
			validation.Required.Error("UUID pengguna tidak valid"),
			validation.By(helper.ValidateUUIDv4),
		),
		validation.Field(&cmd.Username,
			validation.Required.Error("Username wajib diisi"),
			validation.By(func(value interface{}) error {
				s, _ := value.(string)
				trimmed := strings.TrimSpace(s)
				if len(trimmed) < 3 {
					return fmt.Errorf("Username minimal 3 karakter")
				}
				if len(trimmed) > 100 {
					return fmt.Errorf("Username maksimal 100 karakter")
				}
				return nil
			}),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.Name,
			validation.Required.Error("Nama lengkap wajib diisi"),
			validation.By(func(value interface{}) error {
				s, _ := value.(string)
				trimmed := strings.TrimSpace(s)
				if trimmed == "" {
					return fmt.Errorf("Nama lengkap wajib diisi")
				}
				if len(trimmed) > 255 {
					return fmt.Errorf("Nama lengkap maksimal 255 karakter")
				}
				return nil
			}),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.Email,
			validation.By(helper.ValidateOptionalEmail),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.CurrentPassword,
			validation.When(hasNewPassword,
				validation.Required.Error("Password saat ini wajib diisi untuk mengubah password"),
			),
		),
		validation.Field(&cmd.NewPassword,
			validation.When(hasNewPassword,
				validation.By(func(value interface{}) error {
					s, _ := value.(string)
					trimmed := strings.TrimSpace(s)
					if len(trimmed) < 3 {
						return fmt.Errorf("Password baru minimal 3 karakter")
					}
					if len(trimmed) > 255 {
						return fmt.Errorf("Password baru maksimal 255 karakter")
					}
					return nil
				}),
				validation.By(helper.NoXSSFullScanWithDecode()),
			),
		),
		validation.Field(&cmd.ConfirmPassword,
			validation.When(hasNewPassword,
				validation.Required.Error("Konfirmasi password baru wajib diisi"),
				validation.By(func(value interface{}) error {
					s, _ := value.(string)
					if strings.TrimSpace(s) != strings.TrimSpace(cmd.NewPassword) {
						return fmt.Errorf("Konfirmasi password baru tidak cocok")
					}
					return nil
				}),
			),
		),
	)
}
