package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/util"
)

func EnsureAuthenticated() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if util.IsDevEnvironment() {
			return c.Next()
		}
		return c.Next()
	}
}
