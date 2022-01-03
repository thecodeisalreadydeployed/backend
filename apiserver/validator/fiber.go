package validator

import "github.com/gofiber/fiber/v2"

func ParseBody(ctx *fiber.Ctx, body interface{}) *fiber.Error {
	if err := ctx.BodyParser(body); err != nil {
		return fiber.ErrBadRequest
	}

	return nil
}

func ParseBodyAndValidate(ctx *fiber.Ctx, body interface{}) *fiber.Error {
	if err := ParseBody(ctx, body); err != nil {
		return err
	}

	if validationErrors := CheckStruct(body); len(validationErrors) > 0 {
		ctx.Status(fiber.StatusBadRequest).JSON(validationErrors)
		return fiber.ErrBadRequest
	}

	return nil
}
