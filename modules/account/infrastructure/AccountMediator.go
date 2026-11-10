package infrastructure

import (
	commoninfra "UnpakSiamida/common/infrastructure"
	login "UnpakSiamida/modules/account/application/Login"
	updateprofile "UnpakSiamida/modules/account/application/UpdateProfile"
	who "UnpakSiamida/modules/account/application/Whoami"
	domain "UnpakSiamida/modules/account/domain"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterModuleAccount(db *gorm.DB) error {
	repoAccount := NewAccountRepository(db)

	mediatr.RegisterRequestHandler[
		who.WhoamiCommand,
		*domain.Account,
	](&who.WhoamiCommandHandler{
		Repo: repoAccount,
	})

	mediatr.RegisterRequestHandler[
		login.LoginCommand,
		*domain.LoginResult,
	](&login.LoginCommandHandler{
		Repo: repoAccount,
	})

	mediatr.RegisterRequestHandler[
		updateprofile.UpdateProfileCommand,
		*domain.Account,
	](&updateprofile.UpdateProfileCommandHandler{
		Repo: repoAccount,
	})

	commoninfra.RegisterValidation(login.LoginCommandValidation, "Login.Validation")
	commoninfra.RegisterValidation(updateprofile.UpdateProfileCommandValidation, "Account.UpdateProfile.Validation")

	return nil
}
