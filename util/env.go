package util

import "os"

func IsDevEnvironment() bool {
	return os.Getenv("APP_ENV") == "DEV"
}