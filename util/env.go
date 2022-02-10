package util

import "os"

func IsDevEnvironment() bool {
	return os.Getenv("APP_ENV") == "DEV"
}

func IsTestEnvironment() bool {
	return os.Getenv("APP_ENV") == "TEST" || os.Getenv("CI") == "true"
}

func IsKubernetesTestEnvironment() bool {
	return os.Getenv("APP_ENV") == "TEST_KUBERNETES"
}

func IsDockerTestEnvironment() bool {
	return os.Getenv("APP_ENV") == "TEST_DOCKER"
}
