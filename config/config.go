package config

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/thecodeisalreadydeployed/constant"
	containerregistry "github.com/thecodeisalreadydeployed/containerregistry/types"
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
	if len(viper.GetString(constant.USERSPACE_REPOSITORY)) != 0 {
		repo := viper.GetString(constant.USERSPACE_REPOSITORY)
		zap.L().Info("userspace repository: " + repo)
		return repo
	}

	dir, err := os.MkdirTemp("", uuid.NewString()+"-")
	if err != nil {
		zap.L().Fatal("cannot create temporary directory", zap.Error(err))
	}

	zap.L().Info("created temporary directory: " + dir)

	viper.Set(constant.USERSPACE_REPOSITORY, dir)

	return dir
}

func ArgoCDServerHost() string {
	viper.SetDefault(constant.ARGOCD_SERVER_HOST, "http://argocd-server.argocd.svc.cluster.local")
	return viper.GetString(constant.ARGOCD_SERVER_HOST)
}

func GitServerHost() string {
	return viper.GetString(constant.GITSERVER_HOST)
}

func Auth0Domain() string {
	return viper.GetString(constant.AUTH0_DOMAIN)
}

func Auth0Audience() string {
	return viper.GetString(constant.AUTH0_AUDIENCE)
}

func FirebaseServiceAccountKey() string {
	return viper.GetString(constant.FIREBASE_SERVICE_ACCOUNT_KEY)
}

func DefaultContainerRegistryConfiguration() containerregistry.ContainerRegistryConfiguration {
	return containerregistry.ContainerRegistryConfiguration{
		Type:                 containerregistry.GCR,
		AuthenticationMethod: containerregistry.KubernetesServiceAccount,
		// The value of Secret field depends on the value of AuthenticationMethod field.
		// If the value of AuthenticationMethod field is KubernetesServiceAccount, the
		// value of Secret field is the Kubernetes service account. Otherwise, the value
		// of Secret field is cloud provider's credentials.
		Secret:     "codedeploy-imagebuilder",
		Repository: "senior-project",
		Metadata: map[string]string{
			"GOOGLE_CLOUD_PROJECT": "deploys-dev",
			"GOOGLE_CLOUD_REGION":  "asia-southeast1",
		},
	}
}

func KindLocalRegistryConfiguration() containerregistry.ContainerRegistryConfiguration {
	return containerregistry.ContainerRegistryConfiguration{
		Type:                 containerregistry.LOCAL,
		AuthenticationMethod: containerregistry.Secret,
		Secret:               "",
		Repository:           "",
		Metadata: map[string]string{
			"HOSTNAME": "localhost",
			"PORT":     fmt.Sprintf("%d", 5000),
		},
	}
}

func BindEnv() {
	_ = viper.BindEnv(constant.ARGOCD_SERVER_HOST)
	_ = viper.BindEnv(constant.AUTH0_AUDIENCE)
	_ = viper.BindEnv(constant.AUTH0_DOMAIN)
	_ = viper.BindEnv(constant.DATABASE_HOST)
	_ = viper.BindEnv(constant.DATABASE_NAME)
	_ = viper.BindEnv(constant.DATABASE_PASSWORD)
	_ = viper.BindEnv(constant.DATABASE_PORT)
	_ = viper.BindEnv(constant.DATABASE_USERNAME)
	_ = viper.BindEnv(constant.GITSERVER_HOST)
	_ = viper.BindEnv(constant.FIREBASE_SERVICE_ACCOUNT_KEY)
	_ = viper.BindEnv(constant.USERSPACE_REPOSITORY)
}
