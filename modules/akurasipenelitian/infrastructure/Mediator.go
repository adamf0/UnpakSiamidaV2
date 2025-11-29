package infrastructure

import (
    "UnpakSiamida/modules/akurasipenelitian/application"
    "github.com/mehdihadeli/go-mediatr"
)

func RegisterModule() {
    // Pipeline behavior
    mediatr.RegisterRequestPipelineBehaviors(NewValidationBehavior())

    // Register request handler
    mediatr.RegisterRequestHandler[
        application.CreateUserCommand,
        string,
    ](&application.CreateUserCommandHandler{})
}
