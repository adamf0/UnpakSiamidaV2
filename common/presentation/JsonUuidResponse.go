package presentation

import "github.com/gofiber/fiber/v2"

func JsonUUID(c *fiber.Ctx, uuid string) error {
	return c.JSON(fiber.Map{"uuid": uuid})
}
