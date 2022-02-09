package auth

import (
	"context"
	"net/url"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/thecodeisalreadydeployed/config"
	"go.uber.org/zap"
)

type CustomClaims struct {
	Scope string `json:"scope"`
}

func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func EnsureValidToken() func(c *fiber.Ctx) error {
	issuerURL, err := url.Parse("https://" + config.Auth0Domain() + "/")
	if err != nil {
		zap.L().Error("failed to parse issuer URL", zap.Error(err))
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{config.Auth0Audience()},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)

	if err != nil {
		zap.L().Error("failed to set up the JWT validator", zap.Error(err))
	}
}
