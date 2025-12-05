package main

import (
	"UnpakSiamida/modules/user/infrastructure"
	"UnpakSiamida/modules/user/presentation"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:      true, // gunakan semua CPU cores
		ServerHeader: "Fiber",
		// ReadTimeout: 10 * time.Second,
		// WriteTimeout: 10 * time.Second,
		// IdleTimeout: 10 * time.Second
	})

	infrastructure.RegisterModule()
	presentation.ModuleUser(app)

	app.Listen(":3000")
}
