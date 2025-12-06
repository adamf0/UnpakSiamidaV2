package infrastructure

import (
    domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"
    create "UnpakSiamida/modules/indikatorrenstra/application/CreateIndikatorRenstra"
    update "UnpakSiamida/modules/indikatorrenstra/application/UpdateIndikatorRenstra"
    delete "UnpakSiamida/modules/indikatorrenstra/application/DeleteIndikatorRenstra"
    get "UnpakSiamida/modules/indikatorrenstra/application/GetIndikatorRenstra"
    getAll "UnpakSiamida/modules/indikatorrenstra/application/GetAllIndikatorRenstras"
    "github.com/mehdihadeli/go-mediatr"
    "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func RegisterModuleIndikatorRenstra() {
    dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

    repoIndikatorRenstra := NewIndikatorRenstraRepository(db)
	// if err := db.AutoMigrate(&domainindikatorrenstra.IndikatorRenstra{}); err != nil {
	// 	panic(err)
	// }

    // Pipeline behavior
    // mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorIndikatorRenstra())

    // Register request handler
    mediatr.RegisterRequestHandler[
        create.CreateIndikatorRenstraCommand,
        string,
    ](&create.CreateIndikatorRenstraCommandHandler{
        Repo: repoIndikatorRenstra,
    })

    mediatr.RegisterRequestHandler[
        update.UpdateIndikatorRenstraCommand,
        string,
    ](&update.UpdateIndikatorRenstraCommandHandler{
        Repo: repoIndikatorRenstra,
    })

    mediatr.RegisterRequestHandler[
        delete.DeleteIndikatorRenstraCommand,
        string,
    ](&delete.DeleteIndikatorRenstraCommandHandler{
        Repo: repoIndikatorRenstra,
    })

    mediatr.RegisterRequestHandler[
        get.GetIndikatorRenstraByUuidQuery,
        *domainindikatorrenstra.IndikatorRenstra,
    ](&get.GetIndikatorRenstraByUuidQueryHandler{
        Repo: repoIndikatorRenstra,
    })

    mediatr.RegisterRequestHandler[
        getAll.GetAllIndikatorRenstrasQuery,
        domainindikatorrenstra.PagedIndikatorRenstras,
    ](&getAll.GetAllIndikatorRenstrasQueryHandler{
        Repo: repoIndikatorRenstra,
    })


}
