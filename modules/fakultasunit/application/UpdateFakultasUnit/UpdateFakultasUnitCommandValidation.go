package application

import (
	helper "UnpakSiamida/common/helper"
	"errors"
	"strings"
	"unicode/utf8"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func UpdateFakultasUnitCommandValidation(cmd UpdateFakultasUnitCommand) error {
	errs := validation.Errors{}

	if cmd.Nama == nil && cmd.IsActive == nil {
		errs["payload"] = errors.New("nama or is_active must be provided")
	}

	if err := validation.Validate(cmd.Uuid,
		validation.Required.Error("UUID cannot be blank"),
		validation.By(helper.ValidateUUIDv4),
		validation.By(helper.NoXSSFullScanWithDecode()),
	); err != nil {
		errs["Uuid"] = err
	}

	if cmd.Nama != nil {
		trimmed := strings.TrimSpace(*cmd.Nama)
		if trimmed == "" {
			errs["Nama"] = errors.New("Nama cannot be blank")
		} else if utf8.RuneCountInString(trimmed) > 100 {
			errs["Nama"] = errors.New("Nama must be at most 100 characters")
		} else if err := helper.NoXSSFullScanWithDecode()(trimmed); err != nil {
			errs["Nama"] = err
		}
	}

	return errs.Filter()
}
