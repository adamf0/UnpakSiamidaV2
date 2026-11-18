package infrastructure

import (
	createCmd "UnpakSiamida/modules/mataprogramv2/application/CreateMataProgramV2"
	deleteCmd "UnpakSiamida/modules/mataprogramv2/application/DeleteMataProgramV2"
	getAllQuery "UnpakSiamida/modules/mataprogramv2/application/GetAllMataProgramV2"
	getByUuidQuery "UnpakSiamida/modules/mataprogramv2/application/GetMataProgramV2"
	restoreCmd "UnpakSiamida/modules/mataprogramv2/application/RestoreMataProgramV2"
	updateCmd "UnpakSiamida/modules/mataprogramv2/application/UpdateMataProgramV2"
	domain "UnpakSiamida/modules/mataprogramv2/domain"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterModuleMataProgramV2(db *gorm.DB) error {
	repo := NewMataProgramV2Repository(db)

	mediatr.RegisterRequestHandler[
		getAllQuery.GetAllMataProgramV2Query,
		*getAllQuery.GetAllMataProgramV2Result,
	](&getAllQuery.GetAllMataProgramV2QueryHandler{
		Repo: repo,
	})

	mediatr.RegisterRequestHandler[
		getByUuidQuery.GetMataProgramV2ByUuidQuery,
		*domain.MataProgramV2Detail,
	](&getByUuidQuery.GetMataProgramV2ByUuidQueryHandler{
		Repo: repo,
	})

	mediatr.RegisterRequestHandler[
		createCmd.CreateMataProgramV2Command,
		*createCmd.CreateMataProgramV2Result,
	](&createCmd.CreateMataProgramV2CommandHandler{
		Repo: repo,
	})

	mediatr.RegisterRequestHandler[
		updateCmd.UpdateMataProgramV2Command,
		*updateCmd.UpdateMataProgramV2Result,
	](&updateCmd.UpdateMataProgramV2CommandHandler{
		Repo: repo,
	})

	mediatr.RegisterRequestHandler[
		deleteCmd.DeleteMataProgramV2Command,
		*deleteCmd.DeleteMataProgramV2Result,
	](&deleteCmd.DeleteMataProgramV2CommandHandler{
		Repo: repo,
	})

	mediatr.RegisterRequestHandler[
		restoreCmd.RestoreMataProgramV2Command,
		*restoreCmd.RestoreMataProgramV2Result,
	](&restoreCmd.RestoreMataProgramV2CommandHandler{
		Repo: repo,
	})

	return nil
}
