package config

import (
	"os"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/thecodeisalreadydeployed/constant"
	"go.uber.org/zap"
)

func DefaultGitSignature() *object.Signature {
	return &object.Signature{
		Name:  "GitHub Action",
		Email: "action@github.com",
		When:  time.Now(),
	}
}

const (
	DefaultKanikoWorkingDirectory string = "/__w/kaniko/"
)

func DefaultUserspaceRepository() string {
	if len(viper.GetString(constant.UserspaceRepository)) != 0 {
		repo := viper.GetString(constant.UserspaceRepository)
		zap.L().Info("userspace repository: " + repo)
		return repo
	}

	dir, err := os.MkdirTemp("", uuid.NewString()+"-")
	if err != nil {
		zap.L().Fatal("cannot create temporary directory", zap.Error(err))
	}

	zap.L().Info("created temporary directory: " + dir)

	viper.Set(constant.UserspaceRepository, dir)

	return dir
}

func ArgoCDServerHost() string {
	viper.SetDefault(constant.ArgoCDServerHostEnv, "http://argocd-server.argocd.svc.cluster.local")
	return viper.GetString(constant.ArgoCDServerHostEnv)
}

func Auth0Domain() string {
	return viper.GetString("AUTH0_DOMAIN")
}

func Auth0Audience() string {
	return viper.GetString("AUTH0_AUDIENCE")
}

func SetDefault() {
	viper.SetDefault(constant.UserspaceRepository, "/__w/userspace")
}
