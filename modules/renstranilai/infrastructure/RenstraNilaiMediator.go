package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	delete "UnpakSiamida/modules/renstranilai/application/DeleteRenstraNilai"
	getAll "UnpakSiamida/modules/renstranilai/application/GetAllRenstraNilais"
	get "UnpakSiamida/modules/renstranilai/application/GetRenstraNilai"
	setupUuid "UnpakSiamida/modules/renstranilai/application/SetupUuidRenstraNilai"
	update "UnpakSiamida/modules/renstranilai/application/UpdateRenstraNilai"
	domainrenstranilai "UnpakSiamida/modules/renstranilai/domain"

	infraRenstra "UnpakSiamida/modules/renstra/infrastructure"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterModuleRenstraNilai(db *gorm.DB) error {
	repoRenstraNilai := NewRenstraNilaiRepository(db)
	repoRenstra := infraRenstra.NewRenstraRepository(db)

	mediatr.RegisterRequestHandler[
		update.UpdateRenstraNilaiCommand,
		string,
	](&update.UpdateRenstraNilaiCommandHandler{
		Repo:        repoRenstraNilai,
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
		commondomain.Paged[domainrenstranilai.RenstraNilaiDefault],
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
