package infrastructure

import (
    domaindokumentambahan "UnpakSiamida/modules/dokumentambahan/domain"
    update "UnpakSiamida/modules/dokumentambahan/application/UpdateDokumenTambahan"
    delete "UnpakSiamida/modules/dokumentambahan/application/DeleteDokumenTambahan"
    get "UnpakSiamida/modules/dokumentambahan/application/GetDokumenTambahan"
    getAll "UnpakSiamida/modules/dokumentambahan/application/GetAllDokumenTambahans"
    setupUuid "UnpakSiamida/modules/dokumentambahan/application/SetupUuidDokumenTambahan"

    infraRenstra "UnpakSiamida/modules/renstra/infrastructure"
    "github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterModuleDokumenTambahan(db *gorm.DB) error{
    repoDokumenTambahan := NewDokumenTambahanRepository(db)
    repoRenstra := infraRenstra.NewRenstraRepository(db)

    mediatr.RegisterRequestHandler[
        update.UpdateDokumenTambahanCommand,
        string,
    ](&update.UpdateDokumenTambahanCommandHandler{
        Repo: repoDokumenTambahan,
        RepoRenstra: repoRenstra,
    })

    mediatr.RegisterRequestHandler[
        delete.DeleteDokumenTambahanCommand,
        string,
    ](&delete.DeleteDokumenTambahanCommandHandler{
        Repo: repoDokumenTambahan,
    })

    mediatr.RegisterRequestHandler[
        get.GetDokumenTambahanByUuidQuery,
        *domaindokumentambahan.DokumenTambahan,
    ](&get.GetDokumenTambahanByUuidQueryHandler{
        Repo: repoDokumenTambahan,
    })

    mediatr.RegisterRequestHandler[
        getAll.GetAllDokumenTambahansQuery,
        domaindokumentambahan.PagedDokumenTambahans,
    ](&getAll.GetAllDokumenTambahansQueryHandler{
        Repo: repoDokumenTambahan,
    })

    mediatr.RegisterRequestHandler[
        setupUuid.SetupUuidDokumenTambahanCommand,
        string,
    ](&setupUuid.SetupUuidDokumenTambahanCommandHandler{
        Repo: repoDokumenTambahan,
    })

    return nil
}
