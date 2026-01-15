package infrastructure

import (
	create "UnpakSiamida/modules/user/application/CreateUser"
	delete "UnpakSiamida/modules/user/application/DeleteUser"
	getAll "UnpakSiamida/modules/user/application/GetAllUsers"
	get "UnpakSiamida/modules/user/application/GetUser"
	setupUuid "UnpakSiamida/modules/user/application/SetupUuidUser"
	update "UnpakSiamida/modules/user/application/UpdateUser"
	domainuser "UnpakSiamida/modules/user/domain"

	"github.com/mehdihadeli/go-mediatr"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "fmt"

	commoninfra "UnpakSiamida/common/infrastructure"
	eventUser "UnpakSiamida/modules/user/event"

	infraFakultasUnit "UnpakSiamida/modules/fakultasunit/infrastructure"
)

func RegisterModuleUser(db *gorm.DB, tg commoninfra.TelegramSender) error {
	// dsn := "root:@tcp(127.0.0.1:3306)/unpak_sijamu_server?charset=utf8mb4&parseTime=true&loc=Local"

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	//     return fmt.Errorf("User DB connection failed: %w", err)
	// 	// panic(err)
	// }

	repoUser := NewUserRepository(db)
	repoFakultasUnit := infraFakultasUnit.NewFakultasUnitRepository(db)
	// if err := db.AutoMigrate(&domainuser.User{}); err != nil {
	// 	panic(err)
	// }

	// Pipeline behavior
	// mediatr.RegisterRequestPipelineBehaviors(NewValidationBehaviorUser())

	// Register request handler
	mediatr.RegisterRequestHandler[
		create.CreateUserCommand,
		string,
	](&create.CreateUserCommandHandler{
		Repo:             repoUser,
		RepoFakultasUnit: repoFakultasUnit,
	})

	mediatr.RegisterRequestHandler[
		update.UpdateUserCommand,
		string,
	](&update.UpdateUserCommandHandler{
		Repo:             repoUser,
		RepoFakultasUnit: repoFakultasUnit,
	})

	mediatr.RegisterRequestHandler[
		delete.DeleteUserCommand,
		string,
	](&delete.DeleteUserCommandHandler{
		Repo: repoUser,
	})

	mediatr.RegisterRequestHandler[
		get.GetUserByUuidQuery,
		*domainuser.User,
	](&get.GetUserByUuidQueryHandler{
		Repo: repoUser,
	})

	mediatr.RegisterRequestHandler[
		getAll.GetAllUsersQuery,
		domainuser.PagedUsers,
	](&getAll.GetAllUsersQueryHandler{
		Repo: repoUser,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidUserCommand,
		string,
	](&setupUuid.SetupUuidUserCommandHandler{
		Repo: repoUser,
	})

	// =========================
	// Domain Event Handler
	// =========================
	commoninfra.RegisterDomainEvent(&eventUser.UserUpdatedEvent{})

	mediatr.RegisterNotificationHandler[eventUser.UserUpdatedEvent](
		eventUser.NewUserUpdatedEventHandler(tg),
	)

	commoninfra.RegisterDomainEvent(&eventUser.UserCreatedEvent{})

	mediatr.RegisterNotificationHandler[eventUser.UserCreatedEvent](
		eventUser.NewUserCreatedEventHandler(tg),
	)

	return nil
}
