package group

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"github.com/thecodeisalreadydeployed/datastore"
	"github.com/thecodeisalreadydeployed/preset"
)

func Preset(c *fiber.Ctx) error {
	framework := c.Params("framework")
	text := preset.Text(preset.Framework(framework))
	return c.SendString(text)
}

func Health(c *fiber.Ctx) error {
	ok := datastore.IsReady()
	return c.JSON(map[string]string{"ok": cast.ToString(ok)})
}
