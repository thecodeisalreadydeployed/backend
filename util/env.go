package util

import (
	"os"
	"strings"
)

func IsDevEnvironment() bool {
	return os.Getenv("APP_ENV") == "DEV"
}

func IsProductionEnvironment() bool {
	return os.Getenv("APP_ENV") == "PROD"
}

func IsTestEnvironment() bool {
	return strings.HasPrefix(os.Getenv("APP_ENV"), "TEST")
}

func IsKubernetesTestEnvironment() bool {
	return os.Getenv("APP_ENV") == "TEST_KUBERNETES"
}

func IsDockerTestEnvironment() bool {
	return os.Getenv("APP_ENV") == "TEST_DOCKER"
}
