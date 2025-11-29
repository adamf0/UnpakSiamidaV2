package presentation

import (
    "context"
    "errors"
    
    "github.com/gofiber/fiber/v2"
    "github.com/mehdihadeli/go-mediatr"

    "UnpakSiamida/modules/akurasipenelitian/application"
    commoninfra "UnpakSiamida/common/infrastructure"
    domain "UnpakSiamida/common/domain"
)

func ModuleUser(app *fiber.App) {
    app.Post("/user", func(c *fiber.Ctx) error {
            // Ambil form values
        nama := c.FormValue("nama")
        skor := c.FormValue("skor") 

        cmd := application.CreateUserCommand{
            Nama: nama,
            Skor: skor,
        }

        id, err := mediatr.Send[application.CreateUserCommand, string](context.Background(), cmd)
        if err != nil {
            // 1) coba unwrap ResponseError
            var respErr *commoninfra.ResponseError
            if errors.As(err, &respErr) {
                // map code -> http status jika perlu. contoh default 400 untuk validation
                status := 400
                switch respErr.Code {
                case "AkurasiPenelitian.Validation":
                    status = 400
                case "AkurasiPenelitian.NotFound":
                    status = 404
                case "AkurasiPenelitian.Conflict":
                    status = 409
                default:
                    status = 400
                }
                return c.Status(status).JSON(respErr)
            }

            // 2) coba unwrap domain.Error (jika behavior mengembalikan domain.Error secara langsung)
            var derr domain.Error
            if errors.As(err, &derr) {
                re := commoninfra.NewResponseError(derr.Code, derr.Description)
                // map domain.Error type -> status (NotFound -> 404, ...)
                status := 400
                switch derr.Type {
                case domain.NotFound:
                    status = 404
                case domain.Conflict:
                    status = 409
                case domain.Validation:
                    status = 400
                default:
                    status = 500
                }
                return c.Status(status).JSON(re)
            }

            // 3) fallback: unknown/wrapped error -> internal server error
            return c.Status(500).JSON(commoninfra.NewInternalError(err))
        }

        return c.JSON(fiber.Map{"id": id})
    })
}
