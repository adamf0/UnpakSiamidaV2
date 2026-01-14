package infrastructure

import (
    domainrenstranilai "UnpakSiamida/modules/renstranilai/domain"
    update "UnpakSiamida/modules/renstranilai/application/UpdateRenstraNilai"
    delete "UnpakSiamida/modules/renstranilai/application/DeleteRenstraNilai"
    get "UnpakSiamida/modules/renstranilai/application/GetRenstraNilai"
    getAll "UnpakSiamida/modules/renstranilai/application/GetAllRenstraNilais"
    setupUuid "UnpakSiamida/modules/renstranilai/application/SetupUuidRenstraNilai"

    infraRenstra "UnpakSiamida/modules/renstra/infrastructure"
    "github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterModuleRenstraNilai(db *gorm.DB) error{
    repoRenstraNilai := NewRenstraNilaiRepository(db)
    repoRenstra := infraRenstra.NewRenstraRepository(db)

    mediatr.RegisterRequestHandler[
        update.UpdateRenstraNilaiCommand,
        string,
    ](&update.UpdateRenstraNilaiCommandHandler{
        Repo: repoRenstraNilai,
        RepoRenstra: repoRenstra,
    })

    mediatr.RegisterRequestHandler[
        delete.DeleteRenstraNilaiCommand,
        string,
    ](&delete.DeleteRenstraNilaiCommandHandler{
        Repo: repoRenstraNilai,
    })

    mediatr.RegisterRequestHandler[
        get.GetRenstraNilaiByUuidQuery,
        *domainrenstranilai.RenstraNilai,
    ](&get.GetRenstraNilaiByUuidQueryHandler{
        Repo: repoRenstraNilai,
    })

    mediatr.RegisterRequestHandler[
        getAll.GetAllRenstraNilaisQuery,
        domainrenstranilai.PagedRenstraNilais,
    ](&getAll.GetAllRenstraNilaisQueryHandler{
        Repo: repoRenstraNilai,
    })

    mediatr.RegisterRequestHandler[
        setupUuid.SetupUuidRenstraNilaiCommand,
        string,
    ](&setupUuid.SetupUuidRenstraNilaiCommandHandler{
        Repo: repoRenstraNilai,
    })

    return nil
}
