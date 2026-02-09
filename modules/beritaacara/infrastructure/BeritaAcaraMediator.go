package infrastructure

import (
	create "UnpakSiamida/modules/beritaacara/application/CreateBeritaAcara"
	delete "UnpakSiamida/modules/beritaacara/application/DeleteBeritaAcara"
	export "UnpakSiamida/modules/beritaacara/application/ExportBeritaAcara"
	getAll "UnpakSiamida/modules/beritaacara/application/GetAllBeritaAcaras"
	get "UnpakSiamida/modules/beritaacara/application/GetBeritaAcara"
	setupUuid "UnpakSiamida/modules/beritaacara/application/SetupUuidBeritaAcara"
	update "UnpakSiamida/modules/beritaacara/application/UpdateBeritaAcara"
	domainBeritaAcara "UnpakSiamida/modules/beritaacara/domain"

	commonDomain "UnpakSiamida/common/domain"
	commonInfra "UnpakSiamida/common/infrastructure"
	eventBeritaAcara "UnpakSiamida/modules/beritaacara/event"
	infraFakultasUnit "UnpakSiamida/modules/fakultasunit/infrastructure"
	infraUser "UnpakSiamida/modules/user/infrastructure"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"
)

func RegisterModuleBeritaAcara(db *gorm.DB, redis *commonDomain.IRedisStore) error {

	// dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("Indikator Renstra DB connection failed: %w", err)
	// 	// panic(err)
	// }

	repoBeritaAcara := NewBeritaAcaraRepository(db)
	repoFakultasUnit := infraFakultasUnit.NewFakultasUnitRepository(db)
	repoUser := infraUser.NewUserRepository(db)
	// if err := db.AutoMigrate(&domainBeritaAcara.BeritaAcara{}); err != nil {
	// 	panic(err)
	// }

	// Pipeline behavior
	// mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorBeritaAcara())

	// Register request handler
	mediatr.RegisterRequestHandler[
		export.PublishBeritaAcaraCommand,
		string,
	](&export.PublishBeritaAcaraCommandHandler{
		Repo: repoBeritaAcara,
	})

	mediatr.RegisterRequestHandler[
		export.ExportBeritaAcaraCommand,
		[]byte,
	](&export.ExportBeritaAcaraCommandHandler{
		Repo:  repoBeritaAcara,
		Redis: *redis,
	})

	mediatr.RegisterRequestHandler[
		create.CreateBeritaAcaraCommand,
		string,
	](&create.CreateBeritaAcaraCommandHandler{
		Repo:             repoBeritaAcara,
		RepoFakultasUnit: repoFakultasUnit,
		RepoUser:         repoUser,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateBeritaAcaraCommand,
		string,
	](&update.UpdateBeritaAcaraCommandHandler{
		Repo:             repoBeritaAcara,
		RepoFakultasUnit: repoFakultasUnit,
		RepoUser:         repoUser,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteBeritaAcaraCommand,
		string,
	](&delete.DeleteBeritaAcaraCommandHandler{
		Repo: repoBeritaAcara,
	})

	mediatr.RegisterRequestHandler[
		get.GetBeritaAcaraByUuidQuery,
		*domainBeritaAcara.BeritaAcara,
	](&get.GetBeritaAcaraByUuidQueryHandler{
		Repo: repoBeritaAcara,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllBeritaAcarasQuery,
		domainBeritaAcara.PagedBeritaAcaras,
	](&getAll.GetAllBeritaAcarasQueryHandler{
		Repo: repoBeritaAcara,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidBeritaAcaraCommand,
		string,
	](&setupUuid.SetupUuidBeritaAcaraCommandHandler{
		Repo: repoBeritaAcara,
	})

	// =========================
	// Domain Event Handler
	// =========================
	commonInfra.RegisterDomainEvent(&eventBeritaAcara.BeritaAcaraPdfRequestedEvent{})

	mediatr.RegisterNotificationHandler[eventBeritaAcara.BeritaAcaraPdfRequestedEvent](
		eventBeritaAcara.NewBeritaAcaraPdfRequestedEventHandler(*redis),
	)

	return nil
}
