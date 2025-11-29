package main

import (
	"UnpakSiamida/modules/akurasipenelitian/infrastructure"
	"UnpakSiamida/modules/akurasipenelitian/presentation"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	infrastructure.RegisterModule()
	presentation.ModuleUser(app)

	app.Listen(":3000")
}
