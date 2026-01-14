package infrastructure

import (
    login "UnpakSiamida/modules/account/application/Login"
    who "UnpakSiamida/modules/account/application/Whoami"
    "github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
    domain "UnpakSiamida/modules/account/domain"
)

func RegisterModuleAccount(db *gorm.DB) error{
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

    return nil
}
