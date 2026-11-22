package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	commoninfra "UnpakSiamida/common/infrastructure"
	getAll "UnpakSiamida/modules/fakultasunit/application/GetAllFakultasUnits"
	get "UnpakSiamida/modules/fakultasunit/application/GetFakultasUnit"
	setupUuid "UnpakSiamida/modules/fakultasunit/application/SetupUuidFakultasUnit"
	update "UnpakSiamida/modules/fakultasunit/application/UpdateFakultasUnit"
	domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
	"context"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func RegisterModuleFakultasUnit(db *gorm.DB) error {
	repoFakultasUnit := NewFakultasUnitRepository(db)

	if err := repoFakultasUnit.EnsureSchema(context.Background()); err != nil {
		return err
	}

	mediatr.RegisterRequestHandler[
		get.GetFakultasUnitByUuidQuery,
		*domainfakultasunit.FakultasUnit,
	](&get.GetFakultasUnitByUuidQueryHandler{
		Repo: repoFakultasUnit,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllFakultasUnitsQuery,
		commondomain.Paged[domainfakultasunit.FakultasUnit],
	](&getAll.GetAllFakultasUnitsQueryHandler{
		Repo: repoFakultasUnit,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateFakultasUnitCommand,
		*domainfakultasunit.FakultasUnit,
	](&update.UpdateFakultasUnitCommandHandler{
		Repo: repoFakultasUnit,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidFakultasUnitCommand,
		string,
	](&setupUuid.SetupUuidFakultasUnitCommandHandler{
		Repo: repoFakultasUnit,
	})

	commoninfra.RegisterValidation(update.UpdateFakultasUnitCommandValidation, "FakultasUnitUpdate.Validation")

	return nil
}
