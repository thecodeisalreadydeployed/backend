package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/util"
)

func EnsureAuthenticated(firebaseAuth *auth.Client) func(c *fiber.Ctx) error {
	if util.IsDevEnvironment() || util.IsTestEnvironment() {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	return func(c *fiber.Ctx) error {
		isInternalRequest := string(c.Request().Header.Peek("X-CodeDeploy-Internal-Request")) == "True" && len(c.Request().Header.Peek("X-Forwarded-For")) == 0
		if isInternalRequest {
			c.Next()
		}

		parts := strings.Split(c.GetReqHeaders()["Authorization"], " ")
		if len(parts) != 2 {
			fmt.Printf("parts: %v\n", parts)
			return c.SendStatus(http.StatusUnauthorized)
		}

		_, err := firebaseAuth.VerifyIDToken(context.Background(), parts[1])
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return c.SendStatus(http.StatusUnauthorized)
		}

		return c.Next()
	}
}
