package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	create "UnpakSiamida/modules/indikatorrenstra/application/CreateIndikatorRenstra"
	delete "UnpakSiamida/modules/indikatorrenstra/application/DeleteIndikatorRenstra"
	getAll "UnpakSiamida/modules/indikatorrenstra/application/GetAllIndikatorRenstras"
	getAllMasterIku "UnpakSiamida/modules/indikatorrenstra/application/GetAllMasterIku"
	getDashboard "UnpakSiamida/modules/indikatorrenstra/application/GetDashboardIkuAnggaran"
	get "UnpakSiamida/modules/indikatorrenstra/application/GetIndikatorRenstra"
	getIndikatorIku "UnpakSiamida/modules/indikatorrenstra/application/GetIndikatorRenstraIku"
	saveIndikatorIku "UnpakSiamida/modules/indikatorrenstra/application/SaveIndikatorRenstraIku"
	setupUuid "UnpakSiamida/modules/indikatorrenstra/application/SetupUuidIndikatorRenstra"
	getTree "UnpakSiamida/modules/indikatorrenstra/application/TreeIndikatorRenstra"
	update "UnpakSiamida/modules/indikatorrenstra/application/UpdateIndikatorRenstra"
	domainindikatorrenstra "UnpakSiamida/modules/indikatorrenstra/domain"

	infraStandarRenstra "UnpakSiamida/modules/standarrenstra/infrastructure"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"
)

func RegisterModuleIndikatorRenstra(db *gorm.DB, dbs ...*gorm.DB) error {
	var dbIku *gorm.DB
	if len(dbs) > 0 {
		dbIku = dbs[0]
	}
	repoIndikatorRenstra := NewIndikatorRenstraRepository(db, dbIku)
	repoStandarRenstra := infraStandarRenstra.NewStandarRenstraRepository(db)
	// if err := db.AutoMigrate(&domainindikatorrenstra.IndikatorRenstra{}); err != nil {
	// 	panic(err)
	// }

	// Register request handler
	mediatr.RegisterRequestHandler[
		create.CreateIndikatorRenstraCommand,
		string,
	](&create.CreateIndikatorRenstraCommandHandler{
		Repo:               repoIndikatorRenstra,
		RepoStandarRenstra: repoStandarRenstra,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateIndikatorRenstraCommand,
		string,
	](&update.UpdateIndikatorRenstraCommandHandler{
		Repo:               repoIndikatorRenstra,
		RepoStandarRenstra: repoStandarRenstra,
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
		commondomain.Paged[domainindikatorrenstra.IndikatorRenstraDefault],
	](&getAll.GetAllIndikatorRenstrasQueryHandler{
		Repo: repoIndikatorRenstra,
	})

	mediatr.RegisterRequestHandler[
		getTree.TreeIndikatorRenstraQuery,
		[]domainindikatorrenstra.IndikatorTree,
	](&getTree.TreeIndikatorRenstraQueryHandler{
		Repo: repoIndikatorRenstra,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidIndikatorRenstraCommand,
		string,
	](&setupUuid.SetupUuidIndikatorRenstraCommandHandler{
		Repo: repoIndikatorRenstra,
	})

	mediatr.RegisterRequestHandler[
		getAllMasterIku.GetAllMasterIkuQuery,
		[]domainindikatorrenstra.MasterIku,
	](&getAllMasterIku.GetAllMasterIkuQueryHandler{
		Repo: repoIndikatorRenstra,
	})

	mediatr.RegisterRequestHandler[
		getIndikatorIku.GetIndikatorRenstraIkuQuery,
		*getIndikatorIku.GetIndikatorRenstraIkuResult,
	](&getIndikatorIku.GetIndikatorRenstraIkuQueryHandler{
		Repo: repoIndikatorRenstra,
	})

	mediatr.RegisterRequestHandler[
		saveIndikatorIku.SaveIndikatorRenstraIkuCommand,
		*saveIndikatorIku.SaveIndikatorRenstraIkuResult,
	](&saveIndikatorIku.SaveIndikatorRenstraIkuCommandHandler{
		Repo: repoIndikatorRenstra,
	})

	mediatr.RegisterRequestHandler[
		getDashboard.GetDashboardIkuAnggaranQuery,
		*getDashboard.GetDashboardIkuAnggaranResult,
	](&getDashboard.GetDashboardIkuAnggaranQueryHandler{
		Repo: repoIndikatorRenstra,
	})

	commoninfra.RegisterValidation(create.CreateIndikatorRenstraCommandValidation, "IndikatorRenstraCreate.Validation")
	commoninfra.RegisterValidation(update.UpdateIndikatorRenstraCommandValidation, "IndikatorRenstraUpdate.Validation")
	commoninfra.RegisterValidation(delete.DeleteIndikatorRenstraCommandValidation, "IndikatorRenstraDelete.Validation")
	commoninfra.RegisterValidation(saveIndikatorIku.SaveIndikatorRenstraIkuCommandValidation, "IndikatorRenstraSaveIku.Validation")

	return nil
}
