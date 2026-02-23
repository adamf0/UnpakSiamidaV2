package application

import (
	"UnpakSiamida/common/helper"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func PublishKtsCommandValidation(cmd PublishKtsCommand) error {
	rules := []*validation.FieldRules{
		validation.Field(&cmd.Uuid,
			validation.Required.Error("Uuid cannot be blank"),
			validation.By(helper.ValidateUUIDv4),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
		validation.Field(&cmd.Token,
			validation.Required.Error("Token cannot be blank"),
			validation.By(helper.NoXSSFullScanWithDecode()),
		),
	}
	return validation.ValidateStruct(&cmd, rules...)
}
