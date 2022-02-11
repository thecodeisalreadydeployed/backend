package auth

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/thecodeisalreadydeployed/config"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

func SetupFirebase() *auth.Client {
	opt := option.WithCredentialsJSON([]byte(config.FirebaseServiceAccountKey()))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		zap.L().Fatal("cannot initialize Firebase Admin SDK", zap.Error(err))
	}

	auth, err := app.Auth(context.Background())
	if err != nil {
		zap.L().Fatal("cannot initialize Firebase Admin SDK", zap.Error(err))
	}

	return auth
}
