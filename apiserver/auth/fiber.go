package auth

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/util"
)

func EnsureAuthenticated() func(c *fiber.Ctx) error {
	if util.IsDevEnvironment() || util.IsTestEnvironment() {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	return adaptor.HTTPMiddleware(ensureValidToken())
}
