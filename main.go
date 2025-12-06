package main

import (
	userInfrastructure "UnpakSiamida/modules/user/infrastructure"
	userPresentation "UnpakSiamida/modules/user/presentation"

	standarrenstraInfrastructure "UnpakSiamida/modules/standarrenstra/infrastructure"
	standarrenstraPresentation "UnpakSiamida/modules/standarrenstra/presentation"
	
	createUser "UnpakSiamida/modules/user/application/CreateUser"
	updateUser "UnpakSiamida/modules/user/application/UpdateUser"
	deleteUser "UnpakSiamida/modules/user/application/DeleteUser"

	createStandarRenstra "UnpakSiamida/modules/standarrenstra/application/CreateStandarRenstra"
	updateStandarRenstra "UnpakSiamida/modules/standarrenstra/application/UpdateStandarRenstra"
	deleteStandarRenstra "UnpakSiamida/modules/standarrenstra/application/DeleteStandarRenstra"

	"github.com/gofiber/fiber/v2"
	"context"
	"strings"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/go-ozzo/ozzo-validation/v4"
	commoninfra "UnpakSiamida/common/infrastructure"
	commonpresentation "UnpakSiamida/common/presentation"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:      true, // gunakan semua CPU cores
		ServerHeader: "Fiber",
		// ReadTimeout: 10 * time.Second,
		// WriteTimeout: 10 * time.Second,
		// IdleTimeout: 10 * time.Second
	})
	app.Use(commonpresentation.LoggerMiddleware)

	userInfrastructure.RegisterModuleUser()
	userPresentation.ModuleUser(app)

	mediatr.RegisterRequestPipelineBehaviors(NewValidationBehavior())

	standarrenstraInfrastructure.RegisterModuleStandarRenstra()
	standarrenstraPresentation.ModuleStandarRenstra(app)

	app.Listen(":3000")
}

type ValidationBehavior struct{}

func NewValidationBehavior() *ValidationBehavior {
	return &ValidationBehavior{}
}

func (b *ValidationBehavior) Handle(
	ctx context.Context,
	request interface{},
	next mediatr.RequestHandlerFunc,
) (interface{}, error) {

	switch cmd := request.(type) {
		// === User Commands ===
		case createUser.CreateUserCommand:
			if err := createUser.CreateUserCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("User.Validation", err)
			}
		case updateUser.UpdateUserCommand:
			if err := updateUser.UpdateUserCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("User.Validation", err)
			}
		case deleteUser.DeleteUserCommand:
			if err := deleteUser.DeleteUserCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("User.Validation", err)
			}

		// === StandarRenstra Commands ===
		case createStandarRenstra.CreateStandarRenstraCommand:
			if err := createStandarRenstra.CreateStandarRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("StandarRenstra.Validation", err)
			}
		case updateStandarRenstra.UpdateStandarRenstraCommand:
			if err := updateStandarRenstra.UpdateStandarRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("StandarRenstra.Validation", err)
			}
		case deleteStandarRenstra.DeleteStandarRenstraCommand:
			if err := deleteStandarRenstra.DeleteStandarRenstraCommandValidation(cmd); err != nil {
				return nil, wrapValidationError("StandarRenstra.Validation", err)
			}

		// default:
		// 	// request lain → skip validation
	}

	return next(ctx)
}

func wrapValidationError(code string, err error) error {
	if ve, ok := err.(validation.Errors); ok {
		msgs := make(map[string]string)
		for field, ferr := range ve {
			key := strings.ToLower(field)
			msgs[key] = ferr.Error()
		}
		return commoninfra.NewResponseError(code, msgs)
	}
	return commoninfra.NewResponseError(code, err.Error())
}