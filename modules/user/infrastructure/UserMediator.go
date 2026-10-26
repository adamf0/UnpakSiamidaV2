package infrastructure

import (
	commondomain "UnpakSiamida/common/domain"
	create "UnpakSiamida/modules/user/application/CreateUser"
	delete "UnpakSiamida/modules/user/application/DeleteUser"
	getAllOptions "UnpakSiamida/modules/user/application/GetAllUserOptions"
	getAll "UnpakSiamida/modules/user/application/GetAllUsers"
	get "UnpakSiamida/modules/user/application/GetUser"
	merge "UnpakSiamida/modules/user/application/MergeUser"
	setupUuid "UnpakSiamida/modules/user/application/SetupUuidUser"
	update "UnpakSiamida/modules/user/application/UpdateUser"
	domainuser "UnpakSiamida/modules/user/domain"

	createAbac "UnpakSiamida/modules/user/application/CreateUserAbac"
	deleteAbac "UnpakSiamida/modules/user/application/DeleteUserAbac"
	getAllAbac "UnpakSiamida/modules/user/application/GetAllUserAbac"
	updateAbac "UnpakSiamida/modules/user/application/UpdateUserAbac"

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

	if db != nil {
		_ = db.Exec("ALTER TABLE `users` MODIFY COLUMN `level` VARCHAR(100) NOT NULL DEFAULT 'user'")
		_ = db.Exec("ALTER TABLE `user_abac` MODIFY COLUMN `id_fakultas_unit` int(10) UNSIGNED NULL DEFAULT NULL")
	}

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
		commondomain.Paged[domainuser.User],
	](&getAll.GetAllUsersQueryHandler{
		Repo: repoUser,
	})

	mediatr.RegisterRequestHandler[
		getAllOptions.GetAllUserOptionsQuery,
		[]domainuser.UserOptions,
	](&getAllOptions.GetAllUserOptionsQueryHandler{
		Repo: repoUser,
	})

	mediatr.RegisterRequestHandler[
		setupUuid.SetupUuidUserCommand,
		string,
	](&setupUuid.SetupUuidUserCommandHandler{
		Repo: repoUser,
	})

	mediatr.RegisterRequestHandler[
		merge.MergeUserCommand,
		string,
	](&merge.MergeUserCommandHandler{
		Repo: repoUser,
	})

	// User ABAC Handlers
	mediatr.RegisterRequestHandler[
		getAllAbac.GetAllUserAbacQuery,
		[]domainuser.UserAbac,
	](&getAllAbac.GetAllUserAbacQueryHandler{
		Repo: repoUser,
	})

	mediatr.RegisterRequestHandler[
		createAbac.CreateUserAbacCommand,
		*createAbac.CreateUserAbacResult,
	](&createAbac.CreateUserAbacCommandHandler{
		Repo: repoUser,
	})

	mediatr.RegisterRequestHandler[
		updateAbac.UpdateUserAbacCommand,
		string,
	](&updateAbac.UpdateUserAbacCommandHandler{
		Repo: repoUser,
	})

	mediatr.RegisterRequestHandler[
		deleteAbac.DeleteUserAbacCommand,
		string,
	](&deleteAbac.DeleteUserAbacCommandHandler{
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

	commoninfra.RegisterValidation(create.CreateUserCommandValidation, "UserCreate.Validation")
	commoninfra.RegisterValidation(update.UpdateUserCommandValidation, "UserUpdate.Validation")
	commoninfra.RegisterValidation(delete.DeleteUserCommandValidation, "UserDelete.Validation")

	return nil
}
