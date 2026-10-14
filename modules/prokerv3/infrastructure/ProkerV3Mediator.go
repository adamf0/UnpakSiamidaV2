package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	create "UnpakSiamida/modules/prokerv3/application/CreateProkerV3"
	delete "UnpakSiamida/modules/prokerv3/application/DeleteProkerV3"
	getAll "UnpakSiamida/modules/prokerv3/application/GetAllProkerV3"
	get "UnpakSiamida/modules/prokerv3/application/GetProkerV3"
	update "UnpakSiamida/modules/prokerv3/application/UpdateProkerV3"
	domain "UnpakSiamida/modules/prokerv3/domain"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterModuleProkerV3(db *gorm.DB) error {
	repoProker := NewProkerV3Repository(db)

	mediatr.RegisterRequestHandler[
		create.CreateProkerV3Command,
		string,
	](&create.CreateProkerV3CommandHandler{
		Repo: repoProker,
		DB:   db,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateProkerV3Command,
		string,
	](&update.UpdateProkerV3CommandHandler{
		Repo: repoProker,
		DB:   db,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteProkerV3Command,
		string,
	](&delete.DeleteProkerV3CommandHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		get.GetProkerV3ByUuidQuery,
		*domain.ProkerV3Default,
	](&get.GetProkerV3ByUuidQueryHandler{
		Repo: repoProker,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllProkerV3Query,
		commondomain.Paged[domain.ProkerV3Default],
	](&getAll.GetAllProkerV3QueryHandler{
		Repo: repoProker,
	})

	return nil
}
