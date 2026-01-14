package infrastructure

import (
    domainfakultasunit "UnpakSiamida/modules/fakultasunit/domain"
    get "UnpakSiamida/modules/fakultasunit/application/GetFakultasUnit"
    getAll "UnpakSiamida/modules/fakultasunit/application/GetAllFakultasUnits"
    setupUuid "UnpakSiamida/modules/fakultasunit/application/SetupUuidFakultasUnit"
    "github.com/mehdihadeli/go-mediatr"
    // "gorm.io/driver/mysql"
	"gorm.io/gorm"
    // "fmt"
)

func RegisterModuleFakultasUnit(db *gorm.DB) error{
    // dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Standar Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

    repoFakultasUnit := NewFakultasUnitRepository(db)
	// if err := db.AutoMigrate(&domainfakultasunit.FakultasUnit{}); err != nil {
	// 	panic(err)
	// }

    // Pipeline behavior
    // mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorFakultasUnit())

    // Register request handler
    mediatr.RegisterRequestHandler[
        get.GetFakultasUnitByUuidQuery,
        *domainfakultasunit.FakultasUnit,
    ](&get.GetFakultasUnitByUuidQueryHandler{
        Repo: repoFakultasUnit,
    })

    mediatr.RegisterRequestHandler[
        getAll.GetAllFakultasUnitsQuery,
        domainfakultasunit.PagedFakultasUnits,
    ](&getAll.GetAllFakultasUnitsQueryHandler{
        Repo: repoFakultasUnit,
    })

    mediatr.RegisterRequestHandler[
        setupUuid.SetupUuidFakultasUnitCommand,
        string,
    ](&setupUuid.SetupUuidFakultasUnitCommandHandler{
        Repo: repoFakultasUnit,
    })

    return nil
}
